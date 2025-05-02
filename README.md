# TinyURL Service

A URL shortening service built with Go, GORM, and MySQL.

## Setup

1. Install dependencies:
```bash
go mod tidy
```

2. Set up environment variables (or use defaults):
```bash
export DB_USER=root
export DB_PASSWORD=your_password
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=tinyurl
```

3. Run the service:
```bash
go run main.go
```

## API Documentation

### Health Check Endpoint

The health check endpoint verifies the service's status and database connectivity.

**Endpoint:** `GET /health`

**Response:**
```json
{
    "status": "ok",
    "database": "connected",
    "timestamp": "2024-03-14T14:30:00Z"
}
```

**Status Codes:**
- `200 OK`: Service is healthy and database is connected
- `503 Service Unavailable`: Service is unhealthy or database is disconnected

**Response Fields:**
- `status`: Overall service status ("ok" or "error")
- `database`: Database connection status ("connected" or "disconnected")
- `timestamp`: Current server time in RFC3339 format

## Database Schema

The service uses two main tables:

### Users Table
- `user_id` (int64, primary key)
- `first_name` (string)
- `last_name` (string)
- `username` (string, unique)
- `dob` (string)
- `created_at` (timestamp)
- `updated_at` (timestamp)
- `deleted_at` (timestamp, nullable)

### Generated URLs Table
- `url_id` (int64, primary key)
- `username` (string)
- `original_url` (string)
- `tiny_url` (string, unique)
- `is_active` (boolean)
- `start_time` (timestamp)
- `end_time` (timestamp)
- `created_at` (timestamp)
- `updated_at` (timestamp)
- `deleted_at` (timestamp, nullable) 