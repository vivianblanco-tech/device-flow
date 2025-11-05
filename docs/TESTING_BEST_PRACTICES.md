# Testing Best Practices

**Last Updated:** November 4, 2025  
**Status:** Active Guidelines  
**Applies To:** All Go tests in the project

---

## 📋 Table of Contents

1. [Quick Reference](#quick-reference)
2. [Test Organization](#test-organization)
3. [Database Testing](#database-testing)
4. [Time-Based Testing](#time-based-testing)
5. [Test Data Management](#test-data-management)
6. [Common Patterns](#common-patterns)
7. [What to Avoid](#what-to-avoid)
8. [Examples](#examples)

---

## Quick Reference

### ✅ DO

```go
// Use descriptive test names
func TestUserLogin_WithValidCredentials_ReturnsSuccess(t *testing.T) {}

// Capture reference time once
now := time.Now()
future := now.Add(24 * time.Hour)

// Use unique test data
name := fmt.Sprintf("Test_%s_%d", t.Name(), time.Now().UnixNano())

// Clean up resources
defer cleanup()

// Use subtests for organization
t.Run("happy path", func(t *testing.T) {})
t.Run("error case", func(t *testing.T) {})
```

### ❌ DON'T

```go
// Don't use cryptic names
func TestStuff(t *testing.T) {}

// Don't call time.Now() multiple times
expires1 := time.Now().Add(1 * time.Hour) // ❌
expires2 := time.Now().Add(2 * time.Hour) // ❌

// Don't use static test data
user := &User{Email: "test@example.com"} // ❌ Collision risk

// Don't skip cleanup
// defer cleanup() // ❌ Missing!

// Don't write mega-tests
func TestEverything(t *testing.T) {} // ❌ Too broad
```

---

## Test Organization

### File Structure

```
internal/
  auth/
    password.go           # Implementation
    password_test.go      # Unit tests
    session.go           # Implementation  
    session_test.go      # Unit tests
    integration_test.go  # Integration tests (if needed)
```

### Naming Conventions

```go
// Unit test: Test[FunctionName]
func TestHashPassword(t *testing.T) {}
func TestValidateEmail(t *testing.T) {}

// Integration test: TestIntegration[Feature]
func TestIntegrationUserLogin(t *testing.T) {}
func TestIntegrationDatabaseMigration(t *testing.T) {}

// Subtests: Use descriptive names with spaces
t.Run("returns error when password is too short", func(t *testing.T) {})
t.Run("accepts valid password with special characters", func(t *testing.T) {})
```

### Test Structure (Arrange-Act-Assert)

```go
func TestUserLogin(t *testing.T) {
    // Arrange - Set up test data and dependencies
    db, cleanup := database.SetupTestDB(t)
    defer cleanup()
    
    user := createTestUser(t, db, "test@example.com")
    handler := NewAuthHandler(db, templates)
    
    // Act - Execute the function being tested
    result, err := handler.Login(user.Email, "password123")
    
    // Assert - Verify the results
    if err != nil {
        t.Fatalf("Expected no error, got: %v", err)
    }
    if result == nil {
        t.Error("Expected result, got nil")
    }
}
```

### Using Subtests

```go
func TestPasswordValidation(t *testing.T) {
    tests := []struct {
        name        string
        password    string
        shouldError bool
        errorMsg    string
    }{
        {
            name:        "valid password with all requirements",
            password:    "StrongPass123!",
            shouldError: false,
        },
        {
            name:        "password too short",
            password:    "Short1!",
            shouldError: true,
            errorMsg:    "password must be at least 8 characters",
        },
        {
            name:        "password missing uppercase",
            password:    "weakpass123!",
            shouldError: true,
            errorMsg:    "password must contain an uppercase letter",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidatePassword(tt.password)
            
            if tt.shouldError {
                if err == nil {
                    t.Error("Expected error, got nil")
                }
                if err != nil && !strings.Contains(err.Error(), tt.errorMsg) {
                    t.Errorf("Expected error containing %q, got %q", tt.errorMsg, err.Error())
                }
            } else {
                if err != nil {
                    t.Errorf("Expected no error, got: %v", err)
                }
            }
        })
    }
}
```

---

## Database Testing

### Setup and Cleanup

```go
func TestUserCRUD(t *testing.T) {
    // Use the test helper
    db, cleanup := database.SetupTestDB(t)
    defer cleanup() // Always defer cleanup
    
    // Your test code here
}
```

### Transaction-Based Isolation (Recommended)

```go
func TestUserUpdate(t *testing.T) {
    db := getTestDB(t)
    
    // Start transaction
    tx, err := db.Begin()
    if err != nil {
        t.Fatalf("Failed to start transaction: %v", err)
    }
    defer tx.Rollback() // Always rolls back, even on success
    
    // Use tx for all database operations
    _, err = tx.Exec("INSERT INTO users ...")
    // Test continues with tx
    
    // No need to manually clean up - rollback handles it
}
```

### Avoiding Test Data Collisions

```go
// ❌ BAD - Static data
company := &ClientCompany{
    Name: "Test Company", // Will fail on second run
}

// ✅ GOOD - Unique with timestamp
company := &ClientCompany{
    Name: fmt.Sprintf("Test Company %d", time.Now().UnixNano()),
}

// ✅ BETTER - UUID
company := &ClientCompany{
    Name: fmt.Sprintf("Test Company %s", uuid.New().String()),
}

// ✅ BEST - Test name + timestamp
company := &ClientCompany{
    Name: fmt.Sprintf("Test_%s_%d", 
        strings.ReplaceAll(t.Name(), "/", "_"), 
        time.Now().UnixNano()),
}
```

### Querying Test Data

```go
// Always use explicit context
ctx := context.Background()

// Use QueryRowContext for single results
var user User
err := db.QueryRowContext(ctx, 
    "SELECT id, email FROM users WHERE email = $1",
    "test@example.com",
).Scan(&user.ID, &user.Email)

// Check for no rows explicitly
if err == sql.ErrNoRows {
    t.Error("User not found")
} else if err != nil {
    t.Fatalf("Query failed: %v", err)
}
```

### Test Helpers

```go
// Create reusable test helpers
func createTestUser(t *testing.T, db *sql.DB, email string) *User {
    t.Helper() // Marks this as a helper function
    
    user := &User{
        Email: email,
        PasswordHash: "hashedpassword",
        Role: RoleLogistics,
    }
    user.BeforeCreate()
    
    err := db.QueryRowContext(
        context.Background(),
        `INSERT INTO users (email, password_hash, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5) RETURNING id`,
        user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
    ).Scan(&user.ID)
    
    if err != nil {
        t.Fatalf("Failed to create test user: %v", err)
    }
    
    return user
}

// Usage
func TestSomething(t *testing.T) {
    db, cleanup := database.SetupTestDB(t)
    defer cleanup()
    
    user := createTestUser(t, db, "test@example.com")
    // Use user in test
}
```

---

## Time-Based Testing

### Consistent Time References

```go
// ❌ BAD - Multiple time.Now() calls create race conditions
func TestSessionExpiration(t *testing.T) {
    session1 := &Session{
        ExpiresAt: time.Now().Add(1 * time.Hour), // Different time
    }
    time.Sleep(10 * time.Millisecond)
    session2 := &Session{
        ExpiresAt: time.Now().Add(1 * time.Hour), // Different time
    }
    // session1.ExpiresAt != session2.ExpiresAt (race condition!)
}

// ✅ GOOD - Single reference time
func TestSessionExpiration(t *testing.T) {
    now := time.Now() // Capture once
    
    session1 := &Session{
        ExpiresAt: now.Add(1 * time.Hour), // Consistent
    }
    session2 := &Session{
        ExpiresAt: now.Add(1 * time.Hour), // Consistent
    }
    // session1.ExpiresAt == session2.ExpiresAt ✓
}
```

### Avoiding Date Boundaries

```go
// ❌ BAD - Testing at exact boundaries
func TestShipmentsLast30Days(t *testing.T) {
    now := time.Now()
    shipment := &Shipment{
        CreatedAt: now.AddDate(0, 0, -30), // Exactly 30 days ago
    }
    // Query: "last 30 days" might exclude this due to time-of-day
}

// ✅ GOOD - Safe margin from boundaries
func TestShipmentsLast30Days(t *testing.T) {
    now := time.Now()
    shipment := &Shipment{
        CreatedAt: now.AddDate(0, 0, -29), // 29 days ago (safely within window)
    }
    // Query: "last 30 days" will reliably include this
}
```

### Time Formatting

```go
// Always use consistent format
const DateFormat = "2006-01-02"
const DateTimeFormat = "2006-01-02 15:04:05"

// For date comparisons in tests
expectedDate := now.Format(DateFormat)
actualDate := result.CreatedAt.Format(DateFormat)
if expectedDate != actualDate {
    t.Errorf("Expected date %s, got %s", expectedDate, actualDate)
}
```

### Testing Time-Dependent Code

```go
// Option 1: Use interfaces for time (dependency injection)
type Clock interface {
    Now() time.Time
}

type RealClock struct{}
func (c RealClock) Now() time.Time { return time.Now() }

type MockClock struct {
    CurrentTime time.Time
}
func (c MockClock) Now() time.Time { return c.CurrentTime }

// In production code
type Service struct {
    clock Clock
}

// In tests
func TestSomething(t *testing.T) {
    fixedTime := time.Date(2025, 11, 4, 12, 0, 0, 0, time.UTC)
    service := &Service{
        clock: MockClock{CurrentTime: fixedTime},
    }
    // Test with predictable time
}

// Option 2: Pass time as parameter
func ProcessExpiredSessions(db *sql.DB, currentTime time.Time) error {
    // Use currentTime instead of time.Now()
}

// In tests
func TestProcessExpiredSessions(t *testing.T) {
    fixedTime := time.Date(2025, 11, 4, 12, 0, 0, 0, time.UTC)
    err := ProcessExpiredSessions(db, fixedTime)
    // Predictable behavior
}
```

---

## Test Data Management

### Unique Identifiers

```go
// Helper function for unique names
func uniqueName(prefix string) string {
    return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

// Usage
company := &ClientCompany{
    Name: uniqueName("TestCompany"),
}

user := &User{
    Email: fmt.Sprintf("test_%d@example.com", time.Now().UnixNano()),
}
```

### Test Fixtures

```go
// Define reusable test data structures
var testUsers = []struct {
    email string
    role  string
}{
    {"logistics@test.com", RoleLogistics},
    {"client@test.com", RoleClient},
    {"warehouse@test.com", RoleWarehouse},
}

func TestPermissions(t *testing.T) {
    db, cleanup := database.SetupTestDB(t)
    defer cleanup()
    
    // Create all test users
    for _, tu := range testUsers {
        createTestUser(t, db, tu.email, tu.role)
    }
    
    // Run tests
}
```

### Cleanup Strategies

```go
// Strategy 1: Defer cleanup (recommended)
func TestSomething(t *testing.T) {
    db, cleanup := database.SetupTestDB(t)
    defer cleanup() // Always runs, even on panic
    
    // Test code
}

// Strategy 2: Transaction rollback (best for isolation)
func TestSomething(t *testing.T) {
    db := getTestDB(t)
    tx, _ := db.Begin()
    defer tx.Rollback() // Automatic cleanup
    
    // Use tx for all queries
}

// Strategy 3: Manual cleanup (when defer isn't enough)
func TestSomething(t *testing.T) {
    db, cleanup := database.SetupTestDB(t)
    defer cleanup()
    
    // Create test data
    userID := createTestUser(t, db, "test@example.com")
    
    // Ensure specific cleanup
    defer func() {
        db.Exec("DELETE FROM audit_logs WHERE user_id = $1", userID)
        db.Exec("DELETE FROM users WHERE id = $1", userID)
    }()
    
    // Test code
}
```

---

## Common Patterns

### Testing HTTP Handlers

```go
func TestLoginHandler(t *testing.T) {
    // Setup
    db, cleanup := database.SetupTestDB(t)
    defer cleanup()
    
    handler := NewAuthHandler(db, templates)
    
    // Create request
    form := url.Values{}
    form.Set("email", "test@example.com")
    form.Set("password", "password123")
    
    req := httptest.NewRequest(
        http.MethodPost, 
        "/login", 
        strings.NewReader(form.Encode()),
    )
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    
    // Record response
    w := httptest.NewRecorder()
    
    // Call handler
    handler.Login(w, req)
    
    // Assert response
    if w.Code != http.StatusSeeOther {
        t.Errorf("Expected status 303, got %d", w.Code)
    }
    
    // Check cookies
    cookies := w.Result().Cookies()
    found := false
    for _, cookie := range cookies {
        if cookie.Name == "session" {
            found = true
            break
        }
    }
    if !found {
        t.Error("Session cookie not set")
    }
}
```

### Testing with Context

```go
func TestWithContext(t *testing.T) {
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Create request with context
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    req = req.WithContext(ctx)
    
    // Add values to context
    user := &User{ID: 1, Email: "test@example.com"}
    ctx = context.WithValue(req.Context(), middleware.UserContextKey, user)
    req = req.WithContext(ctx)
    
    // Test with context
    handler(w, req)
}
```

### Table-Driven Tests

```go
func TestEmailValidation(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"missing @", "userexample.com", true},
        {"missing domain", "user@", true},
        {"empty", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

---

## What to Avoid

### 1. Testing Implementation Details

```go
// ❌ BAD - Testing how something works
func TestUserLoginImplementation(t *testing.T) {
    // Tests that bcrypt is used with specific cost
    // Tests exact SQL query structure
    // Breaks when implementation changes
}

// ✅ GOOD - Testing what it does
func TestUserLogin(t *testing.T) {
    // Tests that valid credentials succeed
    // Tests that invalid credentials fail
    // Tests that session is created
    // Resilient to implementation changes
}
```

### 2. Flaky Tests

```go
// ❌ BAD - Depends on timing
func TestCache(t *testing.T) {
    cache.Set("key", "value", 100*time.Millisecond)
    time.Sleep(50*time.Millisecond) // Might fail under load
    if cache.Get("key") == nil {
        t.Error("Key should still be cached")
    }
}

// ✅ GOOD - Explicit control
func TestCache(t *testing.T) {
    cache.Set("key", "value", 1*time.Hour)
    
    if cache.Get("key") == nil {
        t.Error("Key should be cached")
    }
    
    cache.Delete("key")
    if cache.Get("key") != nil {
        t.Error("Key should be removed")
    }
}
```

### 3. Mega Tests

```go
// ❌ BAD - Tests everything in one function
func TestUserManagement(t *testing.T) {
    // Tests user creation
    // Tests user update
    // Tests user deletion
    // Tests user search
    // Tests user permissions
    // 500 lines of test code
}

// ✅ GOOD - Focused tests
func TestUserCreation(t *testing.T) {
    t.Run("with valid data", func(t *testing.T) {})
    t.Run("with duplicate email", func(t *testing.T) {})
}

func TestUserUpdate(t *testing.T) {
    t.Run("updates email successfully", func(t *testing.T) {})
    t.Run("rejects invalid email", func(t *testing.T) {})
}
```

### 4. Ignoring Errors

```go
// ❌ BAD - Ignoring errors
func TestSomething(t *testing.T) {
    db, _ := Connect() // Ignores error
    user, _ := CreateUser(db, "test@example.com") // Ignores error
    // Test continues with potentially nil values
}

// ✅ GOOD - Handling errors
func TestSomething(t *testing.T) {
    db, err := Connect()
    if err != nil {
        t.Fatalf("Failed to connect: %v", err)
    }
    
    user, err := CreateUser(db, "test@example.com")
    if err != nil {
        t.Fatalf("Failed to create user: %v", err)
    }
    
    // Test continues safely
}
```

---

## Examples

### Complete Test Example

```go
package auth_test

import (
    "context"
    "fmt"
    "testing"
    "time"
    
    "github.com/yourusername/laptop-tracking-system/internal/auth"
    "github.com/yourusername/laptop-tracking-system/internal/database"
    "github.com/yourusername/laptop-tracking-system/internal/models"
)

func TestSessionManagement(t *testing.T) {
    // Skip if running in short mode (unit tests only)
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Setup test database
    db, cleanup := database.SetupTestDB(t)
    defer cleanup()
    
    // Create test context
    ctx := context.Background()
    
    // Capture reference time
    now := time.Now()
    
    t.Run("create session successfully", func(t *testing.T) {
        // Create test user
        user := createTestUser(t, db, uniqueEmail())
        
        // Create session
        session, err := auth.CreateSession(ctx, db, user.ID, 24)
        if err != nil {
            t.Fatalf("Failed to create session: %v", err)
        }
        
        // Verify session properties
        if session.Token == "" {
            t.Error("Session token is empty")
        }
        if session.UserID != user.ID {
            t.Errorf("Expected user ID %d, got %d", user.ID, session.UserID)
        }
        if session.ExpiresAt.Before(now) {
            t.Error("Session expires in the past")
        }
    })
    
    t.Run("validate existing session", func(t *testing.T) {
        // Create test user and session
        user := createTestUser(t, db, uniqueEmail())
        session, _ := auth.CreateSession(ctx, db, user.ID, 24)
        
        // Validate session
        validatedUser, err := auth.ValidateSession(ctx, db, session.Token)
        if err != nil {
            t.Fatalf("Failed to validate session: %v", err)
        }
        
        if validatedUser.ID != user.ID {
            t.Errorf("Expected user ID %d, got %d", user.ID, validatedUser.ID)
        }
    })
    
    t.Run("reject expired session", func(t *testing.T) {
        // Create test user
        user := createTestUser(t, db, uniqueEmail())
        
        // Create expired session manually
        expiredSession := &models.Session{
            UserID:    user.ID,
            Token:     "expired-token-" + fmt.Sprintf("%d", now.UnixNano()),
            ExpiresAt: now.Add(-1 * time.Hour), // Expired
        }
        expiredSession.BeforeCreate()
        
        err := db.QueryRowContext(ctx,
            `INSERT INTO sessions (user_id, token, expires_at, created_at)
            VALUES ($1, $2, $3, $4) RETURNING id`,
            expiredSession.UserID, expiredSession.Token, 
            expiredSession.ExpiresAt, expiredSession.CreatedAt,
        ).Scan(&expiredSession.ID)
        
        if err != nil {
            t.Fatalf("Failed to create expired session: %v", err)
        }
        
        // Try to validate expired session
        _, err = auth.ValidateSession(ctx, db, expiredSession.Token)
        if err == nil {
            t.Error("Expected error for expired session, got nil")
        }
    })
}

// Helper functions
func uniqueEmail() string {
    return fmt.Sprintf("test_%d@example.com", time.Now().UnixNano())
}

func createTestUser(t *testing.T, db *sql.DB, email string) *models.User {
    t.Helper()
    
    user := &models.User{
        Email:        email,
        PasswordHash: "hashedpassword",
        Role:         models.RoleLogistics,
    }
    user.BeforeCreate()
    
    err := db.QueryRowContext(
        context.Background(),
        `INSERT INTO users (email, password_hash, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5) RETURNING id`,
        user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
    ).Scan(&user.ID)
    
    if err != nil {
        t.Fatalf("Failed to create test user: %v", err)
    }
    
    return user
}
```

---

## Running Tests

### Basic Commands

```bash
# Run all tests
make test-all

# Run unit tests only
make test-unit

# Run specific package
make test-package PKG=internal/auth

# Run with coverage
make test-coverage
```

### Test Flags

```bash
# Verbose output
go test -v ./...

# Run specific test
go test -v -run TestUserLogin ./internal/auth

# Short mode (skip slow tests)
go test -short ./...

# Show coverage
go test -cover ./...

# Race detection
go test -race ./...

# Parallel execution control
go test -p=1 ./...  # Sequential
go test -p=4 ./...  # 4 packages in parallel
```

---

## Checklist for New Tests

Before submitting a test, ensure:

- [ ] Test has a descriptive name
- [ ] Test uses subtests for organization
- [ ] Test data is unique (no collisions)
- [ ] Time references are captured once
- [ ] Database cleanup is handled with defer
- [ ] Errors are checked and handled
- [ ] Test is focused (single responsibility)
- [ ] Test is deterministic (no randomness or timing issues)
- [ ] Test passes in isolation
- [ ] Test passes in the full suite
- [ ] Test includes both happy path and error cases

---

## Additional Resources

- **Official Go Testing Docs:** https://pkg.go.dev/testing
- **Table-Driven Tests:** https://dave.cheney.net/2019/05/07/prefer-table-driven-tests
- **Test Fixtures:** https://github.com/go-testfixtures/testfixtures
- **Project Docs:** `docs/TEST_DATABASE_SETUP.md`

---

**Maintained By:** Development Team  
**Questions?** See `docs/TEST_FAILURES_RESOLVED.md` for real-world examples


