# Rune

In-memory key-value store over HTTP. Set, get, and delete string keys with configurable size limits, structured JSON logging, and Prometheus metrics.

## Build

```bash
go build -o rune ./cmd/server
```

## Run

```bash
go run ./cmd/server
# or
./rune
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | `8080` | HTTP listen port |
| `-max-key-size` | `256` | Max key length (bytes) |
| `-max-value-size` | `1048576` | Max value length (bytes, 1 MiB) |
| `-log-level` | `debug` | `debug`, `info`, `warn`, or `error` |

Example:

```bash
./rune -port 9090 -log-level info
```

## API

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/set` | Form fields `key`, `value` |
| `GET` | `/get?key=...` | Returns value as plain text |
| `DELETE` | `/delete?key=...` | Removes key |
| `GET` | `/metrics` | Prometheus scrape endpoint |

```bash
curl -X POST http://localhost:8080/set -d "key=foo&value=bar"
curl "http://localhost:8080/get?key=foo"
curl -X DELETE "http://localhost:8080/delete?key=foo"
```

## Test

```bash
go test ./...
```
