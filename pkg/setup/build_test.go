package parser

import (
	"net/http"
	"testing"
	"time"

	"github.com/StevenCyb/ServMock/pkg/ini"
	"github.com/stretchr/testify/assert"
)

func TestBuild_DefaultAndCustomBehavior(t *testing.T) {
	sections := []ini.Section{
		{Name: "default", Properties: []ini.Property{{Key: "status_code", Value: "200"}}},
		{Name: "GET /foo", LineIndex: 1, Properties: []ini.Property{{Key: "body", Value: "bar"}}},
	}
	bs, err := Build(sections)
	assert.NoError(t, err)
	assert.NotNil(t, bs.DefaultBehavior)
	assert.Equal(t, uint16(200), *bs.DefaultBehavior.StatusCode)
	assert.Len(t, bs.Behaviors, 1)
	assert.Equal(t, "bar", *bs.Behaviors[0].ResponseBehavior.Body)
}

func TestBuild_MalformedBehaviorHeaderError(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET", LineIndex: 2, Properties: []ini.Property{}},
	}
	bs, err := Build(sections)
	assert.Nil(t, bs)
	assert.Error(t, err)
	assert.IsType(t, &MalformedBehaviorHeaderError{}, err)
}

func TestBuild_MalformedPropertyError(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /bar", LineIndex: 3, Properties: []ini.Property{{Key: "unknown", Value: "val", LineIndex: 3}}},
	}
	bs, err := Build(sections)
	assert.Nil(t, bs)
	assert.Error(t, err)
	assert.IsType(t, &MalformedPropertyError{}, err)
}

func TestBuild_StatusCodeValidation(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /baz", LineIndex: 4, Properties: []ini.Property{{Key: "status_code", Value: "999", LineIndex: 4}}},
	}
	bs, err := Build(sections)
	assert.Nil(t, bs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid status code")
}

func TestBuild_DelayValidation(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /delay", LineIndex: 5, Properties: []ini.Property{{Key: "delay", Value: "-1s", LineIndex: 5}}},
	}
	bs, err := Build(sections)
	assert.Nil(t, bs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid delay")
}

func TestBuild_HeaderValidation(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /header", LineIndex: 6, Properties: []ini.Property{{Key: "header", Value: "X-Test", LineIndex: 6}}},
	}
	bs, err := Build(sections)
	assert.Nil(t, bs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid header format")
}

func TestBuild_HeaderSuccess(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /header", LineIndex: 7, Properties: []ini.Property{{Key: "header", Value: "X-Test: Value", LineIndex: 7}}},
	}
	bs, err := Build(sections)
	assert.NoError(t, err)
	assert.Equal(t, "Value", bs.Behaviors[0].ResponseBehavior.Headers["X-Test"])
}

func TestBuild_HeaderWithEmptyValue(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /header", LineIndex: 7, Properties: []ini.Property{{Key: "header", Value: ":", LineIndex: 7}}},
	}
	_, err := Build(sections)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Header key and value cannot be empty")
}

func TestBuild_CookieProperties(t *testing.T) {
	dur := "1h"
	expires := time.Now().Add(time.Hour).Format(time.RFC3339)
	sections := []ini.Section{
		{Name: "GET /cookie", LineIndex: 8, Properties: []ini.Property{
			{Key: "cookie.name", Value: "foo"},
			{Key: "cookie.value", Value: "bar"},
			{Key: "cookie.path", Value: "/"},
			{Key: "cookie.domain", Value: "example.com"},
			{Key: "cookie.expires", Value: dur},
			{Key: "cookie.raw_expires", Value: expires},
			{Key: "cookie.max_age", Value: "10"},
			{Key: "cookie.secure", Value: "true"},
			{Key: "cookie.http_only", Value: "true"},
			{Key: "cookie.same_site", Value: "lax"},
			{Key: "cookie.partitioned", Value: "true"},
		}},
	}
	bs, err := Build(sections)
	assert.NoError(t, err)
	assert.Len(t, bs.Behaviors[0].ResponseBehavior.Cookies, 1)
	cookie := bs.Behaviors[0].ResponseBehavior.Cookies[0]
	assert.Equal(t, "foo", cookie.Name)
	assert.Equal(t, "bar", cookie.Value)
	assert.Equal(t, "/", cookie.Path)
	assert.Equal(t, "example.com", cookie.Domain)
	assert.Equal(t, 10, cookie.MaxAge)
	assert.True(t, cookie.Secure)
	assert.True(t, cookie.HttpOnly)
	assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	assert.True(t, cookie.Partitioned)
}

func TestBuild_TwoCookies(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /cookie", LineIndex: 8, Properties: []ini.Property{
			{Key: "cookie.name", Value: "foo"},
			{Key: "cookie.value", Value: "/foo"},
			{Key: "cookie.name", Value: "bar"},
			{Key: "cookie.value", Value: "/bar"},
		}},
	}
	bs, err := Build(sections)
	assert.NoError(t, err)
	assert.Len(t, bs.Behaviors[0].ResponseBehavior.Cookies, 2)
	cookie1 := bs.Behaviors[0].ResponseBehavior.Cookies[0]
	cookie2 := bs.Behaviors[0].ResponseBehavior.Cookies[1]
	assert.Equal(t, "foo", cookie1.Name)
	assert.Equal(t, "/foo", cookie1.Value)
	assert.Equal(t, "bar", cookie2.Name)
	assert.Equal(t, "/bar", cookie2.Value)
}

func TestBuild_NoNameCookies(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /cookie", LineIndex: 8, Properties: []ini.Property{
			{Key: "cookie.value", Value: "/foo"},
		}},
	}
	_, err := Build(sections)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cookie.name must be set before other cookie properties")
}

func TestBuild_RepeatProperty(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /repeat", LineIndex: 9, Properties: []ini.Property{{Key: "repeat", Value: "5", LineIndex: 9}}},
	}
	bs, err := Build(sections)
	assert.NoError(t, err)
	assert.NotNil(t, bs.Behaviors[0].Repeat)
	assert.Equal(t, uint(5), *bs.Behaviors[0].Repeat)
}

func TestBuild_RepeatPropertyInvalid(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /repeat", LineIndex: 10, Properties: []ini.Property{{Key: "repeat", Value: "-1", LineIndex: 10}}},
	}
	bs, err := Build(sections)
	assert.Nil(t, bs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid repeat value")
}

func TestBuild_CookieUnknownProperty(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /cookie", LineIndex: 11, Properties: []ini.Property{{Key: "cookie.name", Value: "foo"}, {Key: "cookie.unknown", Value: "val", LineIndex: 11}}},
	}
	bs, err := Build(sections)
	assert.Nil(t, bs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unknown cookie property")
}

func TestBuild_CookieInvalidMaxAge(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /cookie", LineIndex: 12, Properties: []ini.Property{{Key: "cookie.name", Value: "foo"}, {Key: "cookie.max_age", Value: "notanint", LineIndex: 12}}},
	}
	bs, err := Build(sections)
	assert.Nil(t, bs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid cookie.max_age integer")
}

func TestBuild_CookieInvalidRawExpires(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /cookie", LineIndex: 13, Properties: []ini.Property{{Key: "cookie.name", Value: "foo"}, {Key: "cookie.raw_expires", Value: "notadate", LineIndex: 13}}},
	}
	bs, err := Build(sections)
	assert.Nil(t, bs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid cookie.raw_expires RFC3339 time")
}

func TestBuild_CookieInvalidExpires(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /cookie", LineIndex: 14, Properties: []ini.Property{{Key: "cookie.name", Value: "foo"}, {Key: "cookie.expires", Value: "notaduration", LineIndex: 14}}},
	}
	bs, err := Build(sections)
	assert.Nil(t, bs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid cookie.expires duration")
}

func TestBuild_CookieSameSiteStrict(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /cookie", LineIndex: 15, Properties: []ini.Property{{Key: "cookie.name", Value: "foo"}, {Key: "cookie.same_site", Value: "strict", LineIndex: 15}}},
	}
	bs, err := Build(sections)
	assert.NoError(t, err)
	cookie := bs.Behaviors[0].ResponseBehavior.Cookies[0]
	assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
}

func TestBuild_CookieSameSiteNone(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /cookie", LineIndex: 16, Properties: []ini.Property{{Key: "cookie.name", Value: "foo"}, {Key: "cookie.same_site", Value: "none", LineIndex: 16}}},
	}
	bs, err := Build(sections)
	assert.NoError(t, err)
	cookie := bs.Behaviors[0].ResponseBehavior.Cookies[0]
	assert.Equal(t, http.SameSiteNoneMode, cookie.SameSite)
}

func TestBuild_CookieSameSiteInvalid(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /cookie", LineIndex: 17, Properties: []ini.Property{{Key: "cookie.name", Value: "foo"}, {Key: "cookie.same_site", Value: "invalid", LineIndex: 17}}},
	}
	bs, err := Build(sections)
	assert.Nil(t, bs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid cookie.same_site value")
}

func TestBuild_RedirectProperty(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /redirect", LineIndex: 18, Properties: []ini.Property{{Key: "redirect", Value: "/new", LineIndex: 18}}},
	}
	bs, err := Build(sections)
	assert.NoError(t, err)
	assert.Equal(t, "/new", *bs.Behaviors[0].ResponseBehavior.Body)
}

func TestBuild_StreamPropertyTrue(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /stream", LineIndex: 19, Properties: []ini.Property{{Key: "stream", Value: "true", LineIndex: 19}}},
	}
	bs, err := Build(sections)
	assert.NoError(t, err)
	assert.True(t, bs.Behaviors[0].ResponseBehavior.Stream)
}

func TestBuild_StreamPropertyFalse(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET /stream", LineIndex: 20, Properties: []ini.Property{{Key: "stream", Value: "false", LineIndex: 20}}},
	}
	bs, err := Build(sections)
	assert.NoError(t, err)
	assert.False(t, bs.Behaviors[0].ResponseBehavior.Stream)
}

func TestBuild_InvalidHttpMethod(t *testing.T) {
	sections := []ini.Section{
		{Name: "FOO /bar", LineIndex: 21, Properties: []ini.Property{}},
	}
	bs, err := Build(sections)
	assert.Nil(t, bs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid HTTP method")
}

func TestBuild_InvalidUrl(t *testing.T) {
	sections := []ini.Section{
		{Name: "GET bar", LineIndex: 22, Properties: []ini.Property{}},
	}
	bs, err := Build(sections)
	assert.Nil(t, bs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "URL cannot be empty")
}
