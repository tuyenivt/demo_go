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

- Create database tables by using sql scripts in `migrations` folder (skipped database migrations with Goose).

- Install Go 1.25 or higher
- Install dependencies

```shell
go get -u github.com/go-chi/chi/v5
go get github.com/jackc/pgx/v4/stdlib
```

3. Testing APIs

### Create workout

```shell
curl -X POST http://localhost:8080/workouts \
     -H "Content-Type: application/json" \
     -d '{
           "title": "Title 1",
           "description": "Description 1",
           "duration_minutes": 60,
           "calories_burned": 500,
           "entries": [
                {
                    "exercise_name": "Exercise Name 1",
                    "sets": 1,
                    "reps": 2,
                    "weight": 15.1,
                    "notes": "Note 1",
                    "order_index": 1
                },
                {
                    "exercise_name": "Exercise Name 2",
                    "sets": 2,
                    "reps": 3,
                    "weight": 15.2,
                    "notes": "Note 2",
                    "order_index": 2
                }
           ]
         }'
```

### Get workout

```shell
curl http://localhost:8080/workouts/1
```

### Update workout

```shell
curl -X PUT http://localhost:8080/workouts/1 \
     -H "Content-Type: application/json" \
     -d '{
           "title": "Title 1.1",
           "description": "Description 1.1",
           "duration_minutes": 61,
           "calories_burned": 501,
           "entries": [
                {
                    "exercise_name": "Exercise Name 1.1",
                    "sets": 1,
                    "reps": 2,
                    "weight": 15.11,
                    "notes": "Note 1.1",
                    "order_index": 1
                },
                {
                    "exercise_name": "Exercise Name 2.1",
                    "sets": 2,
                    "reps": 3,
                    "weight": 15.21,
                    "notes": "Note 2.1",
                    "order_index": 2
                }
           ]
         }'
```

### Delete workout

```shell
curl -X DELETE http://localhost:8080/workouts/1
```

### Create user

```shell
curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{
           "username": "username1",
           "email": "username1@example.com",
           "password": "password1",
           "bio": "bio1"
         }'
```

### Get user

```shell
curl http://localhost:8080/users/username1
```
