# TinyURL Development Workflow

## Development Process

### 1. Initial Setup
- Created project structure with proto, server, service, and db directories
- Set up Go modules and dependencies
- Configured PostgreSQL database connection

### 2. Protocol Buffer Definition
- Defined service interfaces in `proto/tinyurl.proto`
- Implemented four main services:
  - RegisterUser
  - LoginService
  - ShortenURL
  - ExtendedURL
- Generated Go code from proto definitions

### 3. Database Schema
```sql
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    username VARCHAR(50) UNIQUE
);

CREATE TABLE generatedurl (
    url_id SERIAL PRIMARY KEY,
    original_url TEXT,
    username VARCHAR(50),
    tiny_url VARCHAR(100),
    is_active BOOLEAN,
    FOREIGN KEY (username) REFERENCES users(username)
);
```

### 4. Service Implementation
1. **User Registration**
   - Automatic username generation from first and last name
   - Unique constraint on usernames
   - Error handling for duplicate usernames

2. **Login Service**
   - Username validation
   - User information retrieval
   - Error handling for non-existent users

3. **URL Shortening**
   - SHA-256 hash generation
   - Base64 encoding for URL shortening
   - User association with shortened URLs
   - Active/inactive URL tracking

4. **URL Extension**
   - Original URL retrieval
   - User validation
   - Error handling for non-existent URLs

### 5. API Gateway
- Implemented gRPC-Gateway for REST endpoints
- Added HTTP handlers for all services
- Configured CORS and content type headers

### 6. Testing Workflow
1. **User Registration**
   ```bash
   curl -X POST -H "Content-Type: application/json" \
   -d '{"first_name": "John", "last_name": "Doe"}' \
   http://localhost:8080/v1/registeruser
   ```

2. **Login**
   ```bash
   curl -X POST -H "Content-Type: application/json" \
   -d '{"username": "JohDoe"}' \
   http://localhost:8080/v1/login
   ```

3. **URL Shortening**
   ```bash
   curl -X POST -H "Content-Type: application/json" \
   -d '{"url": "https://example.com", "username": "JohDoe"}' \
   http://localhost:8080/v1/shortenurl
   ```

4. **URL Extension**
   ```bash
   curl -X POST -H "Content-Type: application/json" \
   -d '{"url": "shortened_url", "username": "JohDoe"}' \
   http://localhost:8080/v1/extendurl
   ```

## Key Decisions and Improvements

### 1. Username Generation
- Implemented algorithm to generate usernames from first and last names
- Added unique constraint to prevent duplicates
- Considered adding retry logic for duplicate usernames

### 2. URL Shortening
- Used SHA-256 for consistent hashing
- Added "https://" prefix to shortened URLs
- Implemented user association for security

### 3. Error Handling
- Standardized error response format
- Added detailed error messages
- Implemented proper HTTP status codes

### 4. Database Design
- Used foreign key constraints for data integrity
- Added is_active flag for URL management
- Implemented proper indexing for performance

## Future Improvements
1. Add rate limiting
2. Implement URL expiration
3. Add analytics tracking
4. Implement caching layer
5. Add user authentication with JWT
6. Implement URL validation
7. Add bulk URL operations
8. Implement URL categories/tags 