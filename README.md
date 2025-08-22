# Caching Proxy Server

This project is a simple caching proxy server written in Go. It acts as a reverse proxy between clients and an origin server, caching responses in memory to improve performance and reduce load on the origin server.

## Features
- Reverse proxy for HTTP requests
- In-memory caching of responses
- Configurable origin server and port
- Option to clear cache on startup

## Usage

### Build and Run
```sh
go build -o caching-proxy-server
./caching-proxy-server --port=3002 --origin=http://example.com
```

### Command Line Flags
- `--port` (required): Port to run the proxy server on
- `--origin` (required): Origin server URL (e.g., http://example.com)
- `--clear-cache`: Clear in-memory cache on start

### Example
```sh
./caching-proxy-server --port=3002 --origin=https://jsonplaceholder.typicode.com
```

## Project Source
This project is based on the [Caching Server project from roadmap.sh](https://roadmap.sh/projects/caching-server).

## License
MIT
