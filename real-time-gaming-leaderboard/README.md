# Real-time Gaming Leaderboard System

A high-performance, scalable leaderboard system designed to handle 500 million daily active users (DAU) using Go and Redis with sharding.

## Features

- Real-time score updates and rankings
- Redis sharding for horizontal scalability
- Connection pooling for optimal performance
- Rate limiting to prevent abuse
- Comprehensive error handling and logging
- RESTful API endpoints
- Unit tests for core functionality

## Architecture

The system uses Redis sorted sets for leaderboard data storage, with the following key components:

- Sharded Redis instances for data distribution
- Connection pooling for efficient resource utilization
- Concurrent-safe operations
- Low-latency response times
- Rate limiting middleware

## Prerequisites

- Go 1.25 or higher
- Redis 8.4 or higher (multiple instances for sharding)

## Configuration

The system uses environment variables for configuration. Create a `.env` file with the following variables:

```env
REDIS_SHARDS=3
REDIS_HOSTS=localhost:6379,localhost:6380,localhost:6381
API_PORT=8080
RATE_LIMIT=100
LOG_LEVEL=info
```

## Installation

1. Install dependencies:

```bash
go mod download
```

2. Run the application:

```bash
go run cmd/server/main.go
```

## API Endpoints

### Leaderboard Operations

- `POST /api/v1/scores` - Add or update a player's score
- `GET /api/v1/rank/:player_id` - Get a player's rank
- `GET /api/v1/top/:limit` - Get top N players
- `GET /api/v1/range/:start/:stop` - Get players within a rank range

## Performance Considerations

- The system is designed to handle 500M DAU
- Redis sharding provides horizontal scalability
- Connection pooling optimizes resource usage
- Rate limiting prevents system overload
- Concurrent operations are handled safely
