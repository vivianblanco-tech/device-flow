# Email Testing with Docker

This guide explains how to run email notification tests inside Docker containers.

## Quick Start

### Prerequisites

- Docker Desktop installed and running
- Docker Compose (comes with Docker Desktop)

### Run Tests (Easiest)

**Windows (PowerShell):**
```powershell
.\scripts\test_emails_docker.ps1
```

**Linux/macOS:**
```bash
chmod +x scripts/test_emails_docker.sh
./scripts/test_emails_docker.sh
```

That's it! The script will:
1. ✅ Start PostgreSQL and Mailhog containers
2. ✅ Wait for services to be ready
3. ✅ Build the test container
4. ✅ Run all 6 email notification tests
5. ✅ Show detailed results
6. ✅ Clean up test data

## Usage Options

### Default Mode (Uses Dev Database)

```powershell
# Windows
.\scripts\test_emails_docker.ps1

# Linux/macOS
./scripts/test_emails_docker.sh
```

Tests run against the `laptop_tracking_dev` database. **Note:** This will create and clean up test data in your dev database.

### Isolated Mode (Recommended)

```powershell
# Windows
.\scripts\test_emails_docker.ps1 --isolated

# Linux/macOS
./scripts/test_emails_docker.sh --isolated
```

Creates a **separate test database** (`laptop_tracking_test`) that is completely isolated from your development data. **Recommended for safety!**

### Keep Containers Running

```powershell
# Windows
.\scripts\test_emails_docker.ps1 --no-cleanup

# Linux/macOS
./scripts/test_emails_docker.sh --no-cleanup
```

Leaves containers running after tests complete. Useful for debugging.

## Manual Docker Commands

If you prefer to run Docker commands directly:

### 1. Start Services

```bash
docker-compose up -d postgres mailhog
```

### 2. Build Test Container

```bash
# Default mode
docker-compose -f docker-compose.yml -f docker-compose.test.yml build email-test

# Isolated mode
docker-compose -f docker-compose.yml -f docker-compose.test.yml build email-test-isolated
```

### 3. Run Tests

```bash
# Default mode
docker-compose -f docker-compose.yml -f docker-compose.test.yml --profile test run --rm email-test

# Isolated mode (with separate test DB)
docker-compose -f docker-compose.yml -f docker-compose.test.yml --profile test-isolated up -d postgres-test
docker-compose -f docker-compose.yml -f docker-compose.test.yml --profile test-isolated run --rm email-test-isolated
```

### 4. View Emails

Open Mailhog: http://localhost:8025

### 5. Cleanup

```bash
# Stop all services
docker-compose down

# Remove test database volume (isolated mode)
docker-compose -f docker-compose.yml -f docker-compose.test.yml --profile test-isolated down -v
```

## How It Works

### Docker Networking

The test container runs inside the Docker network and connects to services using their container names:

- **Database**: `postgres` (instead of `localhost`)
- **Mailhog**: `mailhog` (instead of `localhost`)

Environment variables are automatically configured:
```bash
DATABASE_URL=postgres://postgres:password@postgres:5432/laptop_tracking_dev?sslmode=disable
MAILHOG_URL=http://mailhog:8025
SMTP_HOST=mailhog
SMTP_PORT=1025
```

### Container Architecture

```
┌─────────────────────────────────────────┐
│          Docker Network                 │
│                                         │
│  ┌──────────────┐   ┌──────────────┐   │
│  │  PostgreSQL  │   │   Mailhog    │   │
│  │  Container   │   │  Container   │   │
│  │  (postgres)  │   │  (mailhog)   │   │
│  └──────────────┘   └──────────────┘   │
│         ▲                  ▲            │
│         │                  │            │
│         └──────┬───────────┘            │
│                │                        │
│      ┌─────────▼────────┐               │
│      │   Test Container │               │
│      │   (email-test)   │               │
│      └──────────────────┘               │
│                                         │
└─────────────────────────────────────────┘
             │
             ▼
      Your Browser
   http://localhost:8025
```

## Files Created

| File | Purpose |
|------|---------|
| `docker-compose.test.yml` | Test service definitions |
| `Dockerfile.test` | Test container build configuration |
| `scripts/test_emails_docker.ps1` | PowerShell wrapper script |
| `scripts/test_emails_docker.sh` | Bash wrapper script |

## Expected Output

```powershell
PS> .\scripts\test_emails_docker.ps1

========================================
  Email Testing - Docker Environment
========================================

✅ Docker is running
✅ docker-compose is available

Starting required services...
✅ Services started

Waiting for services to be ready...
✅ Services are ready

Building test container...
Using development database
✅ Test container built

========================================
  Running Email Notification Tests
========================================

==================================================
  Email Notifications Test - Mailhog Verification
==================================================

Configuration:
  Database: postgres://postgres:password@postgres:5432/laptop_tracking_dev?sslmode=disable
  Mailhog:  http://mailhog:8025
  SMTP:     mailhog:1025

✅ Connected to database
✅ Connected to Mailhog
✅ Cleared Mailhog messages
✅ Test data created

📧 Test 1: Magic Link Email
  ✅ SUCCESS
     Subject:   Access Your Form - pickup
     Recipient: test@example.com

📧 Test 2: Pickup Confirmation
  ✅ SUCCESS
     Subject:   Pickup Confirmation - CONF-1
     Recipient: noreply@example.com

... (4 more tests)

===========================================
  Test Summary
===========================================

✅ Magic Link
✅ Pickup Confirmation
✅ Pickup Scheduled
✅ Warehouse Pre-Alert
✅ Release Notification
✅ Delivery Confirmation

Total Tests: 6
Passed: 6
Failed: 0

✅ Test data cleaned up

🎉 All email notifications passed!

========================================
  All Tests Passed!
========================================

💡 View emails in Mailhog:
   http://localhost:8025

Cleaning up...
✅ Cleanup complete
```

## Troubleshooting

### Docker Not Running

**Error:** `❌ Docker is not running`

**Solution:** Start Docker Desktop

### Port Already in Use

**Error:** `Bind for 0.0.0.0:5432 failed: port is already allocated`

**Solution:** 
1. Stop local PostgreSQL: `sudo service postgresql stop`
2. Or change port in `docker-compose.yml`

### Services Not Ready

**Error:** `⚠️ Services may not be fully ready`

**Solution:** Wait a few more seconds and try again, or check Docker logs:
```bash
docker-compose logs postgres
docker-compose logs mailhog
```

### Build Failures

**Error:** `❌ Failed to build test container`

**Solution:**
1. Check Docker disk space: `docker system df`
2. Clean up: `docker system prune`
3. Check logs: `docker-compose logs`

### Test Container Exits Immediately

**Error:** Tests don't run, container exits

**Solution:**
```bash
# Check logs
docker logs laptop-tracking-email-test

# Run interactively for debugging
docker-compose -f docker-compose.yml -f docker-compose.test.yml --profile test run --rm email-test
```

### Cannot Connect to Services

**Error:** `Failed to connect to database` or `Failed to connect to Mailhog`

**Solution:**
1. Verify services are running:
   ```bash
   docker ps
   ```
2. Check service health:
   ```bash
   docker inspect laptop-tracking-db
   docker inspect laptop-tracking-mailhog
   ```
3. Verify network connectivity:
   ```bash
   docker network ls
   docker network inspect bdh_default
   ```

## Comparison: Docker vs Local

| Aspect | Docker | Local |
|--------|--------|-------|
| **Setup** | Automatic | Manual (Mailhog, PostgreSQL) |
| **Isolation** | Complete | Shared with host |
| **Networking** | Container names | localhost |
| **Speed** | Slightly slower (build time) | Faster |
| **CI/CD** | Perfect | Requires setup |
| **Debugging** | Container logs | Direct |
| **Cleanup** | Automatic | Manual |

## When to Use Docker Testing

✅ **Use Docker when:**
- Running tests in CI/CD
- Want isolated environment
- Don't want to install Mailhog locally
- Testing on different OS
- Want reproducible results

❌ **Use Local when:**
- Rapid development/debugging
- Need to inspect email rendering in real-time
- Already have local services running
- Want faster iteration

## Advanced Usage

### Run Specific Tests

Modify `scripts/email-test/main.go` to comment out tests you don't want to run, then rebuild:

```bash
docker-compose -f docker-compose.yml -f docker-compose.test.yml build email-test
```

### Interactive Debugging

Get a shell inside the test container:

```bash
docker-compose -f docker-compose.yml -f docker-compose.test.yml --profile test run --rm --entrypoint /bin/sh email-test

# Inside container
/app # psql $DATABASE_URL
/app # curl $MAILHOG_URL/api/v2/messages
```

### Custom Environment Variables

Create `.env.test` file:

```env
DATABASE_URL=postgres://postgres:password@postgres:5432/my_test_db?sslmode=disable
SMTP_FROM=custom@test.com
```

Then load it:

```bash
docker-compose --env-file .env.test -f docker-compose.yml -f docker-compose.test.yml --profile test run --rm email-test
```

### Keep Test Database

```bash
# Start isolated test DB
docker-compose -f docker-compose.yml -f docker-compose.test.yml --profile test-isolated up -d postgres-test

# Run multiple test iterations
docker-compose -f docker-compose.yml -f docker-compose.test.yml --profile test-isolated run --rm email-test-isolated

# Inspect database
docker exec -it laptop-tracking-db-test psql -U postgres -d laptop_tracking_test

# Stop when done (keeps data)
docker-compose -f docker-compose.yml -f docker-compose.test.yml --profile test-isolated stop

# Remove completely (deletes data)
docker-compose -f docker-compose.yml -f docker-compose.test.yml --profile test-isolated down -v
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Email Notification Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Run Email Tests
        run: |
          chmod +x scripts/test_emails_docker.sh
          ./scripts/test_emails_docker.sh --isolated
      
      - name: Upload Mailhog Screenshots (on failure)
        if: failure()
        uses: actions/upload-artifact@v3
        with:
          name: mailhog-logs
          path: /tmp/mailhog-*.log
```

### GitLab CI Example

```yaml
test:email-notifications:
  image: docker:latest
  services:
    - docker:dind
  script:
    - chmod +x scripts/test_emails_docker.sh
    - ./scripts/test_emails_docker.sh --isolated
  artifacts:
    when: on_failure
    paths:
      - logs/
```

## Related Documentation

- Main Email Testing: `scripts/README.md`
- Quick Start: `scripts/EMAIL_TESTING_QUICKSTART.md`
- Docker Compose: `docker-compose.yml`
- Email Implementation: `internal/email/`

