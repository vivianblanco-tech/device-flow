# Email Notification Testing Scripts

This directory contains scripts to test email notifications for the Laptop Tracking System.

## Overview

The email testing script verifies that all email notifications are correctly sent and received by:

1. Setting up test data in the database
2. Triggering each email notification type
3. Verifying emails arrive in Mailhog
4. Checking email content (subject, recipient, etc.)
5. Cleaning up test data
6. Providing a detailed test report

## Prerequisites

### 1. Mailhog

Mailhog is a local email testing tool that captures SMTP emails.

**Installation:**

- **macOS:**
  ```bash
  brew install mailhog
  ```

- **Linux:**
  ```bash
  go install github.com/mailhog/MailHog@latest
  ```

- **Windows:**
  Download from [Mailhog Releases](https://github.com/mailhog/MailHog/releases)

**Running Mailhog:**
```bash
mailhog  # or MailHog on Linux, MailHog.exe on Windows
```

Then access the web UI at: http://localhost:8025

### 2. Database

Ensure you have a test database running. The script uses `laptop_tracking_test` by default.

## Usage

### Quick Start

**Linux/macOS:**
```bash
chmod +x scripts/test_emails.sh
./scripts/test_emails.sh
```

**Windows (PowerShell):**
```powershell
.\scripts\test_emails.ps1
```

### Manual Execution

If you prefer to run the Go script directly:

```bash
# Set environment variables
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/laptop_tracking_test?sslmode=disable"
export MAILHOG_URL="http://localhost:8025"
export SMTP_HOST="localhost"
export SMTP_PORT="1025"

# Build and run
go run scripts/email-test/main.go

# Or build first, then run
go build -o test_email_notifications scripts/email-test
./test_email_notifications
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/laptop_tracking_test?sslmode=disable` | PostgreSQL connection string |
| `MAILHOG_URL` | `http://localhost:8025` | Mailhog API endpoint |
| `SMTP_HOST` | `localhost` | SMTP server hostname |
| `SMTP_PORT` | `1025` | SMTP server port |

## Email Notification Tests

The script tests all 6 email notification types:

### 1. Magic Link Email
- **Trigger**: When a user needs secure form access
- **Recipient**: test@example.com (test user)
- **Verifies**: Subject contains "Access Your Form", correct recipient

### 2. Pickup Confirmation
- **Trigger**: After pickup form submission
- **Recipient**: Client company contact
- **Verifies**: Subject contains "Pickup Confirmation", email sent

### 3. Pickup Scheduled Notification
- **Trigger**: When pickup is officially scheduled
- **Recipient**: Contact email from pickup form (contact@test.com)
- **Verifies**: Subject contains "Pickup Scheduled", correct recipient

### 4. Warehouse Pre-Alert
- **Trigger**: Alert warehouse about incoming shipment
- **Recipient**: Warehouse user (warehouse@test.com)
- **Verifies**: Subject contains "Incoming Shipment Alert", correct recipient

### 5. Release Notification
- **Trigger**: Hardware released from warehouse
- **Recipient**: Logistics user (logistics@test.com)
- **Verifies**: Subject contains "Hardware Release for Pickup", correct recipient

### 6. Delivery Confirmation
- **Trigger**: Device delivered to engineer
- **Recipient**: Software engineer (engineer@test.com)
- **Verifies**: Subject contains "Device Delivered Successfully", correct recipient

## Output

The script provides colored output with:

- ✅ **Green**: Successful tests
- ❌ **Red**: Failed tests
- ⚠️ **Yellow**: Warnings
- 💡 **Blue**: Information

### Example Output

```
==================================================
  Email Notifications Test - Mailhog Verification
==================================================

Configuration:
  Database: postgres://postgres:postgres@localhost:5432/laptop_tracking_test?sslmode=disable
  Mailhog:  http://localhost:8025
  SMTP:     localhost:1025

✅ Connected to database
✅ Connected to Mailhog
✅ Cleared Mailhog messages
✅ Test data created

===========================================
  Running Email Notification Tests
===========================================

📧 Test 1: Magic Link Email
  ✅ SUCCESS
     Subject:   Access Your Form - pickup
     Recipient: test@example.com
     Email ID:  abc123

📧 Test 2: Pickup Confirmation
  ✅ SUCCESS
     Subject:   Pickup Confirmation - CONF-1
     Recipient: noreply@example.com
     Email ID:  def456

... (more tests)

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

💡 You can view the emails in Mailhog:
   http://localhost:8025
```

## Test Data

The script creates temporary test data:
- 1 Client Company
- 1 Software Engineer
- 1 Warehouse User
- 1 Logistics User
- 1 Laptop
- 1 Shipment (with pickup form)

All test data is automatically cleaned up after the tests complete.

## Troubleshooting

### Mailhog Not Running

```
❌ Failed to connect to Mailhog: Get "http://localhost:8025/api/v2/messages": dial tcp 127.0.0.1:8025: connect: connection refused
```

**Solution**: Start Mailhog:
```bash
mailhog
```

### Database Connection Failed

```
❌ Failed to connect to database: dial tcp 127.0.0.1:5432: connect: connection refused
```

**Solution**: Start your PostgreSQL database or update `DATABASE_URL`

### Email Not Found in Mailhog

```
❌ FAILED
   Error: Email not found in Mailhog: no messages found
```

**Possible causes**:
1. Email wasn't sent (check SMTP configuration)
2. Timing issue (script may need longer delay)
3. Mailhog crashed or restarted

**Solution**:
- Check Mailhog is running
- Verify SMTP settings
- Check application logs for errors

### Wrong Database

```
❌ Failed to set up test data: create client company: pq: permission denied for table client_companies
```

**Solution**: Ensure you're using a test database with proper permissions

## Viewing Emails in Mailhog

After running the tests, open http://localhost:8025 to:

1. View all captured emails
2. Check HTML rendering
3. Inspect email headers
4. Download email source
5. Delete emails

## CI/CD Integration

You can integrate this script into your CI/CD pipeline:

```yaml
# Example for GitHub Actions
- name: Setup Mailhog
  run: |
    go install github.com/mailhog/MailHog@latest
    MailHog &
    sleep 2

- name: Test Email Notifications
  run: ./scripts/test_emails.sh
  env:
    DATABASE_URL: postgres://postgres:postgres@localhost:5432/laptop_tracking_test?sslmode=disable
```

## Development

To add a new email notification test:

1. Add the notification method to `internal/email/notifier.go`
2. Add a new test function in `test_email_notifications.go`:
   ```go
   func testNewNotification(ctx context.Context, notifier *email.Notifier, mailhogURL string, shipmentID int64) TestResult {
       result := TestResult{NotificationType: "New Notification"}
       
       err := notifier.SendNewNotification(ctx, shipmentID)
       if err != nil {
           result.Error = err.Error()
           return result
       }
       
       time.Sleep(200 * time.Millisecond)
       
       msg, err := getLatestMailhogMessage(mailhogURL)
       if err != nil {
           result.Error = fmt.Sprintf("Email not found in Mailhog: %v", err)
           return result
       }
       
       // Validate email content...
       
       result.Success = true
       return result
   }
   ```
3. Call the test function in `main()`:
   ```go
   result = testNewNotification(ctx, notifier, mailhogURL, testData.ShipmentID)
   results = append(results, result)
   printTestResult(result)
   ```

## License

Part of the Laptop Tracking System project.
