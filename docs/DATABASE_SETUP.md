# Database Setup Guide

This guide provides multiple options for setting up the PostgreSQL database for the Laptop Tracking System.

---

## Option 1: Docker Compose (Recommended) ⭐

The easiest and most reliable method. No local PostgreSQL installation required.

### Prerequisites
- Docker Desktop installed and running
- Docker Compose available

### Steps

1. **Start PostgreSQL with Docker Compose:**
```powershell
# Start only PostgreSQL
docker-compose up -d postgres

# Or start all services (PostgreSQL + MailHog)
docker-compose up -d
```

2. **Verify PostgreSQL is running:**
```powershell
docker ps
# Should show 'laptop-tracking-db' container running
```

3. **Create .env file:**
```powershell
# Copy the example file
Copy-Item .env.example .env

# The default settings will work with Docker Compose
```

4. **Run migrations:**
```powershell
# Wait a few seconds for PostgreSQL to fully start, then:
make migrate-up
```

Or without Make:
```powershell
$env:DB_USER="postgres"
$env:DB_PASSWORD="postgres"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_NAME="laptop_tracking_dev"

migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/laptop_tracking_dev?sslmode=disable" up
```

5. **Create test database:**
```powershell
# Connect to PostgreSQL container
docker exec -it laptop-tracking-db psql -U postgres

# In psql prompt:
CREATE DATABASE laptop_tracking_test;
\q
```

Or in one command:
```powershell
docker exec -it laptop-tracking-db psql -U postgres -c "CREATE DATABASE laptop_tracking_test;"
```

6. **Run test migrations:**
```powershell
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/laptop_tracking_test?sslmode=disable" up
```

### Managing Docker Database

**Stop database:**
```powershell
docker-compose stop postgres
```

**Start database:**
```powershell
docker-compose start postgres
```

**Remove database (WARNING: deletes all data):**
```powershell
docker-compose down -v
```

**View logs:**
```powershell
docker-compose logs postgres
```

**Access psql directly:**
```powershell
docker exec -it laptop-tracking-db psql -U postgres -d laptop_tracking_dev
```

---

## Option 2: Automated Setup Script (Windows)

Use the provided PowerShell script for automatic setup.

### Prerequisites
- PostgreSQL installed locally
- `psql` and `migrate` commands available in PATH

### Steps

1. **Run the setup script:**
```powershell
.\scripts\setup-database.ps1
```

The script will:
- ✓ Check for PostgreSQL installation
- ✓ Create development and test databases
- ✓ Generate .env file from .env.example
- ✓ Run all migrations
- ✓ Provide next steps

---

## Option 3: Manual Setup (Local PostgreSQL)

For manual control over the setup process.

### Prerequisites
- PostgreSQL 15+ installed locally
- `psql` command available in PATH

### Steps

1. **Verify PostgreSQL is running:**
```powershell
# Check if PostgreSQL service is running
Get-Service -Name postgresql*

# Or try connecting
psql -U postgres -c "SELECT version();"
```

2. **Create databases:**
```powershell
# Development database
psql -U postgres -c "CREATE DATABASE laptop_tracking_dev;"

# Test database
psql -U postgres -c "CREATE DATABASE laptop_tracking_test;"
```

3. **Create .env file:**
```powershell
# Copy example
Copy-Item .env.example .env

# Edit with your password
notepad .env
```

Update these values in `.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_actual_password
DB_NAME=laptop_tracking_dev
DB_SSLMODE=disable

TEST_DATABASE_URL=postgres://postgres:your_actual_password@localhost:5432/laptop_tracking_test?sslmode=disable
```

4. **Install golang-migrate (if not installed):**
```powershell
# Using Go
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Or download binary from:
# https://github.com/golang-migrate/migrate/releases
```

5. **Run migrations:**
```powershell
# Load environment variables
Get-Content .env | ForEach-Object {
    if ($_ -match '^([^=]+)=(.*)$') {
        [System.Environment]::SetEnvironmentVariable($matches[1], $matches[2], 'Process')
    }
}

# Run migrations for development database
migrate -path migrations -database "postgres://${env:DB_USER}:${env:DB_PASSWORD}@${env:DB_HOST}:${env:DB_PORT}/${env:DB_NAME}?sslmode=disable" up

# Run migrations for test database
migrate -path migrations -database "${env:TEST_DATABASE_URL}" up
```

Or using Make:
```powershell
make migrate-up
```

---

## Verification

After setup, verify everything is working:

### 1. Test database connection:
```powershell
psql -U postgres -d laptop_tracking_dev -c "SELECT COUNT(*) FROM users;"
```

### 2. Check migrations:
```powershell
psql -U postgres -d laptop_tracking_dev -c "\dt"
```

You should see these tables:
- users
- client_companies
- software_engineers
- laptops
- shipments
- shipment_laptops
- pickup_forms
- reception_reports
- delivery_forms
- magic_links
- sessions
- notification_logs
- audit_logs
- schema_migrations

### 3. Run tests:
```powershell
# Run all tests including integration tests
go test ./... -v

# Run only unit tests (no database)
go test ./... -v -short
```

### 4. Start the application:
```powershell
go run cmd/web/main.go
```

---

## Common Issues & Solutions

### Issue: "psql: command not found"

**Solution:**
- Add PostgreSQL bin directory to PATH
- Windows: `C:\Program Files\PostgreSQL\15\bin`
- Or use Docker option instead

### Issue: "password authentication failed"

**Solution:**
1. Verify password in `.env` matches PostgreSQL password
2. Reset PostgreSQL password:
```powershell
# Windows (as Administrator)
psql -U postgres
ALTER USER postgres WITH PASSWORD 'new_password';
```

### Issue: "database already exists"

**Solution:**
```powershell
# Drop and recreate (WARNING: loses all data)
psql -U postgres -c "DROP DATABASE laptop_tracking_dev;"
psql -U postgres -c "CREATE DATABASE laptop_tracking_dev;"
```

### Issue: "migration failed: dirty database"

**Solution:**
```powershell
# Check current migration version
migrate -path migrations -database "your_db_url" version

# Force to a specific version (e.g., 10)
migrate -path migrations -database "your_db_url" force 10

# Then run migrations again
make migrate-up
```

### Issue: Docker port 5432 already in use

**Solution:**
```powershell
# Stop local PostgreSQL service
Stop-Service postgresql*

# Or change Docker port in docker-compose.yml:
ports:
  - "5433:5432"  # Use port 5433 instead
```

---

## Database Management Commands

### Reset database (WARNING: deletes all data):
```powershell
make db-reset
```

Or manually:
```powershell
psql -U postgres -c "DROP DATABASE IF EXISTS laptop_tracking_dev;"
psql -U postgres -c "CREATE DATABASE laptop_tracking_dev;"
make migrate-up
```

### Rollback last migration:
```powershell
make migrate-down
```

### Check migration status:
```powershell
migrate -path migrations -database "your_db_url" version
```

### View database size:
```powershell
psql -U postgres -d laptop_tracking_dev -c "SELECT pg_size_pretty(pg_database_size('laptop_tracking_dev'));"
```

---

## Environment Variables Reference

Required variables for database connection:

| Variable | Description | Example |
|----------|-------------|---------|
| `DB_HOST` | Database host | `localhost` or `postgres` (Docker) |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `your_password` |
| `DB_NAME` | Database name | `laptop_tracking_dev` |
| `DB_SSLMODE` | SSL mode | `disable` (dev) or `require` (prod) |
| `TEST_DATABASE_URL` | Full test DB URL | `postgres://...` |

---

## Next Steps

After successful database setup:

1. ✅ Verify all migrations are applied
2. ✅ Run the test suite: `go test ./... -v`
3. ✅ Review test results in `docs/TEST_STATUS_SUMMARY.md`
4. ✅ Start the application: `go run cmd/web/main.go`
5. ✅ Access the application at http://localhost:8080

---

## Additional Resources

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- Project Setup Guide: `docs/SETUP.md`
- Test Status: `docs/TEST_STATUS_SUMMARY.md`

---

**Need Help?**

If you encounter issues not covered here:
1. Check PostgreSQL logs
2. Verify connection with `psql`
3. Ensure all prerequisites are installed
4. Try the Docker option (usually most reliable)















