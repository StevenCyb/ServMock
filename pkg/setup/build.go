package parser

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/StevenCyb/ServMock/pkg/ini"
	"github.com/StevenCyb/ServMock/pkg/model"
)

func Build(sections []ini.Section) (*model.BehaviorSet, error) {
	bs := &model.BehaviorSet{}

	for _, section := range sections {
		behavior, err := buildBehavior(section, bs)
		if err != nil {
			return nil, err
		}
		for _, property := range section.Properties {
			if err := propagateResponseBehavior(behavior, property); err != nil {
				return nil, err
			}
		}
	}

	return bs, nil
}

func buildBehavior(section ini.Section, bs *model.BehaviorSet) (*model.Behavior, error) {
	b := &model.Behavior{ResponseBehavior: &model.ResponseBehavior{}}
	if section.Name == "default" {
		bs.DefaultBehavior = b.ResponseBehavior
	} else {
		if bs.Behaviors == nil {
			bs.Behaviors = []*model.Behavior{}
		}
		if err := parseBehaviorHeader(b, section.Name, section.LineIndex); err != nil {
			return nil, err
		}
		bs.Behaviors = append(bs.Behaviors, b)
	}
	return b, nil
}

func parseBehaviorHeader(behaviors *model.Behavior, line string, lineIndex uint64) error {
	behaviorHeader := strings.SplitN(strings.Trim(line, "[]"), " ", 2)
	if len(behaviorHeader) != 2 {
		return &MalformedBehaviorHeaderError{Line: line, LineIndex: lineIndex}
	}

	url := strings.TrimSpace(behaviorHeader[1])
	if url == "" || !strings.HasPrefix(url, "/") {
		return &MalformedBehaviorHeaderError{Line: line, LineIndex: lineIndex, Details: Ptr("URL cannot be empty")}
	}

	method, match := model.HttpMethodFromString(behaviorHeader[0])
	if !match {
		return &MalformedBehaviorHeaderError{
			Line:      line,
			LineIndex: lineIndex,
			Details:   Ptr("Invalid HTTP method: " + behaviorHeader[0]),
		}
	}

	behaviors.Method = method
	behaviors.URL = url

	return nil
}

func propagateResponseBehavior(behavior *model.Behavior, property ini.Property) error {
	switch property.Key {
	case "status_code":
		if err := parseStatusCode(behavior.ResponseBehavior, property); err != nil {
			return err
		}
	case "body":
		behavior.ResponseBehavior.Body = Ptr(property.Value)
	case "delay":
		if err := parseDelay(behavior.ResponseBehavior, property); err != nil {
			return err
		}
	case "header":
		if err := parseHeaderAttribute(behavior.ResponseBehavior, property); err != nil {
			return err
		}
	case "redirect":
		behavior.ResponseBehavior.Body = Ptr(property.Value)
	case "sse":
		behavior.ResponseBehavior.SSE = strings.ToLower(property.Value) == "true"
	case "repeat":
		if err := praseRepeat(behavior, property); err != nil {
			return err
		}
	default:
		if strings.HasPrefix(property.Key, "cookie") {
			if err := parseCookie(behavior.ResponseBehavior, property); err != nil {
				return err
			}
			return nil
		}

		return &MalformedPropertyError{
			LineIndex: property.LineIndex,
			Line:      property.Key + "=" + property.Value,
			Details:   Ptr("Unknown property: " + property.Key),
		}
	}

	return nil
}

func parseStatusCode(responseBehavior *model.ResponseBehavior, property ini.Property) error {
	statusCode, err := strconv.Atoi(property.Value)
	if err != nil || statusCode < 100 || statusCode > 599 {
		return &MalformedPropertyError{
			LineIndex: property.LineIndex,
			Line:      property.Key + "=" + property.Value,
			Details:   Ptr("Invalid status code, must be an integer between 100 and 599"),
		}
	}
	responseBehavior.StatusCode = Ptr(uint16(statusCode))
	return nil
}

func parseDelay(responseBehavior *model.ResponseBehavior, property ini.Property) error {
	delayDuration, err := time.ParseDuration(property.Value)
	if err != nil || delayDuration < 0 {
		return &MalformedPropertyError{
			LineIndex: property.LineIndex,
			Line:      property.Key + "=" + property.Value,
			Details:   Ptr("Invalid delay, must be a non-negative duration"),
		}
	}
	responseBehavior.Delay = &delayDuration
	return nil
}

func parseHeaderAttribute(responseBehavior *model.ResponseBehavior, property ini.Property) error {
	if responseBehavior.Headers == nil {
		responseBehavior.Headers = make(map[string]string)
	}

	headerParts := strings.SplitN(property.Value, ":", 2)
	if len(headerParts) != 2 {
		return &MalformedPropertyError{
			LineIndex: property.LineIndex,
			Line:      property.Key + "=" + property.Value,
			Details:   Ptr("Invalid header format, expected 'Key: Value'"),
		}
	}

	key := strings.TrimSpace(headerParts[0])
	value := strings.TrimSpace(headerParts[1])
	if key == "" || value == "" {
		return &MalformedPropertyError{
			LineIndex: property.LineIndex,
			Line:      property.Key + "=" + property.Value,
			Details:   Ptr("Header key and value cannot be empty"),
		}
	}

	responseBehavior.Headers[key] = value
	return nil
}

func parseCookie(responseBehavior *model.ResponseBehavior, property ini.Property) error {
	var cookie *http.Cookie
	if property.Key == "cookie.name" {
		cookie = &http.Cookie{}
		responseBehavior.Cookies = append(responseBehavior.Cookies, cookie)
		cookie.Name = property.Value
		return nil
	} else if len(responseBehavior.Cookies) == 0 {
		return &MalformedPropertyError{
			LineIndex: property.LineIndex,
			Line:      property.Key + "=" + property.Value,
			Details:   Ptr("cookie.name must be set before other cookie properties"),
		}
	} else {
		cookie = responseBehavior.Cookies[len(responseBehavior.Cookies)-1]
	}

	switch property.Key {
	case "cookie.value":
		cookie.Value = property.Value
	case "cookie.path":
		cookie.Path = property.Value
	case "cookie.domain":
		cookie.Domain = property.Value
	case "cookie.expires":
		dur, err := time.ParseDuration(property.Value)
		if err != nil {
			return &MalformedPropertyError{
				LineIndex: property.LineIndex,
				Line:      property.Key + "=" + property.Value,
				Details:   Ptr("Invalid cookie.expires duration"),
			}
		}
		cookie.Expires = time.Now().Add(dur)
	case "cookie.raw_expires":
		t, err := time.Parse(time.RFC3339, property.Value)
		if err != nil {
			return &MalformedPropertyError{
				LineIndex: property.LineIndex,
				Line:      property.Key + "=" + property.Value,
				Details:   Ptr("Invalid cookie.raw_expires RFC3339 time"),
			}
		}
		cookie.RawExpires = property.Value
		cookie.Expires = t
	case "cookie.max_age":
		maxAge, err := strconv.Atoi(property.Value)
		if err != nil {
			return &MalformedPropertyError{
				LineIndex: property.LineIndex,
				Line:      property.Key + "=" + property.Value,
				Details:   Ptr("Invalid cookie.max_age integer"),
			}
		}
		cookie.MaxAge = maxAge
	case "cookie.secure":
		cookie.Secure = strings.EqualFold(property.Value, "true")
	case "cookie.http_only":
		cookie.HttpOnly = strings.EqualFold(property.Value, "true")
	case "cookie.same_site":
		switch strings.ToLower(property.Value) {
		case "lax":
			cookie.SameSite = http.SameSiteLaxMode
		case "strict":
			cookie.SameSite = http.SameSiteStrictMode
		case "none":
			cookie.SameSite = http.SameSiteNoneMode
		default:
			return &MalformedPropertyError{
				LineIndex: property.LineIndex,
				Line:      property.Key + "=" + property.Value,
				Details:   Ptr("Invalid cookie.same_site value"),
			}
		}
	case "cookie.partitioned":
		cookie.Partitioned = strings.EqualFold(property.Value, "true")
	default:
		return &MalformedPropertyError{
			LineIndex: property.LineIndex,
			Line:      property.Key + "=" + property.Value,
			Details:   Ptr("Unknown cookie property: " + property.Value),
		}
	}

	return nil
}

func praseRepeat(behavior *model.Behavior, property ini.Property) error {
	repeat, err := strconv.Atoi(property.Value)
	if err != nil || repeat < 0 {
		return &MalformedPropertyError{
			LineIndex: property.LineIndex,
			Line:      property.Key + "=" + property.Value,
			Details:   Ptr("Invalid repeat value, must be a non-negative integer"),
		}
	}
	repeatUint := uint(repeat)
	behavior.Repeat = &repeatUint
	return nil
}
