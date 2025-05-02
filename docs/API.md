# TinyURL API Documentation

## User Registration

Register a new user with first name and last name. The API will generate a unique username.

**Endpoint:** `POST /v1/registeruser`

**Request Body:**
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
    "username": "Johoe"  // Generated username
}
```

## User Login

Login with a username to retrieve user information.

**Endpoint:** `POST /v1/login`

**Request Body:**
```json
{
    "username": "Johoe"
}
```

**Response:**
```json
{
    "status": "success",
    "user": {
        "first_name": "John",
        "last_name": "Doe",
        "username": "Johoe"
    }
}
```

## Error Responses

Both endpoints may return the following error responses:

- `400 Bad Request`: Invalid request body
- `404 Not Found`: User not found (login only)
- `405 Method Not Allowed`: Wrong HTTP method
- `500 Internal Server Error`: Database or server error

## Test Cases

### Registration API

1. Successful Registration:
```bash
curl -X POST -H "Content-Type: application/json" \
     -d '{"first_name": "Bob", "last_name": "Wilson"}' \
     http://localhost:8080/v1/registeruser

Response: {"status":"success","username":"Bson"}
```

2. Duplicate Username Error:
```bash
# Trying to register with same name pattern
curl -X POST -H "Content-Type: application/json" \
     -d '{"first_name": "Bob", "last_name": "Wilson"}' \
     http://localhost:8080/v1/registeruser

Response: ERROR: duplicate key value violates unique constraint "users_username_key"
```

### Login API

1. Successful Login:
```bash
curl -X POST -H "Content-Type: application/json" \
     -d '{"username": "Bson"}' \
     http://localhost:8080/v1/login

Response: {"status":"success","user":{"first_name":"Bob","last_name":"Wilson","username":"Bson"}}
```

2. User Not Found:
```bash
curl -X POST -H "Content-Type: application/json" \
     -d '{"username": "NonExistent"}' \
     http://localhost:8080/v1/login

Response: User not found
```

## Example Usage

## Register User
curl -X POST -H "Content-Type: application/json" -d '{"first_name": "Christopher", "last_name": "Anders2n"}' http://localhost:8080/v1/registeruser

## URL shortening enpoint
curl -X POST -H "Content-Type: application/json" -d '{"url": "https://www.example.com/very/lon", "username": "Chrisrs2n"}' http://localhost:8080/v1/shortenurl

## Login enpoint
curl -X POST -H "Content-Type: application/json" -d '{"username": "Chrisrs2n"}' http://localhost:8080/v1/login

## Extend URL
curl -X POST -H "Content-Type: application/json" -d '{"url": "https://JF-rqVwI", "username": "Chrisrson"}' http://localhost:8080/v1/extendurl

