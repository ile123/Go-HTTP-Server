# http-server

A small Go HTTP server built on the standard library (`net/http`). It runs two
listeners at once(can register more), uses a composable middleware stack, and exposes a handful
of endpoints for experimenting with query params, form data, JSON bodies, and
custom error headers.

It uses no third-party dependencies.

## Features

- Two servers started from one process (`:3333` and `:4444`), each tagging, can add more if needed
  requests with the address it received them on
- Chainable middleware stack (logging, versioning header, security headers)
- JSON endpoint with field validation and structured error reporting via
  response headers
- Health check endpoint
- Sensible server timeouts and a TLS config ready to enable

## Requirements

- Go (see the `go` directive in `go.mod` for the version this was built with)

## Getting started

```bash
# clone the repo, then from the project root:
go run .
```

Or build a binary:

```bash
go build -o http-server .
./http-server
```

The servers start on:

| Name       | Address          |
| ---------- | ---------------- |
| server one | `localhost:3333` |
| server two | `localhost:4444` |

Stop with `Ctrl+C`.

## Endpoints

| Method | Path           | Description                                                  |
| ------ | -------------- | ------------------------------------------------------------ |
| any    | `/`            | Returns a greeting; logs query params and request body       |
| any    | `/hello`       | Returns `Hello, HTTP!`; logs query params and form values    |
| POST   | `/api/message` | Accepts a JSON message, validates it, and echoes it back     |
| any    | `/health`      | Returns `{"status":"healthy"}`                               |

### Examples

```bash
# Root
curl -i "localhost:3333/?name=gopher"

# Hello, with query params and form values
curl -i "localhost:3333/hello?lang=go" -d "user=ada"

# Health check
curl -i localhost:4444/health
```

### `POST /api/message`

Request body:

```json
{ "name": "Ada", "body": "Hello there" }
```

```bash
curl -i -X POST localhost:3333/api/message \
  -H "Content-Type: application/json" \
  -d '{"name":"Ada","body":"Hello there"}'
```

Success (`200 OK`):

```json
{"body":"Hello there","message":"Received message from Ada","status":"success"}
```

Error cases:

| Situation                     | Status | Details                                              |
| ----------------------------- | ------ | ---------------------------------------------------- |
| Method is not `POST`          | `405`  | `Allow: POST` header                                 |
| Body is not valid JSON        | `400`  | `X-Invalid-Value: body` header                       |
| `name` and/or `body` missing  | `400`  | One `X-Missing-Field` header per missing field       |

Validation errors are reported in response **headers** rather than the body.
For example, a request missing `name` returns `400` with `X-Missing-Field: name`.
The error kinds are defined in `constants/response.go`.

## Middleware

Middleware is registered in `main.go` and applied in the order listed, so the
first one added is the outermost:

1. **LoggingMiddleware** logs method, URI, remote address, and request duration
2. **HeaderMiddleware** adds `X-App-Version: 1.0` to every response
3. **SecurityHeadersMiddleware** sets `X-Frame-Options`,
   `X-Content-Type-Options`, `X-XSS-Protection`, `Strict-Transport-Security`,
   and `Content-Security-Policy`

To add your own, write a `func(http.Handler) http.Handler` in `middlewares/`
and register it with `stack.Use(...)` in `main.go`.

## Project structure

```
.
├── main.go                  # server setup, routing, middleware wiring
├── constants/
│   ├── api.go               # request/response types (Message)
│   ├── response.go          # error kinds and ResponseError type
│   └── url.go               # context keys
├── handlers/
│   └── handler.go           # HTTP handlers and response helpers
├── helpers/
│   └── url_helper.go        # query/form/body printing utilities
└── middlewares/
    ├── stack.go             # middleware stack (Use / Then)
    └── middleware.go        # logging, header, and security middleware
```

## Configuration

Ports and timeouts are currently set in `main.go`:

| Setting        | Value                         |
| -------------- | ----------------------------- |
| Ports          | `:3333`, `:4444`              |
| Read timeout   | 10s                           |
| Write timeout  | 10s                           |
| Idle timeout   | 120s                          |
| Min TLS version| TLS 1.2 (when TLS is enabled) |

### Enabling HTTPS

The TLS config is already in place. To serve over HTTPS, generate a
certificate and swap the commented line in `run()`:

```bash
openssl req -x509 -newkey rsa:2048 -nodes \
  -keyout key.pem -out cert.pem -days 365 -subj "/CN=localhost"
```

```go
err := s.ListenAndServeTLS("cert.pem", "key.pem")
```

`cert.pem` and `key.pem` are listed in `.gitignore`. Never commit private keys.

## License

GNU General Public License v3.0