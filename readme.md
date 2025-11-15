# Laptop Tracking System

A comprehensive web application for tracking laptop pickup and delivery from client companies to software engineers, including inventory management and status tracking.

## Features

- 📦 **Shipment Tracking**: Track laptops from client pickup through delivery to software engineers
- 👥 **Multi-Role System**: Support for Logistics, Client, Warehouse, and Project Manager roles
- 🔐 **Dual Authentication**: Username/password and Google OAuth (restricted to @bairesdev.com)
- 🔗 **Magic Links**: Secure one-time access links for form submissions
- 📧 **Email Notifications**: Automated notifications at each process step
- 🎫 **JIRA Integration**: Sync with JIRA tickets for seamless workflow
- 📊 **Dashboard & Analytics**: Real-time statistics and visualization
- 📅 **Calendar View**: Track pickup and delivery schedules
- 📸 **Photo Uploads**: Document device condition at pickup and delivery
- 🔍 **Inventory Management**: Track device serial numbers, availability, and software engineer assignments

## Tech Stack

- **Backend**: Go 1.22+
- **Database**: PostgreSQL 15+
- **Frontend**: HTML templates with Tailwind CSS v4
- **Charts**: Chart.js v4.4.1 (Line, Donut, Bar charts)
- **Authentication**: Google OAuth 2.0, bcrypt for passwords
- **Email**: SMTP (Mailhog for development)
- **Migrations**: golang-migrate
- **Deployment**: Docker, Caddy (reverse proxy with automatic SSL)

## Project Structure

```
.
├── cmd/
│   └── web/                 # Main application entry point
├── internal/                # Private application code
│   ├── auth/               # Authentication logic
│   ├── config/             # Configuration management
│   ├── database/           # Database connection and utilities
│   ├── email/              # Email service
│   ├── handlers/           # HTTP request handlers
│   ├── jira/               # JIRA integration
│   ├── middleware/         # HTTP middleware
│   ├── models/             # Data models
│   └── validator/          # Form validation
├── migrations/             # Database migrations
├── static/                 # Static assets (CSS, JS, images)
├── templates/              # HTML templates
│   ├── layouts/           # Base layouts
│   ├── pages/             # Page templates
│   └── components/        # Reusable components
├── tests/                  # Tests
│   ├── unit/              # Unit tests
│   ├── integration/       # Integration tests
│   └── e2e/               # End-to-end tests
├── uploads/                # Uploaded files (photos)
├── .env.example            # Example environment variables
├── go.mod                  # Go module definition
├── Makefile                # Common commands
└── README.md               # This file
```

## Prerequisites

- Go 1.22 or higher
- PostgreSQL 15 or higher
- Make (optional, for convenience commands)
- Mailhog (for email testing in development)
- golang-migrate CLI tool

## Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd laptop-tracking-system
```

### 2. Install Dependencies

```bash
# Install Go dependencies
make install

# Or manually
go mod download
```

### 3. Set Up Environment Variables

```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your configuration
# Update database credentials, secrets, and OAuth settings
```

### 4. Install golang-migrate

**macOS:**
```bash
brew install golang-migrate
```

**Linux:**
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

**Windows:**
```powershell
# Download from: https://github.com/golang-migrate/migrate/releases
# Add to PATH
```

### 5. Set Up PostgreSQL

```bash
# Create database
createdb laptop_tracking_dev

# Run migrations
make migrate-up
```

### 6. Set Up Mailhog (for email testing)

**macOS:**
```bash
brew install mailhog
mailhog
```

**Linux:**
```bash
go install github.com/mailhog/MailHog@latest
MailHog
```

**Docker:**
```bash
docker run -d -p 1025:1025 -p 8025:8025 mailhog/mailhog
```

Access Mailhog web UI at http://localhost:8025

### 7. Set Up Google OAuth

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable Google+ API
4. Create OAuth 2.0 credentials:
   - Application type: Web application
   - Authorized redirect URIs: `http://localhost:8080/auth/google/callback`
5. Copy Client ID and Client Secret to `.env` file

### 8. Load Sample Data (Optional)

To populate the database with sample data for testing and development:

**Windows (PowerShell):**
```powershell
.\scripts\load-sample-data.ps1
```

**macOS/Linux:**
```bash
make db-seed
```

**Or reset database and load sample data:**
```bash
# macOS/Linux
make db-reset-with-sample

# Windows (PowerShell)
# Run: .\scripts\load-sample-data.ps1 (after resetting database)
```

**Note:** These commands work with the Docker database container (`laptop-tracking-db`). Ensure Docker is running and the container is started.

**Sample users (all passwords: `password123`):**
- Logistics: `logistics@bairesdev.com`
- Client: `client1@techcorp.com`
- Warehouse: `warehouse@bairesdev.com`
- Project Manager: `pm@bairesdev.com`

**Sample data includes:**
- 9 users across all roles
- 5 client companies
- 10 software engineers
- 15 laptops (Dell, HP, Lenovo, Apple, Microsoft)
- 8 shipments in various statuses (pending pickup, in transit, at warehouse, delivered)
- 5 pickup forms with detailed information
- 3 reception reports
- 2 delivery forms

### 9. Run the Application

```bash
# Run with make
make run

# Or run directly
go run cmd/web/main.go
```

The application will be available at http://localhost:8080

## Development

### Available Make Commands

```bash
make help              # Display all available commands
make install           # Install Go dependencies
make build             # Build the application binary
make run               # Run the application
make test              # Run tests
make test-coverage     # Run tests with coverage report
make migrate-up        # Run all database migrations
make migrate-down      # Rollback last migration
make migrate-create    # Create new migration (usage: make migrate-create name=create_table)
make db-reset          # Reset database (drop and recreate)
make db-seed           # Load sample data into database
make db-reset-with-sample  # Reset database and load sample data
make dev-setup         # Set up development environment
make clean             # Clean build artifacts
make fmt               # Format code
make vet               # Run go vet
make lint              # Run linters
make check             # Run all checks (format, vet, lint, test)
```

### Database Migrations

Create a new migration:
```bash
make migrate-create name=create_users_table
```

This creates two files in `migrations/`:
- `000001_create_users_table.up.sql`
- `000001_create_users_table.down.sql`

Apply migrations:
```bash
make migrate-up
```

Rollback last migration:
```bash
make migrate-down
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test -v ./internal/models

# Run with race detection
go test -race ./...

# Run tests with test database (required for integration tests)
# Note: Set TEST_DATABASE_URL environment variable
$env:TEST_DATABASE_URL = "postgres://postgres:password@localhost:5432/laptop_tracking_dev?sslmode=disable"
go test -v ./internal/handlers

# Skip integration tests (run only unit tests)
go test -short ./...
```

### Test-Driven Development (TDD)

This project follows strict TDD methodology. See `docs/tdd.md` for the complete workflow.

**TDD Cycle:**
1. 🟥 **RED**: Write a failing test first
2. 🟩 **GREEN**: Implement minimal code to pass the test
3. 🛠 **REFACTOR**: Improve code structure (after tests pass)

**Example:**
```bash
# 1. Write failing test in *_test.go file
# 2. Run test to verify it fails
go test -v ./internal/handlers -run TestFeature

# 3. Implement feature
# 4. Run test again to verify it passes
go test -v ./internal/handlers -run TestFeature

# 5. Run full test suite to check for regressions
go test ./...

# 6. Commit only after tests pass
git commit -m "feat: implement feature to pass test"
```

### Code Quality

```bash
# Format code
make fmt

# Run go vet
make vet

# Run linters (requires golangci-lint)
make lint

# Run all checks
make check
```

## Docker Deployment

### Build Docker Image

```bash
make docker-build
```

### Run with Docker Compose

```bash
docker-compose up -d
```

## User Roles

1. **Logistics**: Manage shipments, view all data, coordinate pickups and deliveries
2. **Client**: Submit pickup forms for their company
3. **Warehouse**: Receive shipments, create reception reports, release hardware
4. **Project Manager**: View dashboards, reports, and shipment status

## Process Flow

The system enforces a sequential status flow to maintain data integrity:

1. **Pending Pickup from Client** → Client submits pickup form
2. **Pickup from Client Scheduled** → Logistics schedules pickup with courier and tracking number
3. **Picked Up from Client** → Courier confirms pickup
4. **In Transit to Warehouse** → Shipment en route
5. **At Warehouse** → Warehouse receives and creates reception report
6. **Released from Warehouse** → Logistics releases devices for delivery
7. **In Transit to Engineer** → Shipment en route to engineer (with ETA)
8. **Delivered** → Engineer confirms receipt via delivery form

**Note:** Status updates must follow this sequence. Users cannot skip stages or move backwards.

## Environment Variables

See `.env.example` for all available configuration options.

Key variables:
- `DB_*`: Database connection settings
- `GOOGLE_CLIENT_ID/SECRET`: Google OAuth credentials
- `SESSION_SECRET`: Session encryption key (must be kept secret)
- `SMTP_*`: Email server configuration
- `JIRA_*`: JIRA integration settings (optional)

## Security

- All passwords are hashed using bcrypt
- Sessions are encrypted and stored securely
- Magic links expire after single use or timeout
- CSRF protection on all forms
- Role-based access control (RBAC):
  - Dashboard: Logistics users only
  - Pickup forms: Client and Logistics users
  - Reception reports: Warehouse and Logistics users
  - Delivery forms: Software engineers (via magic links)
- Google OAuth restricted to @bairesdev.com domain
- SQL injection prevention via parameterized queries
- File upload validation and sanitization
- **Sequential status validation**: Shipments can only move forward through predefined stages, preventing status skipping or backwards transitions

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Commit Message Convention

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting)
- `refactor:` Code refactoring
- `test:` Test additions or changes
- `chore:` Build process or auxiliary tool changes

## Troubleshooting

### Database Connection Issues

1. Check PostgreSQL is running: `pg_isready`
2. Verify credentials in `.env`
3. Check database exists: `psql -l`
4. Test connection: `psql -h localhost -U postgres -d laptop_tracking_dev`

### Migration Errors

1. Check current migration version: `migrate -path migrations -database "postgresql://..." version`
2. Force specific version if needed: `make migrate-force version=N`
3. Reset database: `make db-reset`

### OAuth Issues

1. Verify redirect URLs match in Google Console and `.env`
2. Check OAuth credentials are correct
3. Ensure Google+ API is enabled
4. Clear browser cookies and try again

### Email Not Sending

1. Check Mailhog is running: `curl http://localhost:8025`
2. Verify SMTP settings in `.env`
3. Check email logs in application output

## License

This project is proprietary and confidential.

## Support

For issues and questions, please contact the development team.

---

**Built with ❤️ for BairesDev**
