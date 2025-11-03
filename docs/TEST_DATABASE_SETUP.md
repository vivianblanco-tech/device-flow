# Test Database Setup Guide

This guide will help you set up the test database required for running integration tests.

## Overview

The test suite includes both unit tests (which don't require a database) and integration tests (which do). Currently, 40 integration tests cannot run because the test database is not configured.

## Quick Setup (Recommended)

### Windows (PowerShell)

```powershell
# 1. Create test database
psql -U postgres -c "CREATE DATABASE laptop_tracking_test;"

# 2. Set password for postgres user (if needed)
psql -U postgres -c "ALTER USER postgres PASSWORD 'password';"

# 3. Run migrations on test database
$env:DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"
migrate -path migrations -database $env:DATABASE_URL up

# 4. Set test database URL (add to .env or set as environment variable)
$env:TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"

# 5. Run all tests including integration tests
go test ./...
```

### Linux/macOS (Bash)

```bash
# 1. Create test database
psql -U postgres -c "CREATE DATABASE laptop_tracking_test;"

# 2. Set password for postgres user (if needed)
psql -U postgres -c "ALTER USER postgres PASSWORD 'password';"

# 3. Run migrations on test database
export DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"
migrate -path migrations -database $DATABASE_URL up

# 4. Set test database URL (add to .env or set as environment variable)
export TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"

# 5. Run all tests including integration tests
go test ./...
```

## Docker-based Setup (Alternative)

If you're using Docker Compose:

```powershell
# Start the database container
docker-compose up -d postgres

# Create test database inside container
docker-compose exec postgres psql -U postgres -c "CREATE DATABASE laptop_tracking_test;"

# Run migrations
$env:DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"
migrate -path migrations -database $env:DATABASE_URL up

# Set environment variable and run tests
$env:TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"
go test ./...
```

## Configuration Details

### Environment Variables

The test suite uses the following environment variable:

- `TEST_DATABASE_URL`: Connection string for test database
  - Default: `postgres://postgres:postgres@localhost:5432/laptop_tracking_test?sslmode=disable`
  - Recommended: `postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable`

### .env File (Optional)

You can add the test database URL to your `.env` file:

```env
# Test Database Configuration
TEST_DATABASE_URL=postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable
```

Note: The `.env` file is primarily for the main application. Tests use environment variables directly.

## Verifying Setup

### 1. Check Database Exists

```powershell
psql -U postgres -l | Select-String "laptop_tracking_test"
```

Expected output should show `laptop_tracking_test` database.

### 2. Check Migrations Applied

```powershell
psql -U postgres -d laptop_tracking_test -c "\dt"
```

Expected output should show 13 tables:
- users
- client_companies
- software_engineers
- laptops
- shipments
- shipment_laptops
- pickup_forms
- reception_reports
- delivery_forms
- sessions
- magic_links
- notification_logs
- audit_logs

### 3. Run Tests

```powershell
# Set environment variable
$env:TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"

# Run specific package tests
go test ./internal/auth -v
go test ./internal/handlers -v
go test ./internal/database -v

# Run all tests
go test ./... -v
```

### Expected Results

After setup, you should see:
- ✅ All 23 auth tests passing
- ✅ All 15 handler tests passing
- ✅ All 2 database tests passing
- ✅ All 33 email tests passing
- ✅ Total: ~254 tests passing

## Troubleshooting

### Issue: "pq: password authentication failed for user postgres"

**Solution**: Your PostgreSQL user password doesn't match. Update the password:

```powershell
psql -U postgres -c "ALTER USER postgres PASSWORD 'password';"
```

Or update `TEST_DATABASE_URL` to use your actual password.

### Issue: "pq: database 'laptop_tracking_test' does not exist"

**Solution**: Create the database:

```powershell
psql -U postgres -c "CREATE DATABASE laptop_tracking_test;"
```

### Issue: Tests fail with "table does not exist"

**Solution**: Run migrations on test database:

```powershell
$env:DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"
migrate -path migrations -database $env:DATABASE_URL up
```

### Issue: PostgreSQL is not running

**Solution**: Start PostgreSQL service:

**Windows**:
```powershell
Start-Service postgresql-x64-15
```

**Linux**:
```bash
sudo systemctl start postgresql
```

**macOS**:
```bash
brew services start postgresql
```

**Docker**:
```bash
docker-compose up -d postgres
```

### Issue: migrate command not found

**Solution**: Install golang-migrate:

**Windows**:
Download from: https://github.com/golang-migrate/migrate/releases

**macOS**:
```bash
brew install golang-migrate
```

**Linux**:
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

## Test Database Best Practices

1. **Separate from Development Database**: Always use a separate database for tests to avoid data corruption.

2. **Clean State**: Tests automatically clean up data after each run using the cleanup functions in `internal/database/testhelpers.go`.

3. **Parallel Tests**: Some tests cannot run in parallel due to database state. Run with `-p 1` if needed:
   ```powershell
   go test ./... -p 1
   ```

4. **CI/CD Integration**: In CI/CD pipelines, create and migrate the test database before running tests:
   ```yaml
   - name: Setup test database
     run: |
       psql -U postgres -c "CREATE DATABASE laptop_tracking_test;"
       migrate -path migrations -database $DATABASE_URL up
   ```

5. **Reset Test Database**: If tests are failing due to corrupted data, reset the database:
   ```powershell
   psql -U postgres -c "DROP DATABASE IF EXISTS laptop_tracking_test;"
   psql -U postgres -c "CREATE DATABASE laptop_tracking_test;"
   migrate -path migrations -database $env:DATABASE_URL up
   ```

## Makefile Commands (Future Enhancement)

Consider adding these to the Makefile:

```makefile
.PHONY: test-db-create
test-db-create:
	@echo "Creating test database..."
	psql -U postgres -c "CREATE DATABASE laptop_tracking_test;"

.PHONY: test-db-migrate
test-db-migrate:
	@echo "Running migrations on test database..."
	migrate -path migrations -database "$(TEST_DATABASE_URL)" up

.PHONY: test-db-reset
test-db-reset:
	@echo "Resetting test database..."
	psql -U postgres -c "DROP DATABASE IF EXISTS laptop_tracking_test;"
	psql -U postgres -c "CREATE DATABASE laptop_tracking_test;"
	migrate -path migrations -database "$(TEST_DATABASE_URL)" up

.PHONY: test-integration
test-integration: test-db-create test-db-migrate
	@echo "Running integration tests..."
	TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test ./... -v
```

## Summary

Once the test database is set up, all 254 tests should pass, giving you confidence in the codebase before moving to Phase 6 (Dashboard & Visualization).

**Estimated Setup Time**: 5-10 minutes

---

**Last Updated**: November 3, 2025  
**Related**: See `docs/PROJECT_STATUS.md` for overall project status

