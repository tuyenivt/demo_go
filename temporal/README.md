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

- Start Temporal with docker

```shell
docker network create devnet
docker run -d --network devnet --name temporal --user root -p 7233:7233 -p 8233:8233 -v temporal-data:/data temporalio/temporal:1.5.1 server start-dev --ip 0.0.0.0 --db-filename /data/temporal.db
```

- Install Go 1.25 or higher
- Install dependencies

```shell
go get -u github.com/gin-gonic/gin
go get -u go.temporal.io/sdk
```

3. Testing APIs

### Health check

```shell
curl http://localhost:8080/health
```
