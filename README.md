# HTTP Helper

A versatile HTTP server for development and testing with built-in network simulation capabilities.

---

[![Go](https://github.com/hurricanerix/http-helper/actions/workflows/go.yml/badge.svg)](https://github.com/hurricanerix/http-helper/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/hurricanerix/http-helper)](https://goreportcard.com/report/github.com/hurricanerix/http-helper)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview

HTTP Helper is a feature-rich HTTP server designed for local development and testing. It provides Python 
SimpleHTTPServer-like functionality with additional middleware for simulating real-world network conditions, 
making it perfect for testing applications under various scenarios.

## Features

- **Static file serving** with directory listing
- **Network simulation** - bandwidth throttling and latency injection
- **Configurable middleware pipeline** - mix and match components
- **Zero dependencies** - single binary distribution
- **Cross-platform** - Windows, macOS (Intel/ARM), and Linux
- **Security features** - directory traversal protection
- **Request logging** with structured output
- **ETag generation** for caching tests
- **CORS support** with flexible configuration

## Installation

### Using Go

```bash
go install github.com/hurricanerix/http-helper/cmd/hs@latest
```

### From Source

```bash
git clone https://github.com/hurricanerix/http-helper.git
cd http-helper
make build
```

### Pre-built Binaries

Download the latest release from the [releases page](https://github.com/hurricanerix/http-helper/releases).


## Quick Start

```bash
# Serve current directory on port 8000
hs

# Serve specific directory on custom port
hs -d ./public -port 8080

# Serve on all interfaces
hs -bind 0.0.0.0

# Simulate slow network (1MB/s bandwidth)
HH_BANDWIDTH_BPS=1000000 hs

# Simulate high latency (500ms TTFB)
HH_TIME_TO_FIRST_BYTE=500 hs
```

## Configuration

HTTP Helper can be configured via environment variables. Configuration is loaded from:
1. `.env` file in the current directory
2. `~/.config/http_helper.env` (global config)
3. Environment variables (highest priority)

See `example.env` for more details.

## Middleware Pipeline

The configuration of its middleware pipeline. but in most cases the default should be used.

## Usage Examples

### Local Development Server
```bash
# Serve your React/Vue/Angular build
hs -d ./dist -port 3000
```

### Network Condition Testing
```bash
# Simulate 3G network (384kbps)
HH_BANDWIDTH_BPS=48000 HH_TIME_TO_FIRST_BYTE=600 hs

# Simulate unstable connection with jitter
HH_BANDWIDTH_BPS=1000000 HH_BANDWIDTH_JITTER=500000 \
HH_TIME_TO_FIRST_BYTE=200 HH_TIME_TO_FIRST_BYTE_JITTER=150 hs
```

### CORS Testing
```bash
# Allow specific origins
HH_CORS_ALLOWED_ORIGINS="http://localhost:3000,https://myapp.com" \
HH_CORS_ALLOWED_METHODS="GET,POST,PUT,DELETE" hs
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](docs/CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by Python's `http.server` module