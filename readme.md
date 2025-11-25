# Align

A comprehensive web application for tracking laptop pickup and delivery from client companies to software engineers, including inventory management and status tracking.

**Project Status**: 🟢 **85% Complete** - Core functionality implemented and operational

## Features

### Core Functionality
- 📦 **Shipment Tracking**: Track laptops through complete lifecycle with three shipment types:
  - **Single Full Journey**: One laptop from client → warehouse → engineer
  - **Bulk to Warehouse**: Multiple laptops from client to warehouse
  - **Warehouse to Engineer**: Direct shipment from warehouse inventory to engineer
- 👥 **Multi-Role System**: Support for Logistics, Client, Warehouse, and Project Manager roles with role-based access control
- 🔐 **Dual Authentication**: Username/password and Google OAuth (restricted to @bairesdev.com)
- 🔗 **Magic Links**: Secure one-time access links for form submissions
- 📧 **Email Notifications**: Automated notifications at each process step (6 email templates)
- 🎫 **JIRA Integration**: Sync with JIRA tickets for seamless workflow (create, update, sync status)

### Dashboard & Visualization
- 📊 **Dashboard**: Real-time statistics, KPIs, and system overview (Logistics users only)
- 📈 **Interactive Charts**: Line, Donut, and Bar charts powered by Chart.js v4.4.1
- 📅 **Calendar View**: Visual timeline of pickups and deliveries with date filtering
- 🔍 **Inventory Management**: Full CRUD operations with search, filter, and status tracking

### Forms & Management
- 📋 **Pickup Forms**: Three types of shipment creation forms with validation
- 📥 **Reception Reports**: Warehouse intake with photo uploads and approval workflow
- 📤 **Delivery Forms**: Engineer confirmation with photo documentation
- 👤 **User Management**: Create and edit users (Logistics only)
- 🏢 **Client Company Management**: Manage client organizations
- 💻 **Software Engineer Management**: Track engineer profiles with address confirmation
- 🚚 **Courier Management**: Manage courier services
- 📊 **Reports**: Shipment status, inventory summary, and timeline reports (Client users)

### Additional Features
- ✏️ **Shipment Editing**: Edit shipment details, add laptops to bulk shipments, assign engineers
- 📸 **Photo Uploads**: Document device condition at pickup and delivery
- 🔄 **Status Management**: Sequential status flow with validation
- 📝 **Audit Logging**: Complete audit trail of all system actions
- 🔔 **Notification Logging**: Track all email notifications sent

## Tech Stack

- **Backend**: Go 1.24+
- **Database**: PostgreSQL 15+
- **Frontend**: HTML templates with Tailwind CSS v4
- **Charts**: Chart.js v4.4.1 (Line, Donut, Bar charts)
- **Authentication**: Google OAuth 2.0, bcrypt for passwords
- **Email**: SMTP (Mailhog for development)
- **Migrations**: golang-migrate
- **Deployment**: Docker, Caddy (reverse proxy with automatic SSL)
- **Testing**: Comprehensive test suite with 258+ test cases (TDD methodology)

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

- Go 1.24 or higher
- PostgreSQL 15 or higher
- Docker (for database container and testing)
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

**Option A: Using Docker (Recommended)**
```bash
# Start PostgreSQL container
docker-compose up -d postgres

# Run migrations
make migrate-up
```

**Option B: Local PostgreSQL**
```bash
# Create database
createdb laptop_tracking_dev

# Run migrations
make migrate-up
```

**Set Up Test Database (for running tests)**
```bash
# Using Docker
make test-db-setup

# Or manually
createdb laptop_tracking_test
migrate -path migrations -database "postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable" up
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

## Available Routes

### Public Routes
- `GET /` - Redirects to dashboard (if authenticated) or login
- `GET /health` - Health check endpoint
- `GET /login` - Login page
- `POST /login` - Login form submission
- `GET /logout` - Logout
- `GET /auth/google` - Google OAuth login
- `GET /auth/google/callback` - Google OAuth callback
- `GET /auth/magic-link` - Magic link authentication

### Protected Routes (Require Authentication)

**Dashboard & Analytics**
- `GET /dashboard` - Main dashboard with statistics (Logistics only)
- `GET /api/charts/shipments-over-time` - Shipments over time chart data
- `GET /api/charts/status-distribution` - Status distribution chart data
- `GET /api/charts/delivery-time-trends` - Delivery time trends chart data

**Calendar**
- `GET /calendar` - Calendar view of pickups and deliveries

**Shipments**
- `GET /shipments` - List all shipments (with filters)
- `GET /shipments/create` - Create new shipment
- `GET /shipments/{id}` - View shipment details
- `POST /shipments/{id}/status` - Update shipment status
- `POST /shipments/{id}/assign-engineer` - Assign engineer to shipment
- `GET /shipments/{id}/edit` - Edit shipment page
- `POST /shipments/{id}/edit` - Update shipment
- `GET /shipments/{id}/form` - Pickup form for shipment
- `POST /shipments/{id}/form` - Submit pickup form
- `POST /shipments/{id}/complete-details` - Complete shipment details
- `POST /shipments/{id}/edit-details` - Edit shipment details
- `POST /shipments/{id}/laptops/add` - Add laptop to bulk shipment

**Shipment Creation Forms**
- `GET /shipments/create/single` - Single full journey form
- `POST /shipments/create/single-minimal` - Create minimal single shipment
- `GET /shipments/create/bulk` - Bulk to warehouse form
- `GET /shipments/create/warehouse-to-engineer` - Warehouse to engineer form
- `GET /pickup-forms` - Pickup forms landing page
- `GET /pickup-form` - Legacy pickup form
- `POST /pickup-form` - Submit pickup form

**Inventory**
- `GET /inventory` - List all laptops
- `GET /inventory/add` - Add laptop page
- `POST /inventory/add` - Create new laptop
- `GET /inventory/{id}` - View laptop details
- `GET /inventory/{id}/edit` - Edit laptop page
- `POST /inventory/{id}/update` - Update laptop
- `POST /inventory/{id}/delete` - Delete laptop

**Reception Reports**
- `GET /reception-reports` - List all reception reports
- `GET /reception-reports/{id}` - View reception report details
- `POST /reception-reports/{id}/approve` - Approve reception report
- `GET /laptops/{id}/reception-report` - Create reception report for laptop
- `POST /laptops/{id}/reception-report` - Submit reception report

**Delivery Forms**
- `GET /delivery-form` - Delivery form page
- `POST /delivery-form` - Submit delivery form

**Forms Management (Logistics Only)**
- `GET /forms` - Forms management dashboard
- `GET /forms/users` - List users
- `GET /forms/users/add` - Add user page
- `POST /forms/users/add` - Create user
- `GET /forms/users/{id}/edit` - Edit user page
- `POST /forms/users/{id}/edit` - Update user
- `GET /forms/client-companies` - List client companies
- `GET /forms/client-companies/add` - Add client company page
- `POST /forms/client-companies/add` - Create client company
- `GET /forms/client-companies/{id}/edit` - Edit client company page
- `POST /forms/client-companies/{id}/edit` - Update client company
- `GET /forms/software-engineers` - List software engineers
- `GET /forms/software-engineers/add` - Add engineer page
- `POST /forms/software-engineers/add` - Create engineer
- `GET /forms/software-engineers/{id}/edit` - Edit engineer page
- `POST /forms/software-engineers/{id}/edit` - Update engineer
- `GET /forms/couriers` - List couriers
- `GET /forms/couriers/add` - Add courier page
- `POST /forms/couriers/add` - Create courier
- `GET /forms/couriers/{id}/edit` - Edit courier page
- `POST /forms/couriers/{id}/edit` - Update courier

**Magic Links (Logistics Only)**
- `GET /magic-links` - List magic links
- `POST /auth/send-magic-link` - Send magic link

**Reports (Client Users)**
- `GET /reports` - Reports index
- `GET /reports/shipment-status` - Shipment status dashboard
- `GET /reports/inventory-summary` - Inventory summary report
- `GET /reports/shipment-timeline` - Shipment timeline report

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

**Test Database Setup (Required for Integration Tests)**
```bash
# Set up test database using Docker
make test-db-setup

# Or manually create test database
createdb laptop_tracking_test
migrate -path migrations -database "postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable" up
```

**Running Tests**
```bash
# Run all tests (sequential, recommended)
make test-all

# Run all tests in parallel (faster but may have conflicts)
make test-parallel

# Run only unit tests (no database required)
make test-unit

# Run tests with coverage report
make test-coverage

# Run specific package tests
make test-package PKG=internal/models

# Run with race detection
go test -race ./...

# Quick test run (unit tests only, no race detection)
make test-quick

# CI mode (with coverage and sequential execution)
make test-ci
```

**Test Statistics:**
- **Total Tests**: 258+ test cases
- **Unit Tests**: 181 tests (passing without database)
- **Integration Tests**: 77 tests (require test database)
- **Test Coverage**: 
  - Models: 97.7%
  - Validators: 95.9%
  - Config: 100%
  - JIRA: 61.8%

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

## User Roles & Permissions

### Logistics
- ✅ Full system access
- ✅ Dashboard with statistics and charts
- ✅ Create and manage all shipment types
- ✅ Edit shipments and assign engineers
- ✅ Manage users, client companies, software engineers, and couriers
- ✅ Create and manage magic links
- ✅ View all shipments and inventory
- ✅ Approve reception reports

### Client
- ✅ View shipments for their company
- ✅ Submit pickup forms for their company
- ✅ View their company's laptop inventory
- ✅ Access reports (shipment status, inventory summary, timeline)
- ✅ View calendar
- ❌ Cannot access dashboard
- ❌ Cannot create reception reports

### Warehouse
- ✅ View shipments
- ✅ Create reception reports for laptops
- ✅ View inventory
- ✅ View calendar
- ❌ Cannot access dashboard
- ❌ Cannot create shipments
- ❌ Cannot approve reception reports

### Project Manager
- ✅ View dashboard with statistics and charts
- ✅ View all shipments
- ✅ View inventory
- ✅ View calendar
- ❌ Cannot create or edit shipments
- ❌ Cannot access forms management

## Process Flow

The system supports three shipment types, each with its own workflow:

### Shipment Types

1. **Single Full Journey** (`single_full_journey`)
   - Complete lifecycle: Client → Warehouse → Engineer
   - All 8 status stages available
   - Best for: Individual laptop shipments

2. **Bulk to Warehouse** (`bulk_to_warehouse`)
   - Client → Warehouse only
   - Statuses: Pending Pickup → Scheduled → Picked Up → In Transit → At Warehouse
   - Best for: Multiple laptops sent to warehouse for later distribution

3. **Warehouse to Engineer** (`warehouse_to_engineer`)
   - Warehouse → Engineer only
   - Statuses: Released from Warehouse → In Transit → Delivered
   - Best for: Direct shipments from warehouse inventory

### Status Flow

The system enforces a sequential status flow to maintain data integrity:

1. **Pending Pickup from Client** → Client submits pickup form
2. **Pickup from Client Scheduled** → Logistics schedules pickup with courier and tracking number
3. **Picked Up from Client** → Courier confirms pickup
4. **In Transit to Warehouse** → Shipment en route
5. **At Warehouse** → Warehouse receives and creates reception report
6. **Released from Warehouse** → Logistics releases devices for delivery
7. **In Transit to Engineer** → Shipment en route to engineer (with ETA)
8. **Delivered** → Engineer confirms receipt via delivery form

**Note:** 
- Status updates must follow this sequence. Users cannot skip stages or move backwards.
- Available statuses depend on shipment type (bulk shipments don't go to engineer, warehouse-to-engineer starts at "Released").
- Status transitions are validated per shipment type.

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
- Role-based access control (RBAC) with granular permissions per role
- Google OAuth restricted to @bairesdev.com domain
- SQL injection prevention via parameterized queries
- File upload validation and sanitization
- **Sequential status validation**: Shipments can only move forward through predefined stages, preventing status skipping or backwards transitions
- **Shipment type validation**: Status transitions validated per shipment type
- Session-based authentication with secure token storage
- Magic links expire after single use or timeout
- CSRF protection on all forms
- Input validation and sanitization on all user inputs

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

## Project Status

### Completed Phases (85% Complete)

- ✅ **Phase 0**: Project Setup & Infrastructure (100%)
- ✅ **Phase 1**: Database Schema & Core Models (100%) - 133 test cases, 97.7% coverage
- ✅ **Phase 2**: Authentication System (100%) - OAuth, RBAC, Magic Links
- ✅ **Phase 3**: Core Forms & Workflows (100%) - Pickup, Reception, Delivery forms
- ✅ **Phase 4**: JIRA Integration (100%) - Full sync capabilities
- ✅ **Phase 5**: Email Notifications (100%) - 6 email templates
- ✅ **Phase 6**: Dashboard & Visualization (100%) - Charts, Calendar, Inventory

### In Progress / Pending

- 🟡 **Phase 7**: Comprehensive Testing (40%) - Integration tests need test database setup
- 🟡 **Phase 8**: Deployment & DevOps (30%) - Docker ready, needs production config
- 🟡 **Phase 9**: Polish & Documentation (20%) - UI/UX improvements, security hardening

### Key Metrics

- **Total Test Cases**: 258+
- **Code Coverage**: 97.7% (models), 95.9% (validators)
- **Database Migrations**: 52 files (26 up, 26 down)
- **Routes**: 50+ endpoints
- **Templates**: 40+ HTML pages
- **Test-Driven Development**: Strict TDD methodology followed throughout

---

**Built with ❤️ for BairesDev**
