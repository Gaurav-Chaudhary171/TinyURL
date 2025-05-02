# URL Shortening Implementation Changes

## Database Schema Changes
1. Created `generatedurl` table with the following structure:
```sql
CREATE TABLE generatedurl (
    url_id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    username VARCHAR(100) NOT NULL,
    tiny_url VARCHAR(100) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (username) REFERENCES users(username)
)
```

## Code Changes

### 1. URL Shortening Algorithm (`service/shorten_url_algo.go`)
```go
func ShortenURL(url string) string {
    hash := sha256.Sum256([]byte(url))
    encoded := base64.URLEncoding.EncodeToString(hash[:])
    return encoded[:8]
}
```
- Uses SHA-256 hashing
- Converts to base64 URL-safe encoding
- Returns first 8 characters for shorter URLs

### 2. Shorten Server Implementation (`service/shorten_server.go`)
Key changes:
- Added URL and username validation
- Modified URL storage format to include "https://" prefix
- Updated database query to use new table structure:
```go
query := `
    INSERT INTO generatedurl (original_url, username, tiny_url, is_active)
    VALUES ($1, $2, $3, true)
    RETURNING url_id`
```

### 3. Response Format
```go
return &proto.ShortenResponse{
    Status:   "success",
    ShortUrl: shortUrl,  // Now includes "https://" prefix
}
```

## API Endpoint
- Method: POST
- URL: `http://localhost:8080/v1/shorten`
- Request Body:
```json
{
    "url": "https://www.example.com/very/long/url",
    "username": "username"
}
```
- Response:
```json
{
    "status": "success",
    "shortUrl": "https://<shortened-url>"
}
```

## Error Handling
- URL validation: Returns error if URL is empty
- Username validation: Returns error if username is empty
- Database errors: Returns error if URL insertion fails

## Testing
To test the URL shortening endpoint:
```bash
curl -X POST -H "Content-Type: application/json" \
     -d '{"url": "https://www.example.com/very/long/url", "username": "Bson"}' \
     http://localhost:8080/v1/shorten
``` 