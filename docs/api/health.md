# Health Check API

## Endpoint

```
GET /health
```

## Description

The health check endpoint provides information about the service's operational status, including the database connection health. This endpoint is useful for monitoring systems, load balancers, and deployment tools to verify that the service is running correctly.

## Request

### Method
- GET

### Headers
- No specific headers required

### Parameters
- No parameters required

## Response

### Success Response (200 OK)

```json
{
    "status": "healthy",
    "message": "Service is healthy"
}
```

### Error Response (503 Service Unavailable)

When the service is unhealthy (e.g., database connection issues):

```json
{
    "status": "error",
    "message": "Database connection error: {error details}"
}
```

or

```json
{
    "status": "error",
    "message": "Database ping failed: {error details}"
}
```

### Error Response (500 Internal Server Error)

When there's an unexpected error processing the request:

```json
{
    "status": "error",
    "message": "Internal server error"
}
```

## Response Fields

| Field   | Type   | Description                                                |
|---------|--------|------------------------------------------------------------|
| status  | string | Current status of the service ("healthy" or "error")        |
| message | string | Detailed message about the service health or error details  |

## Usage Examples

### cURL
```bash
curl -i http://localhost:8080/health
```

### HTTP
```http
GET /health HTTP/1.1
Host: localhost:8080
```

## Notes

- The endpoint checks the database connection by:
  1. Attempting to get the underlying SQL database connection
  2. Performing a ping operation to verify database connectivity
- Response times may vary based on database responsiveness
- Monitoring systems should consider implementing appropriate timeout settings
- This endpoint can be used in Kubernetes liveness and readiness probes 