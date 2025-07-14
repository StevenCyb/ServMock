# ServMock
ServMock is a mock server for testing HTTP clients. It allows you to define responses for specific HTTP requests, making it easier to test your applications without relying on external services and without implementing and maintaining mock API services.

*This project is enough for me but far from production ready.*

## Usage

### Configuration

You can define your mock responses in an INI file format. The following is an example of how to set up a mock server:

```ini
; Default response (supports everything expect "repeat")
status_code = 404

; Define a behavior for a specific HTTP method and path.
; The status code is 200 on a match.
[GET /greeting]
body = Hello, World!
header= Content-Type: text/plain
header = Content-Length: 13

[OPTION /options]
; Define header for this response.
header= Content-Type: text/plain
; Overwrite the default status code (200 OK).
status_code = 201
; Define how often this behavior should be repeated.
; After N times, the next match or default will be used.
repeat = 3
; Body of the response
body = Hello, World!
; Add some delay if needed
delay = 3s
; Respond with a redirect (should also have a matching status code).
redirect = http://example.com
; Respond with a Server-Sent Event (SSE)
; Body will be split by new lines and sent as events.
sse = false
; Response cookies can also bet set with the following properties.
; `cookie.name` must be the first one since it indicates a start of a cookie.
cookie.name = username
cookie.value = steven
cookie.path = /
cookie.domain = example.com
cookie.expires = 24h
cookie.raw_expires = 2023-10-01T00:00:00Z
cookie.max_age = 3600
cookie.secure = true
cookie.http_only = true
cookie.same_site = Lax
```

## Whats next
* Linting
* Documentation
* Pipeline tests
* Docker image build
* Conditional match
* Wildcard path
* Replacement from path, query or req body to response

