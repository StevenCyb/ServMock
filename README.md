# TODO
* Handler
* Server
* Conditional match
* Wildcard path
* Replacement from path, query or req body to response


```ini
; Default
status_code = 404

[GET /greeting]
body = Hello, World!
header= Content-Type: text/plain
header = Content-Length: 13

[GET /is_alive]
status_code = 500
repeat = 1

[GET /is_alive]
status_code = 200
body = Service is alive!

[GET /all]
status_code = 200
body = Service is alive!
repeat = 1
delay = 3s
redirect = http://example.com
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