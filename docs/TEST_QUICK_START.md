# Test Suite Quick Start Guide

## One-Line Setup & Run

**Windows PowerShell:**
```powershell
$env:TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"; go test ./... -p=1 -v -race
```

**Linux/Mac (Bash):**
```bash
export TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable" && go test ./... -p=1 -v -race
```

## Prerequisites Checklist

- [ ] Docker containers running (`docker-compose up -d`)
- [ ] PostgreSQL accessible at `localhost:5432`
- [ ] Password matches Docker config (default: `password`)
- [ ] Test database exists (`laptop_tracking_test`)

## Critical Requirements

1. **MUST run sequentially**: Always use `-p=1` flag
2. **MUST set environment variable**: `TEST_DATABASE_URL` before running
3. **Password must match**: Docker PostgreSQL password

## Common Commands

```bash
# Full suite
go test ./... -p=1 -v -race

# Specific package
go test ./internal/handlers -p=1 -v -race

# Specific test
go test ./internal/handlers -p=1 -v -race -run TestName

# With coverage
go test ./... -p=1 -v -race -cover
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Password auth failed | Check Docker password, set `TEST_DATABASE_URL` |
| Race conditions | Use `-p=1` flag |
| NULL constraint errors | Include all required fields in test data |
| Date validation fails | Use dynamic future dates |

See `docs/TEST_RUN_INSTRUCTIONS.md` for detailed information.

