# Echo HTTP Server

A simple Go HTTP server that echoes back request details, useful for testing and load balancing scenarios.

## Features

- Echoes request method, path, headers, and body
- Configurable response status code
- Simulated latency for load testing
- Response size padding
- Docker support
- Comprehensive logging

## Environment Variables

- `PORT`: Server port (default: 8080)
- `DELAY`: Response delay in milliseconds (e.g., `DELAY=100` for 100ms delay)
- `STATUS_CODE`: HTTP response status code (default: 200)
- `RESPONSE_SIZE`: Minimum response size in bytes (pads with 'x' if needed)

## Usage

### Local Development

```bash
go run main.go
```

### Docker

```bash
docker run -p 8080:8080 felpasl/echo-http
```

With environment variables:

```bash
docker run -p 8080:8080 -e DELAY=500 -e STATUS_CODE=201 felpasl/echo-http
```

### Testing

```bash
go test -v
```

## API Response Format

The server responds with the request details in plain text:

```
Method: GET

Path: /test

Headers:
User-Agent: curl/7.68.0
Accept: */*

Body: (only included if request has body)
```

## Building

### Go Binary

```bash
go build -o echo-http .
```

### Docker Image

```bash
docker build -t felpasl/echo-http .
```

## Deployment

The project includes GitHub Actions for automated building and pushing to GitHub Packages.

## Logging

Each request logs:
```
Method: GET Path:/ status: 200 time: 0.026ms 379µs
```

## Examples

### Simple GET request
```bash
curl http://localhost:8080/
```

### POST with body
```bash
curl -X POST -d "test data" http://localhost:8080/api
```

### With delay
```bash
DELAY=1000 go run main.go
curl http://localhost:8080/
# Response delayed by 1 second
```

### Custom status
```bash
STATUS_CODE=404 go run main.go
curl http://localhost:8080/
# Returns 404 status
```