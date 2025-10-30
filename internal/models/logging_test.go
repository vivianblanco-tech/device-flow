package models

import (
	"encoding/json"
	"testing"
)

// NotificationLog tests

func TestNotificationLog_Validate(t *testing.T) {
	tests := []struct {
		name    string
		log     NotificationLog
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid notification log",
			log: NotificationLog{
				ShipmentID: int64Ptr(1),
				Type:       "email",
				Recipient:  "user@example.com",
				Status:     "sent",
			},
			wantErr: false,
		},
		{
			name: "valid - no shipment ID",
			log: NotificationLog{
				Type:      "email",
				Recipient: "admin@example.com",
				Status:    "sent",
			},
			wantErr: false,
		},
		{
			name: "invalid - missing type",
			log: NotificationLog{
				ShipmentID: int64Ptr(1),
				Recipient:  "user@example.com",
				Status:     "sent",
			},
			wantErr: true,
			errMsg:  "notification type is required",
		},
		{
			name: "invalid - missing recipient",
			log: NotificationLog{
				ShipmentID: int64Ptr(1),
				Type:       "email",
				Status:     "sent",
			},
			wantErr: true,
			errMsg:  "recipient is required",
		},
		{
			name: "invalid - missing status",
			log: NotificationLog{
				ShipmentID: int64Ptr(1),
				Type:       "email",
				Recipient:  "user@example.com",
			},
			wantErr: true,
			errMsg:  "status is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.log.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("NotificationLog.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("NotificationLog.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestNotificationLog_TableName(t *testing.T) {
	log := NotificationLog{}
	expected := "notification_logs"
	if got := log.TableName(); got != expected {
		t.Errorf("NotificationLog.TableName() = %v, want %v", got, expected)
	}
}

func TestNotificationLog_BeforeCreate(t *testing.T) {
	log := &NotificationLog{
		Type:      "email",
		Recipient: "user@example.com",
		Status:    "sent",
	}

	log.BeforeCreate()

	if log.SentAt.IsZero() {
		t.Error("NotificationLog.BeforeCreate() did not set SentAt")
	}
}

func TestNotificationLog_IsSent(t *testing.T) {
	tests := []struct {
		name     string
		log      NotificationLog
		expected bool
	}{
		{
			name: "sent",
			log: NotificationLog{
				Status: "sent",
			},
			expected: true,
		},
		{
			name: "pending",
			log: NotificationLog{
				Status: "pending",
			},
			expected: false,
		},
		{
			name: "failed",
			log: NotificationLog{
				Status: "failed",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.log.IsSent(); got != tt.expected {
				t.Errorf("NotificationLog.IsSent() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// AuditLog tests

func TestAuditLog_Validate(t *testing.T) {
	tests := []struct {
		name    string
		log     AuditLog
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid audit log",
			log: AuditLog{
				UserID:     1,
				Action:     "create",
				EntityType: "shipment",
				EntityID:   10,
				Details:    json.RawMessage(`{"status": "pending"}`),
			},
			wantErr: false,
		},
		{
			name: "valid - minimal fields",
			log: AuditLog{
				UserID:     1,
				Action:     "update",
				EntityType: "laptop",
				EntityID:   20,
			},
			wantErr: false,
		},
		{
			name: "invalid - missing user ID",
			log: AuditLog{
				Action:     "delete",
				EntityType: "user",
				EntityID:   5,
			},
			wantErr: true,
			errMsg:  "user ID is required",
		},
		{
			name: "invalid - missing action",
			log: AuditLog{
				UserID:     1,
				EntityType: "shipment",
				EntityID:   10,
			},
			wantErr: true,
			errMsg:  "action is required",
		},
		{
			name: "invalid - missing entity type",
			log: AuditLog{
				UserID:   1,
				Action:   "create",
				EntityID: 10,
			},
			wantErr: true,
			errMsg:  "entity type is required",
		},
		{
			name: "invalid - missing entity ID",
			log: AuditLog{
				UserID:     1,
				Action:     "create",
				EntityType: "shipment",
			},
			wantErr: true,
			errMsg:  "entity ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.log.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AuditLog.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("AuditLog.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestAuditLog_TableName(t *testing.T) {
	log := AuditLog{}
	expected := "audit_logs"
	if got := log.TableName(); got != expected {
		t.Errorf("AuditLog.TableName() = %v, want %v", got, expected)
	}
}

func TestAuditLog_BeforeCreate(t *testing.T) {
	log := &AuditLog{
		UserID:     1,
		Action:     "create",
		EntityType: "shipment",
		EntityID:   10,
	}

	log.BeforeCreate()

	if log.Timestamp.IsZero() {
		t.Error("AuditLog.BeforeCreate() did not set Timestamp")
	}
}

func TestAuditLog_GetFormattedAction(t *testing.T) {
	tests := []struct {
		name     string
		log      AuditLog
		expected string
	}{
		{
			name: "create action",
			log: AuditLog{
				Action:     "create",
				EntityType: "shipment",
			},
			expected: "Created shipment",
		},
		{
			name: "update action",
			log: AuditLog{
				Action:     "update",
				EntityType: "laptop",
			},
			expected: "Updated laptop",
		},
		{
			name: "delete action",
			log: AuditLog{
				Action:     "delete",
				EntityType: "user",
			},
			expected: "Deleted user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.log.GetFormattedAction(); got != tt.expected {
				t.Errorf("AuditLog.GetFormattedAction() = %v, want %v", got, tt.expected)
			}
		})
	}
}

