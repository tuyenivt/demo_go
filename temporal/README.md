# Temporal

## Overview

- A temporal Go demo project.

## Quick Start

1. Run server

```bash
go run main.go
```

Access the service at `http://localhost:8080/health`

2. Local Development

- Install Go 1.25 or higher
- Install dependencies

```shell
go get -u github.com/gin-gonic/gin
```

3. Testing APIs

### Health check

```shell
curl http://localhost:8080/health
```
