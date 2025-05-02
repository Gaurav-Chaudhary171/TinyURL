# TinyURL Service

A gRPC-based URL shortening service with REST API endpoints.

## Features
- User registration with automatic username generation
- User authentication
- URL shortening with unique hash generation
- URL expansion to original URLs
- PostgreSQL database for data persistence

## Tech Stack
- Go 1.24
- gRPC
- gRPC-Gateway
- PostgreSQL
- Protocol Buffers

## Project Structure
```
TinyURL/
├── proto/           # Protocol buffer definitions
├── server/          # Server implementation
├── service/         # Business logic
└── db/             # Database operations
```

## Setup
1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Generate proto files:
   ```bash
   protoc -I . -I $(go env GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
   --go_out=. --go_opt=paths=source_relative \
   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
   --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
   proto/tinyurl.proto
   ```

3. Build and run:
   ```bash
   cd server
   go run main.go register_user_api.go login_api.go shorten_api.go gateway.go
   ```

## API Documentation

### 1. Register User
**Endpoint:** `POST /v1/registeruser`  
**Request:**
```json
{
    "first_name": "John",
    "last_name": "Doe"
}
```
**Response:**
```json
{
    "status": "success",
    "username": "JohDoe"
}
```

### 2. Login
**Endpoint:** `POST /v1/login`  
**Request:**
```json
{
    "username": "JohDoe"
}
```
**Response:**
```json
{
    "status": "success",
    "user": {
        "firstName": "John",
        "lastName": "Doe",
        "username": "JohDoe"
    }
}
```

### 3. Shorten URL
**Endpoint:** `POST /v1/shortenurl`  
**Request:**
```json
{
    "url": "https://www.example.com/very/long/url",
    "username": "JohDoe"
}
```
**Response:**
```json
{
    "status": "success",
    "shortenurl": "https://abc123",
    "originalurl": "https://www.example.com/very/long/url"
}
```

### 4. Extend URL
**Endpoint:** `POST /v1/extendurl`  
**Request:**
```json
{
    "url": "https://abc123",
    "username": "JohDoe"
}
```
**Response:**
```json
{
    "status": "success",
    "originalurl": "https://www.example.com/very/long/url",
    "extenedurl": "https://abc123"
}
```

## Error Handling
All endpoints return errors in the following format:
```json
{
    "code": 2,
    "message": "Error description",
    "details": []
}
```

Common error scenarios:
- Duplicate username during registration
- Invalid username during login
- URL not found during extension
- Database constraint violations 