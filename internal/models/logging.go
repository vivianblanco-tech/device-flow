package models

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// NotificationLog represents a log entry for sent notifications
type NotificationLog struct {
	ID         int64      `json:"id" db:"id"`
	ShipmentID *int64     `json:"shipment_id,omitempty" db:"shipment_id"`
	Type       string     `json:"type" db:"type"`
	Recipient  string     `json:"recipient" db:"recipient"`
	SentAt     time.Time  `json:"sent_at" db:"sent_at"`
	Status     string     `json:"status" db:"status"`

	// Relations
	Shipment *Shipment `json:"shipment,omitempty" db:"-"`
}

// Validate validates the NotificationLog model
func (n *NotificationLog) Validate() error {
	if n.Type == "" {
		return errors.New("notification type is required")
	}
	if n.Recipient == "" {
		return errors.New("recipient is required")
	}
	if n.Status == "" {
		return errors.New("status is required")
	}
	return nil
}

// TableName returns the table name for the NotificationLog model
func (n *NotificationLog) TableName() string {
	return "notification_logs"
}

// BeforeCreate sets the timestamp before creating a notification log
func (n *NotificationLog) BeforeCreate() {
	n.SentAt = time.Now()
}

// IsSent returns true if the notification was sent successfully
func (n *NotificationLog) IsSent() bool {
	return n.Status == "sent"
}

// AuditLog represents an audit trail entry for important actions
type AuditLog struct {
	ID         int64           `json:"id" db:"id"`
	UserID     int64           `json:"user_id" db:"user_id"`
	Action     string          `json:"action" db:"action"`
	EntityType string          `json:"entity_type" db:"entity_type"`
	EntityID   int64           `json:"entity_id" db:"entity_id"`
	Timestamp  time.Time       `json:"timestamp" db:"timestamp"`
	Details    json.RawMessage `json:"details,omitempty" db:"details"`

	// Relations
	User *User `json:"user,omitempty" db:"-"`
}

// Validate validates the AuditLog model
func (a *AuditLog) Validate() error {
	if a.UserID == 0 {
		return errors.New("user ID is required")
	}
	if a.Action == "" {
		return errors.New("action is required")
	}
	if a.EntityType == "" {
		return errors.New("entity type is required")
	}
	if a.EntityID == 0 {
		return errors.New("entity ID is required")
	}
	return nil
}

// TableName returns the table name for the AuditLog model
func (a *AuditLog) TableName() string {
	return "audit_logs"
}

// BeforeCreate sets the timestamp before creating an audit log
func (a *AuditLog) BeforeCreate() {
	a.Timestamp = time.Now()
}

// GetFormattedAction returns a human-readable formatted action string
func (a *AuditLog) GetFormattedAction() string {
	action := strings.Title(a.Action)
	entityType := strings.ToLower(a.EntityType)
	return action + "d " + entityType
}

