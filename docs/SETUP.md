# Development Environment Setup Guide

This guide will help you set up your development environment for the Laptop Tracking System.

## Prerequisites Installation

### 1. Install Go

**Windows:**
1. Download Go from https://golang.org/dl/
2. Run the installer
3. Verify installation: `go version`

**macOS:**
```bash
brew install go
```

**Linux:**
```bash
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### 2. Install PostgreSQL

**Windows:**
1. Download from https://www.postgresql.org/download/windows/
2. Run the installer
3. Remember the password you set for the `postgres` user
4. Add PostgreSQL bin directory to PATH

**macOS:**
```bash
brew install postgresql@15
brew services start postgresql@15
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

### 3. Install golang-migrate

**Windows (PowerShell as Administrator):**
```powershell
# Download the latest release
$version = "v4.17.0"
$url = "https://github.com/golang-migrate/migrate/releases/download/$version/migrate.windows-amd64.zip"
Invoke-WebRequest -Uri $url -OutFile migrate.zip
Expand-Archive migrate.zip -DestinationPath C:\migrate
Move-Item C:\migrate\migrate.exe C:\Windows\System32\
Remove-Item migrate.zip
Remove-Item C:\migrate -Recurse
```

**macOS:**
```bash
brew install golang-migrate
```

**Linux:**
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

Verify installation:
```bash
migrate -version
```

### 4. Install Mailhog (Email Testing)

**Windows:**
```powershell
# Download from GitHub releases
# https://github.com/mailhog/MailHog/releases
# Or use Docker (recommended)
docker run -d -p 1025:1025 -p 8025:8025 mailhog/mailhog
```

**macOS:**
```bash
brew install mailhog
```

**Linux:**
```bash
go install github.com/mailhog/MailHog@latest
```

**Using Docker (All Platforms):**
```bash
docker run -d -p 1025:1025 -p 8025:8025 --name mailhog mailhog/mailhog
```

### 5. Install Make (Optional but Recommended)

**Windows:**
```powershell
# Using Chocolatey
choco install make

# Or using Scoop
scoop install make
```

**macOS:**
```bash
# Already installed via Xcode Command Line Tools
xcode-select --install
```

**Linux:**
```bash
sudo apt install make
```

## Project Setup

### 1. Clone and Configure

```bash
# Navigate to project directory
cd laptop-tracking-system

# Copy environment file
cp .env.example .env

# Edit .env with your settings
# Use your preferred text editor
notepad .env       # Windows
nano .env          # Linux/macOS
code .env          # VS Code
```

### 2. Configure Database

**Create Database:**

```bash
# Connect to PostgreSQL
psql -U postgres

# In psql prompt:
CREATE DATABASE laptop_tracking_dev;
\q
```

**Windows PowerShell:**
```powershell
# If psql is not in PATH, use full path:
& "C:\Program Files\PostgreSQL\15\bin\psql.exe" -U postgres -c "CREATE DATABASE laptop_tracking_dev;"
```

**Update .env file:**
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_password
DB_NAME=laptop_tracking_dev
DB_SSLMODE=disable
```

### 3. Install Dependencies

```bash
go mod download
go mod tidy
```

### 4. Run Migrations

**Using Make:**
```bash
make migrate-up
```

**Without Make (Windows PowerShell):**
```powershell
$env:DB_USER="postgres"
$env:DB_PASSWORD="your_password"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_NAME="laptop_tracking_dev"
$env:DB_SSLMODE="disable"

migrate -path migrations -database "postgres://${env:DB_USER}:${env:DB_PASSWORD}@${env:DB_HOST}:${env:DB_PORT}/${env:DB_NAME}?sslmode=${env:DB_SSLMODE}" up
```

**Without Make (Linux/macOS):**
```bash
source .env
migrate -path migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" up
```

### 5. Setup Google OAuth (Optional for now)

1. Visit [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project: "Laptop Tracking System"
3. Enable APIs: Google+ API
4. Create OAuth 2.0 Credentials:
   - Application type: Web application
   - Name: Laptop Tracking Dev
   - Authorized JavaScript origins: `http://localhost:8080`
   - Authorized redirect URIs: `http://localhost:8080/auth/google/callback`
5. Copy the Client ID and Client Secret
6. Update `.env`:
   ```env
   GOOGLE_CLIENT_ID=your_client_id_here
   GOOGLE_CLIENT_SECRET=your_client_secret_here
   GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
   ```

### 6. Start Mailhog

```bash
# If installed via brew/apt
mailhog

# If using Docker
docker start mailhog

# Or run new container
docker run -d -p 1025:1025 -p 8025:8025 mailhog/mailhog
```

Access Mailhog UI at: http://localhost:8025

### 7. Generate Secrets

For production, generate secure secrets:

```bash
# Linux/macOS
openssl rand -base64 48

# Windows PowerShell
[Convert]::ToBase64String((1..48 | ForEach-Object { Get-Random -Minimum 0 -Maximum 256 }))
```

Update in `.env`:
```env
SESSION_SECRET=your_generated_secret_here
CSRF_SECRET=your_other_generated_secret_here
```

## Running the Application

### Using Make

```bash
make run
```

### Without Make

```bash
go run cmd/web/main.go
```

### Building and Running Binary

```bash
# Build
go build -o bin/laptop-tracking cmd/web/main.go

# Run
./bin/laptop-tracking        # Linux/macOS
.\bin\laptop-tracking.exe    # Windows
```

## Verify Installation

1. Application should start without errors
2. Visit http://localhost:8080/health - should return "OK"
3. Check logs for "Database connected successfully"
4. Mailhog UI should be accessible at http://localhost:8025

## Common Issues

### PostgreSQL Connection Refused

**Symptom:** `connection refused` error

**Solutions:**
- Check PostgreSQL is running: `pg_isready`
- Start PostgreSQL: `brew services start postgresql@15` (macOS)
- Check port 5432 is not blocked
- Verify credentials in `.env`

### Migration Errors

**Symptom:** `migration failed` or `dirty database`

**Solutions:**
```bash
# Check migration version
migrate -path migrations -database "postgres://..." version

# Force to specific version
migrate -path migrations -database "postgres://..." force 1

# Or reset database
make db-reset
```

### Port Already in Use

**Symptom:** `address already in use: 8080`

**Solutions:**
```bash
# Find process using port (Windows)
netstat -ano | findstr :8080

# Find process using port (Linux/macOS)
lsof -i :8080

# Kill process or change APP_PORT in .env
```

### Go Module Issues

**Symptom:** `cannot find package`

**Solutions:**
```bash
# Clean module cache
go clean -modcache

# Reinstall dependencies
rm go.sum
go mod download
go mod tidy
```

## IDE Setup

### VS Code

Recommended extensions:
- Go (golang.go)
- PostgreSQL (ckolkman.vscode-postgres)
- GitLens (eamodio.gitlens)
- Error Lens (usernamehw.errorlens)

Settings (`.vscode/settings.json`):
```json
{
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "go.formatTool": "gofmt",
  "editor.formatOnSave": true,
  "[go]": {
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  }
}
```

### GoLand

1. Open project
2. GoLand will auto-detect Go modules
3. Configure PostgreSQL data source
4. Enable Go modules integration

## Next Steps

After setup is complete:
1. Run tests: `make test`
2. Check code formatting: `make fmt`
3. Read the main README.md for architecture overview
4. Review `plan.md` for development roadmap
5. Start with Phase 1: Database Schema & Core Models

## Getting Help

- Check `README.md` for general information
- Review `docs/` directory for more documentation
- Check logs in console output
- Verify all environment variables are set correctly

