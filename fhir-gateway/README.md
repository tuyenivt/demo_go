# FHIR API Gateway & Integration Service

## Overview
- Building a **FHIR R4 API Gateway** in Go.  
- It includes **API key security** and more.

## Features
- FHIR R4 Endpoints (Patient, Observation) with easy extension
- API Key middleware
- GORM-based MySQL storage
- HIPAA compliance: encrypt sensitive data at rest using MySQL's encryption features and log access to sensitive resources, capturing events like unauthorized access attempts
- Docker & K8s deployment
- CI/CD example with GitHub Actions
- Structured logging (zap)

## Docker Commands
```shell
docker run -d --name fhir-mysql -e MYSQL_USER=fhiruser -e MYSQL_PASSWORD=fhirpassword -e MYSQL_DATABASE=fhirdb -p 3306:3306 mysql:8.4
```

## Quick Start
1. Environment Variables
```bash
API_SERVER_PORT=8080
API_KEY=supersecret
MYSQL_USER=fhiruser
MYSQL_PASSWORD=fhirpassword
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DATABASE=fhirdb
```

2. Run with Docker Compose
```bash
docker-compose up --build
```
Access the service at http://localhost:8080. Set header X-API-Key: supersecret for requests.

