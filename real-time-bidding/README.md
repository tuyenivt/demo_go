# Real-time Bidding (RTB) System

A high-performance Real-time Bidding system implemented in Go with Redis for ultra-fast caching and concurrent bid processing.

## Features

- Sub-50ms latency for bid requests
- Handles 100,000 bid requests per second
- Redis-based caching for user profiles and campaign data
- Concurrent bid processing with goroutines
- Graceful shutdown and resource cleanup
- Structured logging
- Environment-based configuration

## Prerequisites

- Go 1.24 or later
- Redis 7.0 or later

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Configure environment variables:
```bash
cp .env.example .env
# Edit .env with your Redis configuration
```

3. Start Redis:
```bash
# Using Docker
docker run -d -p 6379:6379 redis:7

# Or using your local Redis installation
redis-server
```

## Running with Docker Compose

1. Build and start the services:
```bash
docker-compose up --build
```

2. To run in detached mode:
```bash
docker-compose up -d
```

3. To stop the services:
```bash
docker-compose down
```

This setup makes it easy to run the entire application stack locally with a single command while maintaining proper isolation between services.

## Running the Application

1. Start the server:
```bash
go run main.go
```

2. The server will start on port 8080.

## API Endpoints

### POST /bid

Process a bid request.

Request body:
```json
{
    "request_id": "req_123",
    "user_id": "user_456",
    "timestamp": "2024-03-20T10:00:00Z",
    "device": {
        "device_type": "mobile",
        "os": "ios",
        "browser": "safari",
        "ip": "192.168.1.1"
    },
    "user_segments": ["premium", "gaming"],
    "context": {
        "page_url": "https://example.com",
        "ad_size": "300x250"
    }
}
```

Response:
```json
{
    "request_id": "req_123",
    "campaign_id": "campaign1",
    "price": 2.50,
    "ad_id": "ad_campaign1",
    "ad_markup": "<div>Sample Ad Content</div>",
    "latency_ms": 45
}
```

## Redis Key Patterns

- User Segments: `user:segments:<userID>`
- Campaign Data: `campaign:<campaignID>`
- Campaign Budget: `campaign:budget:<campaignID>`
- User Profile: `user:profile:<userID>`

## Performance Considerations

1. Redis Configuration:
   - Use Redis cluster for high availability
   - Configure appropriate memory limits
   - Enable persistence if needed

2. Application Tuning:
   - Adjust connection pool size based on load
   - Monitor goroutine usage
   - Use appropriate TTL values for cached data

3. Monitoring:
   - Monitor Redis memory usage
   - Track bid request latency
   - Monitor error rates
