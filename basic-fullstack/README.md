# Basic Fullstack

## Overview

- A fullstack Go practical project.

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
docker run -d --network devnet --name movie-postgres -e POSTGRES_USER=username -e POSTGRES_PASSWORD=password -e POSTGRES_DB=moviedb -p 5432:5432 postgres:17.6-alpine
```

- Load an initial database using a dump file in PostgreSQL: `./migrations/moviedb_backup.sql`

```shell
cd migrations
PGPASSWORD="password" psql -h localhost -p 5432 -U username -d moviedb -f moviedb_backup.sql
```

- Rename `.env.example` to `.env`, check `DATABASE_URL` and change if needed

- Install Go 1.25 or higher
- Install dependencies

```shell
go get github.com/joho/godotenv
go get github.com/lib/pq
```

3. Testing APIs

### Health check

```shell
curl http://localhost:8080/health
```

### Get top movies

```shell
curl http://localhost:8080/api/movies/top
```

### Search movies by name

```shell
curl 'http://localhost:8080/api/movies/search?q=Lion&order=&genre='
```

### Get a movie by ID

```shell
curl http://localhost:8080/api/movies/1
```

### Get all genres

```shell
curl http://localhost:8080/api/genres
```
