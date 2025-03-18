# FHIR API Gateway & Integration Service

## Overview
- Building a **FHIR R4 API Gateway** in Go.  
- It includes **API key security** and more.

## Features
- FHIR R4 Endpoints (Patient) with easy extension
- API Key middleware
- GORM-based PostgreSQL storage
- Structured logging
- Prometheus monitoring
- Cached Patient data using Valkey for faster responses and reduced database load

## Quick Start
1. Run with Docker Compose
```bash
docker-compose up --build
```
Access the service at http://localhost:8080. Set header X-API-Key for requests.

2. Access the API

Create new Patient
```shell
curl -X POST http://localhost:8080/fhir/r4/Patient -H "Content-Type: application/json" -H "X-API-KEY: test-api-key-1" \
-d '{
    "resourceType": "Patient",
    "id": "patient-001",
    "name": [{
        "family": "Smith",
        "given": ["John"]
    }],
    "gender": "male",
    "birthDate": "1980-05-15",
    "address": [{
        "use": "home",
        "line": ["123 Main St"],
        "city": "Boston",
        "state": "MA",
        "postalCode": "02108",
        "country": "USA"
    }],
    "active": true
}'
```

Get Patient created
```shell
curl -H "X-API-KEY: test-api-key-1" http://localhost:8080/fhir/r4/Patient/patient-001
```

Get metrics
```shell
curl http://localhost:8080/metrics
```

3. Local Development
- Install Go 1.24 or higher
- Install dependencies
```shell
go get github.com/gin-gonic/gin
go get github.com/samply/golang-fhir-models/fhir-models
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/swaggo/files
go get github.com/sirupsen/logrus
go get github.com/prometheus/client_golang/prometheus/promhttp
go get github.com/valkey-io/valkey-go
```
- PostgreSQL
```shell
docker run -d --name fhir-postgres -e POSTGRES_USER=fhiruser -e POSTGRES_PASSWORD=fhirpassword -e POSTGRES_DB=fhirdb -p 5432:5432 postgres:17.4
docker run -d --name fhir-valkey -p 6379:6379 valkey/valkey:8.0
```
Import initial database schema: `migrations/000001_create_tables.up.sql`
- Environment variables
```bash
export DB_URL=postgres://fhiruser:fhirpassword@localhost:5432/fhirdb?sslmode=disable
export CACHE_URL=localhost:6379
```
- make build && make run (or go run ./cmd/main.go)
