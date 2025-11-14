# Email Testing Quick Start Guide

This guide will help you quickly test all email notifications in the Laptop Tracking System.

## Prerequisites (5 minutes)

### 1. Install Mailhog

**macOS:**
```bash
brew install mailhog
```

**Linux:**
```bash
go install github.com/mailhog/MailHog@latest
# Add to PATH: export PATH=$PATH:$(go env GOPATH)/bin
```

**Windows:**
1. Download from: https://github.com/mailhog/MailHog/releases
2. Extract `MailHog.exe` to a folder
3. Add to PATH or run from that folder

### 2. Start Mailhog

```bash
mailhog
```

You should see:
```
[HTTP] Binding to address: 0.0.0.0:8025
[SMTP] Binding to address: 0.0.0.0:1025
```

✅ **Verify it's working:** Open http://localhost:8025 in your browser

### 3. Database Setup

Make sure your test database is running:

```bash
# If not already set up
psql -U postgres -c "CREATE DATABASE laptop_tracking_test;"

# Run migrations
go run cmd/migrate/main.go
```

## Running the Tests (1 minute)

### Option 1: Using Shell Scripts (Easiest)

**Linux/macOS:**
```bash
chmod +x scripts/test_emails.sh
./scripts/test_emails.sh
```

**Windows (PowerShell):**
```powershell
.\scripts\test_emails.ps1
```

### Option 2: Manual Execution

```bash
# Set environment (optional - defaults work with standard setup)
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/laptop_tracking_test?sslmode=disable"
export MAILHOG_URL="http://localhost:8025"

# Run the test
go run scripts/email-test/main.go
```

## What to Expect

The script will:

1. ✅ Connect to database and Mailhog
2. 🧹 Clear any existing test emails
3. 📦 Create test data (shipments, users, etc.)
4. 📧 Send 6 different email notifications
5. ✔️ Verify each email arrived correctly
6. 📊 Display detailed results
7. 🧹 Clean up test data

**Typical output:**
```
==================================================
  Email Notifications Test - Mailhog Verification
==================================================

✅ Connected to database
✅ Connected to Mailhog
✅ Test data created

Running Email Notification Tests
===========================================

📧 Test 1: Magic Link Email
  ✅ SUCCESS
     Subject:   Access Your Form - pickup
     Recipient: test@example.com

📧 Test 2: Pickup Confirmation
  ✅ SUCCESS
     Subject:   Pickup Confirmation - CONF-1
     Recipient: noreply@example.com

... (4 more tests)

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

🎉 All email notifications passed!
```

## View the Emails

After the tests complete, you can view all the emails in Mailhog:

**Open:** http://localhost:8025

You'll see all 6 test emails with:
- Full HTML rendering
- Email headers
- Plain text version
- Ability to download/inspect

## Notifications Tested

| # | Notification Type | Recipient | Trigger |
|---|-------------------|-----------|---------|
| 1 | **Magic Link** | Form user | Secure form access needed |
| 2 | **Pickup Confirmation** | Client | After pickup form submission |
| 3 | **Pickup Scheduled** | Contact | When pickup date is set |
| 4 | **Warehouse Pre-Alert** | Warehouse team | Incoming shipment alert |
| 5 | **Release Notification** | Logistics | Hardware ready for pickup |
| 6 | **Delivery Confirmation** | Engineer | Device delivered |

## Troubleshooting

### ❌ "Mailhog is not running"

**Fix:** 
```bash
mailhog
```

Leave it running in a separate terminal.

### ❌ "Cannot connect to database"

**Fix:**
1. Check PostgreSQL is running: `psql -U postgres -l`
2. Verify database exists: `psql -U postgres -c "\l laptop_tracking_test"`
3. Check DATABASE_URL environment variable

### ❌ "Email not found in Mailhog"

**Possible causes:**
- Mailhog crashed (restart it)
- SMTP connection issue (check Mailhog logs)
- Timing issue (emails arrive slowly)

**Fix:**
1. Restart Mailhog
2. Check Mailhog UI shows it's receiving on port 1025
3. Run tests again

### ❌ Tests fail but emails look fine in Mailhog

The tests validate specific subjects and recipients. If the emails are arriving but tests fail:
1. Check the error message for what doesn't match
2. Verify email templates haven't changed
3. Check test data setup

## Configuration

You can customize the test environment:

```bash
# Custom database
export DATABASE_URL="postgres://user:pass@host:5432/dbname?sslmode=disable"

# Custom Mailhog (if running on different port)
export MAILHOG_URL="http://localhost:9025"
export SMTP_PORT="1026"
```

## Next Steps

After verifying emails work:

1. **Update Templates**: Edit `internal/email/templates.go`
2. **Add Notifications**: Implement in `internal/email/notifier.go`
3. **Integration**: Call notifications from handlers
4. **Production Setup**: Configure real SMTP server

## Need Help?

- 📖 Full documentation: `scripts/README.md`
- 🔧 Email implementation: `internal/email/`
- 📋 Phase 5 docs: `docs/PHASE_5_COMPLETE.md`

## Cleaning Up

The script automatically cleans up all test data. If you need to manually clean:

```sql
DELETE FROM notification_logs WHERE recipient LIKE '%test.com';
DELETE FROM pickup_forms WHERE shipment_id IN (SELECT id FROM shipments WHERE tracking_number LIKE 'TEST-%');
DELETE FROM shipment_laptops WHERE shipment_id IN (SELECT id FROM shipments WHERE tracking_number LIKE 'TEST-%');
DELETE FROM shipments WHERE tracking_number LIKE 'TEST-%';
DELETE FROM laptops WHERE serial_number LIKE 'TEST-%';
DELETE FROM users WHERE email LIKE '%test.com';
DELETE FROM software_engineers WHERE email LIKE '%test.com';
DELETE FROM client_companies WHERE name LIKE 'Test %';
```

---

**Time to run:** ~30 seconds  
**Test coverage:** All 6 notification types  
**Database impact:** Temporary (auto-cleaned)

