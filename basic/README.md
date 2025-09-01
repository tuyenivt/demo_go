# Basic

## Overview

- A basic Go practical project.

## Quick Start

1. Run server

```bash
go run main.go -port 8080
```

Access the service at `http://localhost:8080/health`

2. Local Development

- Prepare database with docker (or local installation)

```shell
docker network create devnet
docker run -d --network devnet --name workout-postgres -e POSTGRES_USER=username -e POSTGRES_PASSWORD=password -e POSTGRES_DB=workoutdb -p 5432:5432 postgres:17.6-alpine
```

- Install Go 1.25 or higher
- Install dependencies

```shell
go get -u github.com/go-chi/chi/v5
go get github.com/jackc/pgx/v4/stdlib
```
