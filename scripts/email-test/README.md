# Email Notification Test Script

This standalone Go program tests all email notifications in Align by sending test emails and verifying they arrive correctly in Mailhog.

## Quick Start

### Using the Wrapper Scripts (Recommended)

From the project root:

**Linux/macOS:**
```bash
./scripts/test_emails.sh
```

**Windows:**
```powershell
.\scripts\test_emails.ps1
```

These wrapper scripts will:
- Check that Mailhog is running
- Verify database connectivity
- Build and run the test program
- Display results
- Clean up test data

### Manual Execution

```bash
# From project root
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/laptop_tracking_test?sslmode=disable"
export MAILHOG_URL="http://localhost:8025"
export SMTP_HOST="localhost"
export SMTP_PORT="1025"

# Run directly
go run scripts/email-test/main.go

# Or build first
go build -o test_email scripts/email-test
./test_email
```

## What It Tests

This script verifies all 6 email notification types:

1. **Magic Link** - Secure form access
2. **Pickup Confirmation** - After pickup form submission  
3. **Pickup Scheduled** - When pickup date is set
4. **Warehouse Pre-Alert** - Alert warehouse of incoming shipment
5. **Release Notification** - Hardware ready for courier pickup
6. **Delivery Confirmation** - Device delivered to engineer

## Test Flow

For each notification, the script:

1. ✅ Creates necessary test data in database
2. 📧 Triggers the notification
3. ⏳ Waits for email to arrive
4. 🔍 Queries Mailhog API for the email
5. ✔️ Validates subject, recipient, and content
6. 📊 Records pass/fail result
7. 🧹 Cleans up test data

## Output

The script provides colored console output:

```
==================================================
  Email Notifications Test - Mailhog Verification
==================================================

Configuration:
  Database: postgres://...
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
     Email ID:  abc123def456

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
```

## Requirements

- Go 1.19+ 
- PostgreSQL database (test instance)
- Mailhog running on localhost:8025
- Required Go dependencies (automatically downloaded)

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/laptop_tracking_test?sslmode=disable` | Database connection string |
| `MAILHOG_URL` | `http://localhost:8025` | Mailhog API endpoint |
| `SMTP_HOST` | `localhost` | SMTP server hostname |
| `SMTP_PORT` | `1025` | SMTP server port |

## Exit Codes

- `0` - All tests passed
- `1` - One or more tests failed or setup error

## Troubleshooting

### Mailhog Not Running

**Error:** `❌ Failed to connect to Mailhog`

**Solution:** Start Mailhog:
```bash
mailhog  # or MailHog on Linux, MailHog.exe on Windows
```

### Database Connection Failed  

**Error:** `❌ Failed to connect to database`

**Solution:** 
- Verify PostgreSQL is running
- Check DATABASE_URL is correct
- Ensure database exists: `laptop_tracking_test`

### Emails Not Arriving

**Error:** `Email not found in Mailhog: no messages found`

**Possible Causes:**
- SMTP connection issue
- Mailhog crashed
- Timing/delay issue

**Solutions:**
- Verify Mailhog is running: `curl http://localhost:8025/api/v2/messages`
- Check SMTP settings match Mailhog (port 1025)
- Check application logs for SMTP errors

## Implementation Details

### Test Data Created

The script creates:
- 1 Client Company (`Test Client Company`)
- 1 Software Engineer (`engineer@test.com`)
- 1 Warehouse User (`warehouse@test.com`)
- 1 Logistics User (`logistics@test.com`)
- 1 Laptop (`TEST-SERIAL-123`)
- 1 Shipment (Full Journey type)
- 1 Pickup Form (with contact info)

All test data is automatically deleted after tests complete.

### Email Verification

For each email, the script verifies:
- ✅ Email was sent (no SMTP error)
- ✅ Email arrived in Mailhog
- ✅ Correct recipient
- ✅ Expected subject line
- ✅ Valid email structure

### Mailhog API

The script uses Mailhog's HTTP API:
- `GET /api/v2/messages` - Retrieve messages
- `DELETE /api/v1/messages` - Clear all messages

## Development

To add a new email notification test:

1. Implement the notification in `internal/email/notifier.go`
2. Add a test function in this file:

```go
func testNewNotification(ctx context.Context, notifier *email.Notifier, 
    mailhogURL string, shipmentID int64) TestResult {
    
    result := TestResult{NotificationType: "New Notification"}
    
    // Send notification
    err := notifier.SendNewNotification(ctx, shipmentID)
    if err != nil {
        result.Error = err.Error()
        return result
    }
    
    // Wait for email
    time.Sleep(200 * time.Millisecond)
    
    // Fetch from Mailhog
    msg, err := getLatestMailhogMessage(mailhogURL)
    if err != nil {
        result.Error = fmt.Sprintf("Email not found: %v", err)
        return result
    }
    
    // Extract details
    result.EmailID = msg.ID
    result.TimeSent = msg.Created
    result.Recipient = fmt.Sprintf("%s@%s", 
        msg.To[0].Mailbox, msg.To[0].Domain)
    
    if headers := msg.Content.Headers["Subject"]; len(headers) > 0 {
        result.Subject = headers[0]
    }
    
    // Validate
    if !strings.Contains(result.Subject, "Expected Text") {
        result.Error = fmt.Sprintf("Unexpected subject: %s", result.Subject)
        return result
    }
    
    result.Success = true
    return result
}
```

3. Call it from `main()`:

```go
result = testNewNotification(ctx, notifier, mailhogURL, testData.ShipmentID)
results = append(results, result)
printTestResult(result)
```

## Related Documentation

- Main README: `../../scripts/README.md`
- Email Templates: `../../internal/email/templates.go`
- Email Notifier: `../../internal/email/notifier.go`
- Phase 5 Documentation: `../../docs/PHASE_5_COMPLETE.md`

