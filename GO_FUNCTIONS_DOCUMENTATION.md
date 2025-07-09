# Go Functions and Components Documentation

## Table of Contents
1. [API Service (cmd/api/main.go)](#api-service-cmdapimaingo)
2. [Gateway Service (cmd/gateway/main.go)](#gateway-service-cmdgatewaymaingo)
3. [Gateway Package (pkg/gateway/)](#gateway-package-pkggateway)
4. [Data Structures](#data-structures)
5. [Middleware Functions](#middleware-functions)
6. [Database Operations](#database-operations)
7. [Utility Functions](#utility-functions)

## API Service (cmd/api/main.go)

### Main Function

#### `main()`
The main entry point for the API service.

**Responsibilities:**
- Database connection setup
- HTTP route registration
- Server startup

**Environment Variables:**
- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (default: user)
- `DB_PASSWORD`: Database password (default: pass)
- `DB_NAME`: Database name (default: users)
- `PORT`: API service port (default: 8085)

**Example Usage:**
```bash
DB_HOST=my-db-host DB_PORT=5432 DB_USER=myuser DB_PASSWORD=mypass DB_NAME=mydb PORT=8085 go run cmd/api/main.go
```

### Environment Helper

#### `getEnv(key, defaultValue string) string`
Retrieves environment variable value with fallback to default.

**Parameters:**
- `key`: Environment variable name
- `defaultValue`: Default value if environment variable is not set

**Returns:** Environment variable value or default value

**Example:**
```go
dbHost := getEnv("DB_HOST", "localhost")
```

### Authentication Middleware

#### `jwtAuth(next http.HandlerFunc) http.HandlerFunc`
JWT authentication middleware that validates Bearer tokens.

**Parameters:**
- `next`: Next HTTP handler function

**Returns:** HTTP handler function with JWT validation

**Features:**
- Validates Authorization header format: `Bearer <token>`
- Parses and validates JWT token
- Returns 401 Unauthorized for invalid/missing tokens

**Example Usage:**
```go
http.HandleFunc("/api/customers", jwtAuth(customersHandler))
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid Bearer token
- `401 Unauthorized`: Invalid JWT token

### Authentication Handler

#### `loginHandler(w http.ResponseWriter, r *http.Request)`
Handles user login and JWT token generation.

**HTTP Method:** POST

**Request Body:**
```json
{
  "email": "demo@example.com",
  "password": "demo123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Features:**
- Validates user credentials against database
- Generates JWT token with 24-hour expiration
- Uses hardcoded secret (should be moved to environment)

**Error Responses:**
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: User not found or invalid password
- `405 Method Not Allowed`: Non-POST requests

## Data Structures

### Customer Struct
```go
type Customer struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    CreatedAt time.Time `json:"created_at"`
}
```

**Fields:**
- `ID`: Unique customer identifier (auto-generated)
- `Name`: Customer's full name
- `Email`: Customer's email address (unique)
- `Phone`: Customer's phone number (optional)
- `CreatedAt`: Timestamp when customer was created

### Contact Struct
```go
type Contact struct {
    ID           int    `json:"id"`
    CustomerName string `json:"customer_name"`
    Message      string `json:"message"`
    Date         string `json:"date"`
}
```

**Fields:**
- `ID`: Unique contact identifier (auto-generated)
- `CustomerName`: Name of the associated customer
- `Message`: Contact message content
- `Date`: Formatted date string (YYYY-MM-DD)

## Customer Management Handlers

### `customersHandler(w http.ResponseWriter, r *http.Request)`
Main handler for customer operations.

**Supported Methods:**
- `GET`: Retrieve all customers
- `POST`: Create new customer

**Example Usage:**
```go
http.HandleFunc("/api/customers", jwtAuth(customersHandler))
```

### `getCustomers(w http.ResponseWriter, r *http.Request)`
Retrieves all customers from the database.

**HTTP Method:** GET

**Response:**
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "created_at": "2024-01-15T10:30:00Z"
  }
]
```

**Database Query:**
```sql
SELECT id, name, email, phone, created_at FROM customers ORDER BY id DESC
```

**Error Handling:**
- `500 Internal Server Error`: Database query errors

### `createCustomer(w http.ResponseWriter, r *http.Request)`
Creates a new customer in the database.

**HTTP Method:** POST

**Request Body:**
```json
{
  "name": "Jane Smith",
  "email": "jane@example.com",
  "phone": "+0987654321"
}
```

**Response (201 Created):**
```json
{
  "id": 2,
  "name": "Jane Smith",
  "email": "jane@example.com",
  "phone": "+0987654321",
  "created_at": "2024-01-15T11:00:00Z"
}
```

**Database Query:**
```sql
INSERT INTO customers (name, email, phone) VALUES ($1, $2, $3) RETURNING id, created_at
```

**Error Handling:**
- `400 Bad Request`: Invalid JSON request body
- `500 Internal Server Error`: Database insertion errors

### `customerDetailHandler(w http.ResponseWriter, r *http.Request)`
Handles individual customer operations.

**Supported Methods:**
- `PUT`: Update customer
- `DELETE`: Delete customer

**URL Pattern:** `/api/customers/{id}`

**Example Usage:**
```go
http.HandleFunc("/api/customers/", jwtAuth(customerDetailHandler))
```

### Update Customer (PUT)
Updates an existing customer.

**Request Body:**
```json
{
  "name": "Jane Smith Updated",
  "email": "jane.updated@example.com",
  "phone": "+1111111111"
}
```

**Database Query:**
```sql
UPDATE customers SET name=$1, email=$2, phone=$3 WHERE id=$4
```

**Error Handling:**
- `400 Bad Request`: Missing ID or invalid JSON
- `500 Internal Server Error`: Database update errors

### Delete Customer (DELETE)
Deletes a customer and all associated contacts (CASCADE).

**Database Query:**
```sql
DELETE FROM customers WHERE id=$1
```

**Response:** 204 No Content

**Error Handling:**
- `400 Bad Request`: Missing ID
- `500 Internal Server Error`: Database deletion errors

## Contact Management Handlers

### `contactsHandler(w http.ResponseWriter, r *http.Request)`
Main handler for contact operations.

**Supported Methods:**
- `GET`: Retrieve all contacts
- `POST`: Create new contact

**Example Usage:**
```go
http.HandleFunc("/api/contacts", jwtAuth(contactsHandler))
```

### `getContacts(w http.ResponseWriter, r *http.Request)`
Retrieves all contacts with customer information.

**HTTP Method:** GET

**Response:**
```json
[
  {
    "id": 1,
    "customer_name": "John Doe",
    "message": "Follow up on proposal",
    "date": "2024-01-15"
  }
]
```

**Database Query:**
```sql
SELECT c.id, cu.name, c.message, to_char(c.created_at, 'YYYY-MM-DD') 
FROM contacts c 
JOIN customers cu ON c.customer_id = cu.id 
ORDER BY c.id DESC
```

**Error Handling:**
- `500 Internal Server Error`: Database query errors

### `createContact(w http.ResponseWriter, r *http.Request)`
Creates a new contact for a customer.

**HTTP Method:** POST

**Request Body:**
```json
{
  "customer_id": 1,
  "message": "Customer interested in premium package"
}
```

**Response (201 Created):**
```json
{
  "id": 2,
  "customer_name": "John Doe",
  "message": "Customer interested in premium package",
  "date": "2024-01-15"
}
```

**Database Query:**
```sql
INSERT INTO contacts (customer_id, message) VALUES ($1, $2) 
RETURNING id, (SELECT name FROM customers WHERE id = $1), message, to_char(created_at, 'YYYY-MM-DD')
```

**Validation:**
- `customer_id` must be non-zero
- `message` must not be empty

**Error Handling:**
- `400 Bad Request`: Invalid JSON or missing required fields
- `500 Internal Server Error`: Database insertion errors

## Gateway Service (cmd/gateway/main.go)

### Main Function

#### `main()`
The main entry point for the API Gateway service.

**Responsibilities:**
- HTTP route registration
- Middleware setup
- Server startup

**Routes:**
- `/healthz`: Health check endpoint
- `/api`: Protected API endpoint with rate limiting and IP allowlist

### Rate Limiting

#### `rateLimit(next http.HandlerFunc) http.HandlerFunc`
Rate limiting middleware that limits requests per IP.

**Parameters:**
- `next`: Next HTTP handler function

**Returns:** HTTP handler function with rate limiting

**Features:**
- 10 requests per minute per IP
- Resets every minute
- In-memory storage (not persistent across restarts)

**Error Response:**
- `429 Too Many Requests`: Rate limit exceeded

**Example Usage:**
```go
http.HandleFunc("/api", ipAllowlist(rateLimit(apiHandler)))
```

### IP Allowlist

#### `ipAllowlist(next http.HandlerFunc) http.HandlerFunc`
IP allowlist middleware that restricts access to specific IPs.

**Parameters:**
- `next`: Next HTTP handler function

**Returns:** HTTP handler function with IP filtering

**Allowed IPs:**
- `127.0.0.1` (localhost)
- `::1` (localhost IPv6)

**Error Response:**
- `403 Forbidden`: IP not in allowlist

**Example Usage:**
```go
http.HandleFunc("/api", ipAllowlist(rateLimit(apiHandler)))
```

### Rate Limiter Struct
```go
type rateLimiter struct {
    count     int
    timestamp time.Time
}
```

**Fields:**
- `count`: Number of requests in current window
- `timestamp`: Start time of current window

### API Handler

#### `apiHandler(w http.ResponseWriter, r *http.Request)`
Protected API endpoint handler.

**Response:**
```
JWT dogrulamasi gecici olarak devre disi! Korumali /api endpointindesin.
```

**Note:** Currently returns a static message. In production, this would forward requests to the API service.

## Gateway Package (pkg/gateway/)

### Health Check Handler

#### `HealthzHandler(w http.ResponseWriter, r *http.Request)`
Health check endpoint handler.

**Response:** "ok"

**Usage:**
```go
http.HandleFunc("/healthz", gateway.HealthzHandler)
```

**Example:**
```bash
curl http://localhost:8080/healthz
# Response: ok
```

## Utility Functions

### `atoi(s string) int`
Converts string to integer using `fmt.Sscanf`.

**Parameters:**
- `s`: String to convert

**Returns:** Integer value

**Note:** Ignores conversion errors and returns 0 if conversion fails.

**Example:**
```go
id := atoi("123") // Returns 123
id := atoi("abc") // Returns 0
```

## Database Operations

### Connection Setup
```go
dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
    dbUser, dbPassword, dbHost, dbPort, dbName)

db, err = sql.Open("postgres", dbURL)
if err != nil {
    log.Fatal("DB bağlantı hatası:", err)
}
defer db.Close()

if err = db.Ping(); err != nil {
    log.Fatal("DB ping hatası:", err)
}
```

### Query Patterns

#### Parameterized Queries
All database queries use parameterized statements to prevent SQL injection:

```go
// SELECT with parameters
err := db.QueryRow("SELECT id, email, password_hash FROM users WHERE email=$1", req.Email).Scan(&id, &email, &passwordHash)

// INSERT with RETURNING
err := db.QueryRow("INSERT INTO customers (name, email, phone) VALUES ($1, $2, $3) RETURNING id, created_at", c.Name, c.Email, c.Phone).Scan(&c.ID, &c.CreatedAt)

// UPDATE
_, err := db.Exec("UPDATE customers SET name=$1, email=$2, phone=$3 WHERE id=$4", c.Name, c.Email, c.Phone, id)

// DELETE
_, err := db.Exec("DELETE FROM customers WHERE id=$1", id)
```

#### Row Scanning
```go
rows, err := db.Query("SELECT id, name, email, phone, created_at FROM customers ORDER BY id DESC")
if err != nil {
    // Handle error
}
defer rows.Close()

customers := []Customer{}
for rows.Next() {
    var c Customer
    if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.CreatedAt); err != nil {
        // Handle error
    }
    customers = append(customers, c)
}
```

## Error Handling Patterns

### HTTP Error Responses
```go
// Bad Request
http.Error(w, "Geçersiz istek", http.StatusBadRequest)

// Unauthorized
http.Error(w, "Yetkisiz: Bearer token gerekli", http.StatusUnauthorized)

// Method Not Allowed
http.Error(w, "Yöntem desteklenmiyor", http.StatusMethodNotAllowed)

// Internal Server Error
http.Error(w, "DB Hatasi", http.StatusInternalServerError)

// No Content
w.WriteHeader(http.StatusNoContent)
```

### JSON Error Responses
```go
// For API endpoints that return JSON
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]string{"message": "Error description"})
```

### Logging
```go
log.Printf("DB sorgu hatasi: %v", err)
log.Printf("DB ekleme hatasi: %v", err)
log.Printf("DB contact ekleme hatasi: %v", err)
```

## Security Considerations

### JWT Implementation
```go
var jwtSecret = []byte("supersecret") // Should be moved to environment

token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "email": email,
    "exp":   time.Now().Add(24 * time.Hour).Unix(),
})
tokenStr, _ := token.SignedString(jwtSecret)
```

### Password Security
**Current Implementation (Demo):**
```go
// Demo için şifre hash kontrolü yok, gerçek projede bcrypt kullan!
if req.Password != "demo123" && req.Password != passwordHash {
    http.Error(w, "Şifre hatalı", http.StatusUnauthorized)
    return
}
```

**Recommended Implementation:**
```go
import "golang.org/x/crypto/bcrypt"

// Hash password
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// Compare password
err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
```

## Performance Considerations

### Database Connection Pooling
The current implementation uses the default connection pool. For production:

```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)
```

### Rate Limiting Storage
Current rate limiting uses in-memory storage. For production, consider:
- Redis for distributed rate limiting
- Database storage for persistence
- More sophisticated algorithms (token bucket, sliding window)

## Testing Considerations

### Unit Testing
```go
func TestGetCustomers(t *testing.T) {
    // Mock database
    // Test handler
    // Assert response
}
```

### Integration Testing
```go
func TestCustomerCRUD(t *testing.T) {
    // Setup test database
    // Create customer
    // Read customer
    // Update customer
    // Delete customer
    // Cleanup
}
```

## Monitoring and Observability

### Health Checks
```go
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
    // Check database connectivity
    if err := db.Ping(); err != nil {
        http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))
}
```

### Metrics
Consider adding Prometheus metrics:
```go
import "github.com/prometheus/client_golang/prometheus"

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
)
```

## Future Enhancements

1. **Structured Logging**: Implement structured logging with correlation IDs
2. **Circuit Breaker**: Add circuit breaker pattern for external dependencies
3. **Caching**: Implement Redis caching for frequently accessed data
4. **API Versioning**: Add versioning support for API endpoints
5. **OpenAPI Documentation**: Generate OpenAPI/Swagger documentation
6. **Graceful Shutdown**: Implement graceful shutdown with context cancellation
7. **Configuration Management**: Use Viper for configuration management
8. **Dependency Injection**: Implement proper dependency injection