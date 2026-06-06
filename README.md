# Health Checker Website Backend

A production-ready, containerized, and highly-scalable health monitoring backend written in **Golang**, adhering to **Clean Architecture** principles.

---

## 1. Project Overview

The Health Checker Website Backend monitors application status and verifies connectivity to backend data layers (MySQL and Redis). It is built with Go's standard library for routing/serving HTTP, uses Structured JSON Logging via `log/slog`, implements graceful shutdowns, handles connection pooling, runs as a secure non-root Docker container, and is orchestrated with Docker Compose.

---

## 2. Architecture Diagram

The project is structured under **Clean Architecture** guidelines:

```mermaid
graph TD
    subgraph Drivers & Frameworks (External)
        API_Client([HTTP Client])
        MySQL_DB[(MySQL 8)]
        Redis_Cache[(Redis 7)]
    end

    subgraph Interface Adapters (Transport & Storage)
        Mux[HTTP Multiplexer / Router]
        Middlewares[Logging / Recovery Middlewares]
        HealthHandler[Health Handler]
        CacheHandler[Cache Handler]
        MySQLClient[MySQL Connection Pool]
        RedisClient[Redis Client Pool]
    end

    subgraph Application Business Logic
        HealthService[Health Service]
        CacheService[Cache Service]
    end

    API_Client --> Middlewares
    Middlewares --> Mux
    Mux --> HealthHandler
    Mux --> CacheHandler

    HealthHandler --> HealthService
    CacheHandler --> CacheService

    HealthService --> MySQLClient
    HealthService --> RedisClient
    CacheService --> RedisClient

    MySQLClient --> MySQL_DB
    RedisClient --> Redis_Cache
```

- **Entities / Core Models**: Internal logic structures.
- **Services (Use Cases)**: Domain logic for ping checks and cache storage.
- **Handlers / Controllers**: Decodes HTTP JSON payloads and maps responses.
- **Database / Infrastructure**: Physical drivers for database connections.

---

## 3. Setup Instructions

### Prerequisites
* [Docker](https://docs.docker.com/get-docker/) (v20+)
* [Docker Compose](https://docs.docker.com/compose/install/) (v2+)
* [Golang](https://golang.org/dl/) (v1.25+ if running locally outside Docker)

### Run with Docker Compose (Recommended)
1. Clone the repository.
2. Initialize the `.env` configuration file (already created at project root).
3. Build and launch all services:
   ```bash
   docker compose up --build
   ```
4. Verify all containers are running successfully:
   ```bash
   docker compose ps
   ```

---

## 4. Environment Variables

Configure application settings by modifying the `.env` file at the root directory:

| Variable | Description | Default |
| :--- | :--- | :--- |
| `APP_NAME` | Name of the API application | `Health Checker` |
| `APP_ENV` | Application environment stage | `development` |
| `APP_PORT` | Listening port for the Golang server | `8000` |
| `DB_HOST` | Database server hostname | `mysql` |
| `DB_PORT` | Database server port | `3306` |
| `DB_NAME` | Application database schema | `health_checker` |
| `DB_USER` | MySQL database user | `health_user` |
| `DB_PASSWORD` | MySQL database password | `health_password` |
| `MYSQL_ROOT_PASSWORD` | MySQL root administrative password | `root_password` |
| `REDIS_HOST` | Redis hostname | `redis` |
| `REDIS_PORT` | Redis port | `6379` |
| `REDIS_PASSWORD` | Redis password | *(Empty)* |
| `REDIS_DB` | Redis database slot | `0` |

---

## 5. API Documentation

### Access URLs
* **API Service**: `http://localhost:8000`
* **phpMyAdmin**: `http://localhost:8080`
* **MySQL Database**: `localhost:3306`
* **Redis Cache**: `localhost:6379`

### Endpoints

#### 1. System Health
Verify api status, mysql, and redis connectivity.

* **URL**: `/health`
* **Method**: `GET`
* **Response Status**: `200 OK` (Healthy) or `503 Service Unavailable` (If MySQL/Redis is down)
* **Response Body**:
  ```json
  {
    "status": "ok",
    "service": "health-checker",
    "mysql": "connected",
    "redis": "connected",
    "timestamp": "2026-06-06T10:00:00Z"
  }
  ```

#### 2. Detailed Health Check
Fetch latency status and detailed uptime records.

* **URL**: `/health/details`
* **Method**: `GET`
* **Response Status**: `200 OK` (Healthy) or `503 Service Unavailable`
* **Response Body**:
  ```json
  {
    "api": "healthy",
    "mysql": {
      "status": "healthy"
    },
    "redis": {
      "status": "healthy"
    },
    "uptime": "24h0m0s"
  }
  ```

#### 3. Save Cache Key
Store a string key-value pair in Redis.

* **URL**: `/cache`
* **Method**: `POST`
* **Headers**: `Content-Type: application/json`
* **Request Body**:
  ```json
  {
    "key": "health",
    "value": "ok"
  }
  ```
* **Response Status**: `200 OK`
* **Response Body**:
  ```json
  {
    "message": "cached successfully"
  }
  ```

#### 4. Retrieve Cache Key
Retrieve a stored key from Redis.

* **URL**: `/cache?key=<name>`
* **Method**: `GET`
* **Response Status**: `200 OK` or `404 Not Found`
* **Response Body**:
  ```json
  {
    "key": "health",
    "value": "ok"
  }
  ```

---

## 6. Docker Commands

Commonly used commands for administration and maintenance:

* **Build and Start Services (in detached mode)**:
  ```bash
  docker compose up -d --build
  ```
* **View Container Statuses**:
  ```bash
  docker compose ps
  ```
* **View Real-Time Logs**:
  ```bash
  docker compose logs -f
  ```
* **Stop Services (keep volume storage)**:
  ```bash
  docker compose down
  ```
* **Stop Services and Clear Volumes (deletes database)**:
  ```bash
  docker compose down -v
  ```
* **Run Database CLI shell**:
  ```bash
  docker compose exec mysql mysql -u health_user -phealth_password health_checker
  ```
* **Run Redis CLI shell**:
  ```bash
  docker compose exec redis redis-cli
  ```

---

## 7. Development Workflow

### Local Development (Without Docker)
If you wish to run the app directly on your host machine:
1. Ensure MySQL and Redis servers are running locally.
2. Edit database settings in `.env` (e.g. `DB_HOST=localhost`, `REDIS_HOST=localhost`).
3. Run the compiler:
   ```bash
   go run cmd/server/main.go
   ```

### Dependency Management
To add packages, execute:
```bash
go get <package_name>
go mod tidy
```

---

## 8. Health Check Guide

* **MySQL Health Check**: Handled by `mysqladmin ping` inside the database container.
* **Redis Health Check**: Handled by `redis-cli ping` inside the cache container.
* **API Service Dependability**: The API service relies on both MySQL and Redis health checks passing before booting up (`depends_on: { condition: service_healthy }`).

---

## 9. Troubleshooting

* **MySQL Connection Refused**:
  * Check if MySQL is still starting up (can take up to 10 seconds on first run to setup tables).
  * Ensure user credentials match between `.env` and `docker-compose.yml`.
* **Redis Connection Refused**:
  * Run `docker compose logs redis` to check for configuration or memory dump failures.
* **Ports Already In Use**:
  * If ports `8000`, `8080`, `3306`, or `6379` are bound by local software, stop them or change ports in the `.env` file before executing `docker compose up`.

---

## 10. Production Deployment Guide

For secure cloud deployment:
1. **Change Default Credentials**: Set secure passwords for MySQL Root, MySQL User, and Redis.
2. **Expose Ports Carefully**: Do not expose ports `3306`, `6379`, or `8080` to the public internet. Only the API service (port `8000` or via reverse proxy) should be exposed.
3. **Use Reverse Proxy**: Deploy Nginx or Caddy to terminate HTTPS traffic and proxy calls to port `8000`.
4. **Environment settings**: Set `APP_ENV=production` inside your configuration manager.
5. **Logs aggregation**: Pipe JSON output logs from standard out to cloud logs engines (e.g. AWS CloudWatch, Datadog).

## API Documentation
## http://localhost:8000/health/details#
## http://localhost:8000/health
## sql databse
http://localhost:8080/index.php?route=/database/structure&db=health_checker
## Standard User (Recommended)
Username: health_user
Password: health_password
## Root User
Username: root
Password: root_password#