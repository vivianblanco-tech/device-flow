# Email Notifications Implementation Plan

**Date Created:** November 14, 2025  
**Status:** Planning Complete - Ready for Implementation

---

## Executive Summary

This document outlines the comprehensive plan for implementing all required email notifications in Align. The analysis shows that 6 email notification templates exist but only 2 are properly wired up, and 5 additional notifications need to be created from scratch.

---

## Current State Analysis

### Existing Email Infrastructure

**Location:** `internal/email/`
- `client.go` - SMTP email client with TLS support
- `notifier.go` - Notification service with 6 implemented functions
- `templates.go` - HTML email templates with responsive design

**Database Tables:**
- `notification_logs` - Tracks all sent notifications
- `audit_logs` - Tracks system actions

**Configuration:**
- SMTP settings in `internal/config/config.go`
- Email client initialized in `cmd/web/main.go`

---

## Implemented Email Notifications

### 1. Magic Link Email ✅ WORKING
- **Template:** `magic_link`
- **Function:** `SendMagicLink(ctx, recipientEmail, recipientName, magicLink, formType, expiresAt)`
- **Trigger:** Manual - Logistics user sends via UI (`/auth/send-magic-link`)
- **Recipient:** Email address entered by logistics
- **Purpose:** Allow clients to access pickup form without login

**Status:** Fully functional but manual process. Email is shown in UI instead of being sent automatically.

### 2. Pickup Confirmation ✅ IMPLEMENTED BUT NEEDS FIX
- **Template:** `pickup_confirmation`
- **Function:** `SendPickupConfirmation(ctx, shipmentID)`
- **Trigger:** When pickup form is submitted (all shipment types)
- **Recipient:** Currently sends to client company user (WRONG)
- **Purpose:** Confirm pickup form submission

**Issues:**
- Sends to wrong email - should use `contact_email` from pickup form
- Currently queries users table instead of pickup form data
- Needs to be fixed to use form contact information

**Location:** `internal/handlers/pickup_form.go:521-527`

### 3. Pickup Scheduled Notification ✅ FULLY WORKING
- **Template:** `pickup_scheduled`
- **Function:** `SendPickupScheduledNotification(ctx, shipmentID)`
- **Trigger:** When status changes from `pending_pickup_from_client` → `pickup_from_client_scheduled`
- **Recipient:** Contact email from pickup form
- **Purpose:** Confirm pickup has been scheduled with date/time

**Status:** Fully functional, properly wired up.

**Location:** `internal/handlers/shipments.go:544-569`

### 4. Warehouse Pre-Alert ⚠️ EXISTS BUT NOT WIRED UP
- **Template:** `warehouse_pre_alert`
- **Function:** `SendWarehousePreAlert(ctx, shipmentID)`
- **Trigger:** NONE - needs to be added
- **Recipient:** Warehouse user email
- **Purpose:** Alert warehouse of incoming shipment

**Status:** Template and function exist, needs to be connected to status change handler.

### 5. Release Notification ⚠️ EXISTS BUT NOT WIRED UP
- **Template:** `release_notification`
- **Function:** `SendReleaseNotification(ctx, shipmentID)`
- **Trigger:** NONE - needs to be added
- **Recipient:** Logistics user (`international@bairesdev.com`)
- **Purpose:** Notify logistics that hardware is ready for pickup from warehouse

**Status:** Template and function exist, needs to be connected to status change handler.

### 6. Delivery Confirmation ⚠️ EXISTS BUT NOT WIRED UP
- **Template:** `delivery_confirmation`
- **Function:** `SendDeliveryConfirmation(ctx, shipmentID)`
- **Trigger:** NONE - needs to be added
- **Recipient:** Software engineer email
- **Purpose:** Confirm device delivered to engineer

**Status:** Template and function exist, needs to be connected to status change handler.

---

## Missing Email Notifications (Need to be Created)

### 7. Shipment Picked Up - To Client ❌ NOT IMPLEMENTED
- **Template:** NEEDS CREATION - `shipment_picked_up`
- **Function:** NEEDS CREATION - `SendShipmentPickedUpNotification()`
- **Trigger:** Status change to `picked_up_from_client`
- **Recipient:** Contact email from pickup form
- **Content:**
  - Tracking number
  - Courier information
  - Expected arrival at warehouse (estimate 3 days)
  - Link to track shipment
- **Applies to:** All shipment types

### 8. Delivered to Engineer - To Client ❌ NOT IMPLEMENTED
- **Template:** NEEDS CREATION - `engineer_delivery_notification_to_client`
- **Function:** NEEDS CREATION - `SendEngineerDeliveryNotificationToClient()`
- **Trigger:** Status change to `delivered`
- **Recipient:** Contact email from pickup form
- **Content:**
  - Confirmation laptop delivered to engineer
  - Engineer name
  - Delivery date
  - Project/ticket information
- **Applies to:** Single Full Journey and Warehouse to Engineer ONLY

### 9. Pickup Form Submitted - To Logistics ❌ NOT IMPLEMENTED
- **Template:** NEEDS CREATION - `pickup_form_submitted_logistics`
- **Function:** NEEDS CREATION - `SendPickupFormSubmittedNotification()`
- **Trigger:** When any pickup form is submitted
- **Recipient:** `international@bairesdev.com`
- **Content:**
  - Shipment type
  - Client company information
  - Pickup address and details
  - Contact information
  - Number of devices
  - JIRA ticket number
  - Link to shipment details
- **Applies to:** All shipment types

### 10. Reception Report Created - To Logistics ❌ NOT IMPLEMENTED
- **Template:** NEEDS CREATION - `reception_report_approval_request`
- **Function:** NEEDS CREATION - `SendReceptionReportApprovalRequest()`
- **Trigger:** When reception report is created by warehouse
- **Recipient:** `international@bairesdev.com`
- **Content:**
  - Reception report details
  - Photos/attachments
  - Warehouse notes
  - Serial numbers (if tracked)
  - Request for approval
  - Link to review/approve
- **Applies to:** All shipment types arriving at warehouse
- **Note:** Currently has TODO comment in `reception_report.go:303`

### 11. In Transit to Engineer - To Engineer ❌ NOT IMPLEMENTED
- **Template:** NEEDS CREATION - `in_transit_to_engineer`
- **Function:** NEEDS CREATION - `SendInTransitToEngineerNotification()`
- **Trigger:** Status change to `in_transit_to_engineer`
- **Recipient:** Software engineer email
- **Content:**
  - ETA (from shipment.eta_to_engineer field)
  - Tracking number
  - Courier information
  - What to expect
  - Contact information for issues
- **Applies to:** Single Full Journey and Warehouse to Engineer ONLY

---

## Shipment Types & Email Notification Matrix

The system supports 3 shipment types with different notification flows:

### Shipment Type 1: Single Full Journey (`single_full_journey`)
**Flow:** Client → Pickup → Warehouse → Engineer

| Status | Email Notifications |
|--------|-------------------|
| `pending_pickup_from_client` | Magic Link (manual) |
| Form Submitted | ✉️ Pickup Confirmation (Client)<br>✉️ Form Submitted (Logistics) |
| `pickup_from_client_scheduled` | ✉️ Pickup Scheduled (Client) |
| `picked_up_from_client` | ✉️ Shipment Picked Up (Client)<br>✉️ Warehouse Pre-Alert (Warehouse) |
| `in_transit_to_warehouse` | - |
| `at_warehouse` | ✉️ Reception Report Approval (Logistics) |
| `released_from_warehouse` | ✉️ Release Notification (Logistics) |
| `in_transit_to_engineer` | ✉️ In Transit (Engineer) |
| `delivered` | ✉️ Delivery Confirmation (Engineer)<br>✉️ Delivered (Client) |

### Shipment Type 2: Bulk to Warehouse (`bulk_to_warehouse`)
**Flow:** Client → Pickup → Warehouse (ENDS HERE)

| Status | Email Notifications |
|--------|-------------------|
| `pending_pickup_from_client` | Magic Link (manual) |
| Form Submitted | ✉️ Pickup Confirmation (Client)<br>✉️ Form Submitted (Logistics) |
| `pickup_from_client_scheduled` | ✉️ Pickup Scheduled (Client) |
| `picked_up_from_client` | ✉️ Shipment Picked Up (Client)<br>✉️ Warehouse Pre-Alert (Warehouse) |
| `in_transit_to_warehouse` | - |
| `at_warehouse` | ✉️ Reception Report Approval (Logistics) |

**Note:** No engineer-related notifications for this type.

### Shipment Type 3: Warehouse to Engineer (`warehouse_to_engineer`)
**Flow:** Warehouse → Engineer (STARTS FROM WAREHOUSE)

| Status | Email Notifications |
|--------|-------------------|
| `released_from_warehouse` | ✉️ Release Notification (Logistics) |
| `in_transit_to_engineer` | ✉️ In Transit (Engineer) |
| `delivered` | ✉️ Delivery Confirmation (Engineer)<br>✉️ Delivered (Client) |

**Note:** No pickup-related notifications for this type (no client pickup involved).

---

## Implementation Plan

### Phase 1: Fix & Wire Up Existing Notifications (3-4 hours)

#### Task 1.1: Fix Pickup Confirmation Email
**File:** `internal/email/notifier.go` - `SendPickupConfirmation()` method

**Current Code Issue (lines 36-53):**
```go
// Fetch client company
var clientName, clientCompany, clientEmail string
err = n.db.QueryRowContext(ctx,
    `SELECT name, contact_info FROM client_companies WHERE id = $1`,
    shipment.ClientCompanyID,
).Scan(&clientName, &clientCompany)

// Get client user email
err = n.db.QueryRowContext(ctx,
    `SELECT email FROM users WHERE id IN (
        SELECT id FROM users WHERE role = 'client' LIMIT 1
    )`,
).Scan(&clientEmail)
```

**Required Changes:**
1. Fetch contact email from pickup form instead of users table
2. Use pickup form data for all contact information
3. Ensure it works for all 3 shipment types

**Implementation Steps:**
- Query pickup form: `SELECT form_data FROM pickup_forms WHERE shipment_id = $1`
- Parse JSON to extract `contact_email`, `contact_name`
- Update template data to use form contact info
- Test with all shipment types

#### Task 1.2: Wire Up Warehouse Pre-Alert
**File:** `internal/handlers/shipments.go` - `UpdateShipmentStatus()` method

**Trigger Point:** When status changes to `picked_up_from_client`

**Implementation:**
```go
// Add after line 569 (after pickup scheduled notification)
if newStatus == models.ShipmentStatusPickedUpFromClient {
    if h.EmailNotifier != nil {
        go func() {
            ctx := context.Background()
            if err := h.EmailNotifier.SendWarehousePreAlert(ctx, shipmentID); err != nil {
                fmt.Printf("Warning: failed to send warehouse pre-alert: %v\n", err)
            }
        }()
    }
}
```

**Also trigger in:** `internal/handlers/pickup_form.go` if pickup happens immediately

#### Task 1.3: Wire Up Release Notification
**File:** `internal/handlers/shipments.go` - `UpdateShipmentStatus()` method

**Trigger Point:** When status changes to `released_from_warehouse`

**Implementation:**
```go
if newStatus == models.ShipmentStatusReleasedFromWarehouse {
    if h.EmailNotifier != nil {
        go func() {
            ctx := context.Background()
            if err := h.EmailNotifier.SendReleaseNotification(ctx, shipmentID); err != nil {
                fmt.Printf("Warning: failed to send release notification: %v\n", err)
            }
        }()
    }
}
```

#### Task 1.4: Wire Up Delivery Confirmation
**File:** `internal/handlers/shipments.go` - `UpdateShipmentStatus()` method

**Trigger Point:** When status changes to `delivered`

**Implementation:**
```go
if newStatus == models.ShipmentStatusDelivered {
    if h.EmailNotifier != nil {
        // Only for shipment types that involve engineer delivery
        var shipmentType models.ShipmentType
        err := h.DB.QueryRowContext(r.Context(),
            `SELECT shipment_type FROM shipments WHERE id = $1`,
            shipmentID,
        ).Scan(&shipmentType)
        
        if err == nil && (shipmentType == models.ShipmentTypeSingleFullJourney || 
                         shipmentType == models.ShipmentTypeWarehouseToEngineer) {
            go func() {
                ctx := context.Background()
                if err := h.EmailNotifier.SendDeliveryConfirmation(ctx, shipmentID); err != nil {
                    fmt.Printf("Warning: failed to send delivery confirmation: %v\n", err)
                }
            }()
        }
    }
}
```

---

### Phase 2: Implement New Notifications (8-12 hours)

#### Task 2.1: Shipment Picked Up Notification (To Client)

**Step 1:** Add template data structure to `internal/email/templates.go`
```go
type ShipmentPickedUpData struct {
    ContactName       string
    ClientCompany     string
    TrackingNumber    string
    CourierName       string
    PickedUpDate      string
    ExpectedArrival   string
    TrackingURL       string
    ShipmentType      string
}
```

**Step 2:** Add template to `templates.go` `loadTemplates()` method
- Subject: "Shipment Picked Up - [Tracking Number]"
- Include tracking number, courier, expected arrival
- Add tracking URL
- Professional, informative tone

**Step 3:** Add rendering case to `RenderTemplate()` and `GetSubject()`

**Step 4:** Create notification function in `internal/email/notifier.go`
```go
func (n *Notifier) SendShipmentPickedUpNotification(ctx context.Context, shipmentID int64) error
```

**Step 5:** Wire up in `internal/handlers/shipments.go` (same location as warehouse pre-alert)

**Step 6:** Add unit test in `internal/email/notifier_test.go`

#### Task 2.2: Delivered to Engineer Notification (To Client)

**Step 1:** Add template data structure
```go
type EngineerDeliveryClientData struct {
    ContactName      string
    ClientCompany    string
    EngineerName     string
    DeliveryDate     string
    TrackingNumber   string
    JiraTicket       string
    ProjectName      string
}
```

**Step 2:** Create template - warm, success-focused message

**Step 3:** Create notification function with shipment type check

**Step 4:** Wire up in status change handler (with type check)

**Step 5:** Test with both applicable shipment types

#### Task 2.3: Pickup Form Submitted Notification (To Logistics)

**Step 1:** Add template data structure
```go
type PickupFormSubmittedData struct {
    ShipmentID       int64
    ShipmentType     string
    ClientCompany    string
    ContactName      string
    ContactEmail     string
    ContactPhone     string
    PickupAddress    string
    PickupDate       string
    NumberOfDevices  int
    JiraTicket       string
    ShipmentURL      string
}
```

**Step 2:** Create template - professional, actionable for logistics team

**Step 3:** Create notification function

**Step 4:** Wire up in `internal/handlers/pickup_form.go` after form submission (line 521)

**Step 5:** Ensure fires for all 3 shipment types

#### Task 2.4: Reception Report Approval Request (To Logistics)

**Step 1:** Add template data structure
```go
type ReceptionReportApprovalData struct {
    ShipmentID       int64
    TrackingNumber   string
    ClientCompany    string
    ReceivedDate     string
    WarehouseUser    string
    Notes            string
    PhotoURLs        []string
    SerialNumbers    []string
    ReportURL        string
    ApprovalURL      string
}
```

**Step 2:** Create template - include photos, request action

**Step 3:** Create notification function

**Step 4:** Wire up in `internal/handlers/reception_report.go:303` (replace TODO comment)

**Step 5:** Consider approval workflow (may need additional implementation)

#### Task 2.5: In Transit to Engineer Notification (To Engineer)

**Step 1:** Add template data structure
```go
type InTransitToEngineerData struct {
    EngineerName     string
    DeviceModel      string
    TrackingNumber   string
    CourierName      string
    ETA              string
    ShipmentURL      string
    ContactInfo      string
}
```

**Step 2:** Create template - helpful, informative for engineer

**Step 3:** Create notification function with ETA extraction

**Step 4:** Wire up in status change handler (with type check)

**Step 5:** Ensure ETA is properly displayed

---

### Phase 3: Configuration & Infrastructure (2-3 hours)

#### Task 3.1: Add Logistics Email Configuration

**File:** `internal/config/config.go`

Add to SMTPConfig:
```go
type SMTPConfig struct {
    Host            string
    Port            string
    User            string
    Password        string
    From            string
    FromName        string
    LogisticsEmail  string  // NEW
    WarehouseEmail  string  // NEW (for fallback)
}
```

Update Load() method:
```go
SMTP: SMTPConfig{
    // ... existing fields
    LogisticsEmail: getEnv("LOGISTICS_EMAIL", "international@bairesdev.com"),
    WarehouseEmail: getEnv("WAREHOUSE_EMAIL", "warehouse@bairesdev.com"),
},
```

#### Task 3.2: Update Notifier to Use Config

Pass config to Notifier or update to query from database with fallback to config.

#### Task 3.3: Add Helper Methods

Add to `internal/email/notifier.go`:
```go
// getLogisticsEmail retrieves the logistics team email
func (n *Notifier) getLogisticsEmail(ctx context.Context) (string, error)

// getWarehouseEmail retrieves a warehouse user email
func (n *Notifier) getWarehouseEmail(ctx context.Context) (string, error)

// getContactEmailFromForm retrieves contact email from pickup form
func (n *Notifier) getContactEmailFromForm(ctx context.Context, shipmentID int64) (string, error)
```

---

### Phase 4: Testing (4-6 hours)

#### Task 4.1: Update Email Test Script

**File:** `cmd/emailtest/main.go`

Add test functions for all 11 notifications:
1. Magic Link (existing)
2. Address Confirmation (existing)
3. Pickup Confirmation (existing)
4. Pickup Scheduled (new)
5. Warehouse Pre-Alert (existing)
6. Release Notification (existing)
7. Delivery Confirmation (existing)
8. Shipment Picked Up (new)
9. Engineer Delivery to Client (new)
10. Pickup Form Submitted to Logistics (new)
11. Reception Report Approval (new)
12. In Transit to Engineer (new)

#### Task 4.2: Create Unit Tests

**Files to create/update:**
- `internal/email/notifier_test.go` - Test all notification functions
- `internal/email/templates_test.go` - Test all template rendering
- `internal/handlers/shipments_test.go` - Test email triggers on status changes
- `internal/handlers/pickup_form_test.go` - Test email triggers on form submission
- `internal/handlers/reception_report_test.go` - Test email trigger on report creation

#### Task 4.3: Integration Testing

Test complete flows for each shipment type:

**Single Full Journey Flow:**
1. Create shipment → No email
2. Send magic link → Magic link email
3. Submit pickup form → Pickup confirmation + Form submitted to logistics
4. Status: pickup_scheduled → Pickup scheduled email
5. Status: picked_up → Shipment picked up + Warehouse pre-alert
6. Submit reception report → Reception approval to logistics
7. Status: released → Release notification to logistics
8. Status: in_transit_to_engineer → In transit to engineer
9. Status: delivered → Delivery confirmation + Delivered to client

**Bulk to Warehouse Flow:**
1-6 same as above, then stops at warehouse

**Warehouse to Engineer Flow:**
Steps 7-9 only (starts from warehouse)

#### Task 4.4: Manual Testing with Mailhog

Set up Mailhog and test all notifications:
```bash
docker-compose up mailhog
```

Access at `http://localhost:8025`

---

## Technical Considerations

### Email Template Best Practices
- Responsive design (works on mobile)
- Clear call-to-action buttons
- Professional branding
- Plain text fallback
- Proper subject lines
- Unsubscribe footer (if needed)

### Error Handling
- All email sends should be async (goroutines)
- Log failures but don't block operations
- Record all attempts in `notification_logs` table
- Consider retry mechanism for critical notifications

### Performance
- Send emails asynchronously to avoid blocking HTTP responses
- Use fresh context in goroutines (not request context)
- Consider rate limiting if sending bulk emails
- Monitor email sending performance

### Security
- Validate email addresses before sending
- Don't expose sensitive data in emails
- Use secure SMTP connection (TLS)
- Be careful with magic links (expiration, one-time use)

---

## Open Questions & Decisions Needed

### 1. Reception Report Approval Workflow
**Question:** Is there a formal approval process for reception reports?
- Should logistics approve via email link or UI?
- What happens after approval/rejection?
- Should rejection require warehouse to redo the report?

**Current State:** No approval workflow exists in database

**Recommendation:** 
- Add `approval_status` field to `reception_reports` table
- Add `approved_by_user_id` and `approved_at` fields
- Create approval endpoint
- Email includes approval link

### 2. Magic Link Auto-Send
**Question:** Should magic links be sent automatically via email?

**Current State:** Link is shown in UI, user must manually copy/send it

**Recommendation:** 
- Add checkbox: "Send via email" (default checked)
- Auto-send email when magic link created
- Still show link in UI for reference

### 3. Multiple Recipients & CC
**Question:** Should some emails CC additional recipients?

**Recommendations:**
- Logistics emails: CC warehouse manager?
- Client emails: CC logistics for visibility?
- Engineer emails: CC project manager?

### 4. Client Email Determination
**Question:** For "To Client" notifications, which email should we use?

**Options:**
1. Contact email from pickup form (most recent, accurate)
2. Client company user who created shipment
3. Both (send to contact, CC client user)

**Recommendation:** Use contact email from pickup form as primary, with option to CC client user

### 5. Email Failure Handling
**Question:** What should happen if email sending fails?

**Options:**
1. Silent failure (log only)
2. Show warning to user in UI
3. Retry mechanism
4. Queue for later sending

**Recommendation:** 
- Log failure in notification_logs with status='failed'
- Show warning in UI if critical email fails
- Retry 3 times for critical emails (magic link, pickup scheduled)

### 6. Notification Preferences
**Question:** Should users be able to configure which emails they receive?

**Future Enhancement:** Add notification preferences table

### 7. Warehouse Email Recipient
**Question:** Which warehouse email should receive pre-alerts?

**Options:**
1. Query first warehouse user from users table
2. Use configured warehouse email from config
3. Send to all warehouse users

**Recommendation:** Query warehouse users and send to all (likely 1-2 people)

---

## Database Schema Changes Needed

### Optional: Add Approval to Reception Reports

```sql
ALTER TABLE reception_reports 
ADD COLUMN approval_status VARCHAR(20) DEFAULT 'pending',
ADD COLUMN approved_by_user_id BIGINT REFERENCES users(id),
ADD COLUMN approved_at TIMESTAMP,
ADD COLUMN approval_notes TEXT;

CREATE INDEX idx_reception_reports_approval_status 
ON reception_reports(approval_status);
```

### Optional: Add Notification Preferences

```sql
CREATE TABLE notification_preferences (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    notification_type VARCHAR(50) NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, notification_type)
);
```

---

## Testing Checklist

### Before Implementation
- [ ] Review this plan
- [ ] Answer open questions
- [ ] Make decisions on recommendations
- [ ] Set up test email server (Mailhog)

### During Implementation
- [ ] Create each template
- [ ] Test template rendering
- [ ] Create notification function
- [ ] Add unit tests
- [ ] Wire up triggers
- [ ] Test with all shipment types

### After Implementation
- [ ] Run full email test suite
- [ ] Manual testing of all 11 notifications
- [ ] Test with real email addresses
- [ ] Verify all recipients receive correct emails
- [ ] Check email formatting on mobile
- [ ] Verify notification logs are created
- [ ] Test failure scenarios
- [ ] Performance testing with multiple emails

---

## Files That Will Be Modified

### Existing Files to Modify
1. `internal/email/notifier.go` - Add 5 new notification functions, fix 1 existing
2. `internal/email/templates.go` - Add 5 new templates and data structures
3. `internal/handlers/shipments.go` - Add email triggers for status changes
4. `internal/handlers/pickup_form.go` - Add email trigger for form submission
5. `internal/handlers/reception_report.go` - Add email trigger for report creation
6. `internal/config/config.go` - Add logistics/warehouse email config
7. `cmd/emailtest/main.go` - Add tests for new notifications
8. `cmd/web/main.go` - Pass config to email notifier (if needed)

### Test Files to Create/Modify
1. `internal/email/notifier_test.go` - Add tests for new functions
2. `internal/email/templates_test.go` - Add template tests
3. `internal/handlers/shipments_test.go` - Add email trigger tests
4. `internal/handlers/pickup_form_test.go` - Add email trigger tests
5. `internal/handlers/reception_report_test.go` - Add email trigger tests

### Optional New Files
1. `migrations/000021_add_reception_report_approval.up.sql`
2. `migrations/000021_add_reception_report_approval.down.sql`
3. `migrations/000022_create_notification_preferences.up.sql`
4. `migrations/000022_create_notification_preferences.down.sql`

---

## Estimated Timeline

| Phase | Tasks | Estimated Time |
|-------|-------|----------------|
| Phase 1 | Fix & wire up existing notifications | 3-4 hours |
| Phase 2 | Implement 5 new notifications | 8-12 hours |
| Phase 3 | Configuration & infrastructure | 2-3 hours |
| Phase 4 | Testing & verification | 4-6 hours |
| **Total** | | **17-25 hours** |

**Breakdown by notification:**
- Each new notification: ~2 hours (template + function + tests)
- Each existing notification wire-up: ~30-45 minutes
- Configuration and helpers: ~2 hours
- Comprehensive testing: ~4 hours

---

## Success Criteria

### Functional Requirements
- [ ] All 11 email notifications implemented and working
- [ ] Emails sent to correct recipients for each shipment type
- [ ] All templates render correctly with proper data
- [ ] Emails are sent asynchronously without blocking UI
- [ ] Notification logs created for all sent emails
- [ ] Failures logged appropriately

### Quality Requirements
- [ ] All unit tests passing
- [ ] Integration tests cover all flows
- [ ] Code follows TDD principles
- [ ] Templates are responsive and professional
- [ ] Error handling is robust
- [ ] Performance is acceptable (< 100ms to trigger)

### Documentation Requirements
- [ ] Code is well-commented
- [ ] Email templates are documented
- [ ] Configuration options documented
- [ ] Testing procedure documented

---

## Next Steps

1. **Review this plan** and answer open questions
2. **Make decisions** on recommendations
3. **Set up development environment** with Mailhog
4. **Start with Phase 1** - fix existing notifications
5. **Proceed to Phase 2** - implement new notifications
6. **Complete Phase 3** - configuration
7. **Finish with Phase 4** - comprehensive testing

---

## References

### Related Documentation
- `docs/DOCKER_EMAIL_TESTING.md` - Email testing setup
- `docs/JIRA_INTEGRATION_GUIDE.md` - JIRA integration
- `docs/DATABASE_SETUP.md` - Database configuration

### Code References
- `internal/email/` - Email implementation
- `internal/handlers/` - Event handlers
- `internal/models/shipment.go` - Shipment types and statuses
- `migrations/` - Database schema

### External Resources
- Email template best practices
- SMTP configuration guides
- HTML email design guidelines

---

**Document Version:** 1.0  
**Last Updated:** November 14, 2025  
**Author:** AI Assistant with User Collaboration

