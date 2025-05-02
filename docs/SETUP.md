# TinyURL Project Setup Guide

This document provides instructions for setting up and running the TinyURL project, which includes a gRPC server for user registration and a client application.

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 12 or higher
- Protocol Buffers compiler (protoc)
- Go plugins for protoc

## Project Structure

```
TinyURL/
├── client/
│   └── main.go         # gRPC client implementation
├── server/
│   └── main.go         # gRPC server implementation
├── proto/
│   └── tiny_url_api.proto  # Protocol buffer definitions
├── db/
│   └── db.go          # Database connection and configuration
├── go.mod
└── SETUP.md
```

## Database Setup

1. Create a PostgreSQL database:
```sql
CREATE DATABASE tinyurl;
```

2. Create the users table:
```sql
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

3. Configure database connection in `db/db.go`:
```go
const (
    host     = "localhost"
    port     = 5432
    user     = "your_username"
    password = "your_password"
    dbname   = "tinyurl"
)
```

## Protocol Buffer Setup

1. Install Protocol Buffers compiler:
```bash
# For macOS
brew install protobuf

# For Ubuntu/Debian
sudo apt-get install protobuf-compiler
```

2. Install Go plugins for protoc:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

3. Get Google API proto files:
```bash
git clone https://github.com/googleapis/googleapis.git
```

4. Generate Go code from proto files:
```bash
# Generate gRPC code
protoc -I . -I ./googleapis \
    --go_out . --go_opt paths=source_relative \
    --go-grpc_out . --go-grpc_opt paths=source_relative \
    proto/tiny_url_api.proto

# Generate gRPC-Gateway code
protoc -I . -I ./googleapis \
    --grpc-gateway_out . --grpc-gateway_opt paths=source_relative \
    proto/tiny_url_api.proto
```

## Running the Application

### 1. Start the gRPC Server

```bash
cd server
go run main.go
```

The server will start on port 7777.

### 2. Run the Client

```bash
cd client
go run main.go
```

The client will connect to the server and attempt to register a test user.

## API Documentation

### HTTP Endpoints (via gRPC-Gateway)

#### Register User
- **Endpoint**: `POST /v1/register`
- **Content-Type**: `application/json`
- **Request Body**:
  ```json
  {
    "first_name": "John",
    "last_name": "Doe"
  }
  ```
- **Success Response**:
  ```json
  {
    "status": "success",
    "username": "JoDoe"  // Generated from first half of first_name + last half of last_name
  }
  ```
- **Error Response**:
  ```json
  {
    "status": "error",
    "username": "JoDoe"  // Shows attempted username even in case of error
  }
  ```

Example using curl:
```bash
curl -X POST http://localhost:8080/v1/register \
  -H "Content-Type: application/json" \
  -d '{"first_name": "John", "last_name": "Doe"}'
```

### Username Generation Logic
The username is automatically generated using the following rules:
1. Takes the first half of the first name
2. Combines it with the last half of the last name
3. Example: "John" + "Smith" = "Joith"

Note: Usernames must be unique in the system. If a generated username already exists, the registration will fail with a duplicate key error.

### gRPC Service

#### RegisterUser

- **Service**: RegisterUser
- **Method**: RegisterUser
- **Request**:
  ```protobuf
  message RegisterUserRequest {
      string first_name = 1;
      string last_name = 2;
  }
  ```
- **Response**:
  ```protobuf
  message RegisteredUserResponse {
      string status = 1;
  }
  ```

## Error Handling

The application includes error handling for:
- Database connection failures
- User registration failures
- gRPC communication errors

## Testing

To test the setup:

1. Ensure PostgreSQL is running
2. Start the gRPC server
3. Run the client
4. Check the database to verify user registration

## Test Cases

### Testing User Registration

1. Start the server:
```bash
cd server && go run main.go
```

2. Test HTTP Registration Endpoint:
```bash
curl -X POST http://localhost:8080/v1/register \
  -H "Content-Type: application/json" \
  -d '{"first_name": "John", "last_name": "Doe"}'
```

Expected Response:
```json
{"status":"success"}
```

3. Test gRPC Registration Endpoint:
```bash
cd client && go run main.go
```

Expected Output:
```
User registered successfully
```

### Troubleshooting

If you encounter port conflicts:
1. Check for running processes:
```bash
lsof -i :7777,8080
```

2. Kill existing processes:
```bash
lsof -i :7777,8080 | grep LISTEN | awk '{print $2}' | xargs kill -9
```

3. Restart the server:
```bash
cd server && go run main.go
```

Note: The server runs two services:
- gRPC server on port 7777
- HTTP server on port 8080

## Latest Implementation Changes

### RPC Type Conversion

The latest implementation includes a change from streaming RPC to unary RPC for the RegisterUser service. This change was made to better match the service requirements and improve performance.

#### Before (Streaming RPC):
```go
// Server implementation
func (s *UserServer) RegisterUser(req *proto.RegisterUserRequest, stream proto.RegisterUser_RegisterUserServer) error {
    // ... streaming implementation
}

// Client implementation
stream, err := client.RegisterUser(ctx, req)
if err != nil {
    log.Fatalf("Error calling RegisterUser: %v", err)
}
response, err := stream.Recv()
```

#### After (Unary RPC):
```go
// Server implementation
func (s *UserServer) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.RegisteredUserResponse, error) {
    // ... unary implementation
    return &proto.RegisteredUserResponse{
        Status: "success",
    }, nil
}

// Client implementation
response, err := client.RegisterUser(ctx, req)
if err != nil {
    log.Fatalf("Error calling RegisterUser: %v", err)
}
```

#### Benefits of the Change:
1. Simpler implementation with direct request-response pattern
2. Reduced overhead by eliminating stream management
3. Better alignment with the service's actual requirements
4. Improved error handling with immediate feedback
5. More efficient resource usage

## Server Architecture

The server implements both gRPC and HTTP endpoints, running on different ports:
- gRPC server: `:7777`
- HTTP server: `:8080`

## Login Function Implementation

The login functionality is implemented through a gRPC service with the following components:

### Function Definition
```go
func (s *LoginServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error)
```

#### Components:
1. **Method Receiver**: `(s *LoginServer)`
   - Indicates this is a method of the `LoginServer` struct
   - Part of the gRPC service interface defined in proto file

2. **Parameters**:
   - `ctx context.Context`: Carries request-scoped values, deadlines, and cancellation signals
   - `req *proto.LoginRequest`: Pointer to the login request containing username

3. **Return Values**:
   - `*proto.LoginResponse`: Pointer to response containing user details
   - `error`: Any error that occurred during login process

### Example Flow:
```go
// Client Request:
{
    "username": "Bson"
}

// Server Response:
{
    "status": "success",
    "user": {
        "first_name": "Bob",
        "last_name": "Wilson",
        "username": "Bson"
    }
}
```

## Server Setup

1. **Database Initialization**:
```go
if err := db.InitDB(); err != nil {
    log.Fatalf("Failed to initialize database: %v", err)
}
defer db.CloseDB()
```

2. **Service Instances**:
```go
userServer := &service.UserServer{}
loginServer := &service.LoginServer{}
```

3. **gRPC Server Setup**:
```go
lis, err := net.Listen("tcp", ":7777")
grpcServer := grpc.NewServer()
```

4. **Service Registration**:
```go
proto.RegisterRegisterUserServer(grpcServer, userServer)
proto.RegisterLoginServiceServer(grpcServer, loginServer)
```

5. **HTTP Handlers**:
```go
http.HandleFunc("/v1/register", RegisterUserHandler(userServer))
http.HandleFunc("/v1/login", LoginHandler(loginServer))
```

## Running the Server

To start the server with all necessary components:
```bash
go run server/main.go server/register_user_api.go server/login_api.go
```

Note: Make sure no other process is using ports 7777 (gRPC) or 8080 (HTTP) before starting the server.

## Troubleshooting

1. **Port Already in Use**:
   - Error: `Failed to listen: listen tcp :7777: bind: address already in use`
   - Solution: Kill existing processes using the port
   ```bash
   pkill -f "server/main.go"
   # or
   kill -9 <PID>
   ```

2. **Database Connection**:
   - Ensure PostgreSQL is running
   - Check database credentials in configuration
   - Verify database schema is properly initialized

## Table
- ALTER TABLE generatedurl
ADD CONSTRAINT fk_generatedurl_username
FOREIGN KEY (username)
REFERENCES users(username);


## Database Health Check API testing
- curl -i http://localhost:8080/health
- HTTP/1.1 200 OK
- Content-Type: application/json
- Date: Fri, 02 May 2025 08:17:05 GMT
- Content-Length: 39
- {"status":"ok","database":"connected"}