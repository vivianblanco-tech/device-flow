# Email Notifications Implementation - Continuation Prompt

Use this prompt to continue the email notifications implementation:

---

## PROMPT FOR AI ASSISTANT

I need to implement email notifications for Align. A comprehensive analysis and implementation plan has already been completed and documented in `docs/EMAIL_NOTIFICATIONS_IMPLEMENTATION_PLAN.md`.

### Context

The application is a Go web application using:
- **Backend:** Go with `database/sql` and PostgreSQL
- **Email:** Custom SMTP client in `internal/email/`
- **Templates:** HTML templates with Go `html/template`
- **Handlers:** HTTP handlers in `internal/handlers/`

### Current State

- **6 email notification templates exist** but only 2 are fully wired up
- **5 additional notifications** need to be created from scratch
- Analysis is complete and documented in the implementation plan

### Requirements

The application needs the following 11 email notifications:

**To Client:**
1. Magic Link Email (when created) ✅ EXISTS
2. Pickup Confirmation (form submitted) ⚠️ EXISTS BUT BROKEN
3. Pickup Scheduled Confirmation ✅ WORKING
4. Shipment Picked Up (send tracking) ❌ NEEDS CREATION
5. Delivered to Engineer notification ❌ NEEDS CREATION

**To Warehouse:**
6. Pre-alert for Incoming Shipment ⚠️ EXISTS, NOT WIRED UP

**To Logistics (international@bairesdev.com):**
7. Pickup Form Submitted ❌ NEEDS CREATION
8. Reception Report Approval Request ❌ NEEDS CREATION
9. Hardware Release for Pickup ⚠️ EXISTS, NOT WIRED UP

**To Software Engineer:**
10. Delivery Confirmation ⚠️ EXISTS, NOT WIRED UP
11. In Transit to Engineer (with ETA) ❌ NEEDS CREATION

### Shipment Types

The system has 3 shipment types with different flows:
1. **Single Full Journey** - Client → Warehouse → Engineer (all notifications apply)
2. **Bulk to Warehouse** - Client → Warehouse (stops at warehouse)
3. **Warehouse to Engineer** - Warehouse → Engineer (no pickup notifications)

### What I Need

Please implement the email notifications according to the plan in `docs/EMAIL_NOTIFICATIONS_IMPLEMENTATION_PLAN.md`.

**Implementation Order:**
1. **Phase 1:** Fix and wire up existing notifications (Tasks 1.1-1.4)
2. **Phase 2:** Implement new notifications (Tasks 2.1-2.5)
3. **Phase 3:** Configuration and infrastructure (Tasks 3.1-3.3)
4. **Phase 4:** Testing (Tasks 4.1-4.4)

### Before You Start

Please review the implementation plan and:
1. Ask me to clarify any of the "Open Questions" in the plan
2. Confirm the approach for any recommendations that need decisions
3. Let me know if you need any additional information

### Key Files to Work With

**Email System:**
- `internal/email/notifier.go` - Notification functions
- `internal/email/templates.go` - Email templates
- `internal/email/client.go` - SMTP client

**Handlers (where emails are triggered):**
- `internal/handlers/shipments.go` - Status change handlers
- `internal/handlers/pickup_form.go` - Form submission handlers
- `internal/handlers/reception_report.go` - Reception report handlers

**Configuration:**
- `internal/config/config.go` - Application config
- `cmd/web/main.go` - Application initialization

**Testing:**
- `cmd/emailtest/main.go` - Email testing utility
- `internal/email/*_test.go` - Unit tests

### Testing Setup

Email testing uses Mailhog (SMTP mock server):
```bash
# Start Mailhog
docker-compose up mailhog

# Access web UI
http://localhost:8025

# Run email tests
go run cmd/emailtest/main.go
```

### Guidelines

- **Use TDD methodology** - write tests first
- **Follow existing patterns** in the codebase
- **Send emails asynchronously** using goroutines with fresh context
- **Log all notification attempts** to `notification_logs` table
- **Handle errors gracefully** - log but don't fail operations
- **Consider shipment types** - not all notifications apply to all types

### Success Criteria

- All 11 notifications working correctly
- Proper recipients for each notification type
- All email triggers wired up to correct events
- Comprehensive tests with >80% coverage
- All templates render correctly
- Notification logs created for tracking

### Questions to Answer First

Before implementing, please help me decide on these open questions from the plan:

1. **Reception Report Approval:** Should we implement a formal approval workflow in the database, or just send notification emails?

2. **Magic Link Sending:** Should magic links be sent automatically via email, or continue showing in UI for manual sending?

3. **Client Email Selection:** For "To Client" emails, should we use:
   - Contact email from pickup form (most current)
   - Client company user who created shipment
   - Both (primary + CC)

4. **Email Failure Handling:** What should happen when email sending fails?
   - Silent failure with logging only
   - Show warning in UI
   - Retry mechanism
   - Combination of above

Please review `docs/EMAIL_NOTIFICATIONS_IMPLEMENTATION_PLAN.md` and let me know:
1. Any questions or clarifications you need
2. Your recommendations for the open questions
3. Your approach to implementation

Then we can proceed with the implementation phase by phase.

---

## ADDITIONAL CONTEXT (if needed)

### Example: How Current Notifications Work

From `internal/handlers/shipments.go` (lines 544-569):
```go
// Send email notification if status changed to pickup_scheduled
notificationSent := false
if oldStatus == string(models.ShipmentStatusPendingPickup) && newStatus == models.ShipmentStatusPickupScheduled {
    if h.EmailNotifier != nil {
        // Check if pickup form exists before sending notification
        var pickupFormExists bool
        err := h.DB.QueryRowContext(r.Context(),
            `SELECT EXISTS(SELECT 1 FROM pickup_forms WHERE shipment_id = $1)`,
            shipmentID,
        ).Scan(&pickupFormExists)
        
        if err == nil && pickupFormExists {
            go func() {
                // Use a fresh context for the background goroutine
                ctx := context.Background()
                if err := h.EmailNotifier.SendPickupScheduledNotification(ctx, shipmentID); err != nil {
                    fmt.Printf("Warning: failed to send pickup scheduled notification: %v\n", err)
                } else {
                    fmt.Printf("Pickup scheduled notification sent successfully for shipment %d\n", shipmentID)
                }
            }()
            notificationSent = true
        }
    }
}
```

This pattern should be followed for new notification triggers.

### Shipment Status Flow

```
Single Full Journey:
pending_pickup_from_client → pickup_from_client_scheduled → 
picked_up_from_client → in_transit_to_warehouse → 
at_warehouse → released_from_warehouse → 
in_transit_to_engineer → delivered

Bulk to Warehouse:
pending_pickup_from_client → pickup_from_client_scheduled → 
picked_up_from_client → in_transit_to_warehouse → 
at_warehouse (ENDS)

Warehouse to Engineer:
released_from_warehouse → in_transit_to_engineer → 
delivered (STARTS)
```

Each status transition is an opportunity to trigger email notifications.

---

**Ready to begin?** Start by reviewing the plan and asking any clarifying questions!

