package models

import (
	"encoding/json"
	"errors"
	"time"
)

// PickupForm represents a form submitted by client for pickup scheduling
type PickupForm struct {
	ID                int64           `json:"id" db:"id"`
	ShipmentID        int64           `json:"shipment_id" db:"shipment_id"`
	SubmittedByUserID int64           `json:"submitted_by_user_id" db:"submitted_by_user_id"`
	SubmittedAt       time.Time       `json:"submitted_at" db:"submitted_at"`
	FormData          json.RawMessage `json:"form_data" db:"form_data"`

	// Relations
	Shipment *Shipment `json:"shipment,omitempty" db:"-"`
	User     *User     `json:"user,omitempty" db:"-"`
}

// Validate validates the PickupForm model
func (p *PickupForm) Validate() error {
	if p.ShipmentID == 0 {
		return errors.New("shipment ID is required")
	}
	if p.SubmittedByUserID == 0 {
		return errors.New("submitted by user ID is required")
	}
	return nil
}

// TableName returns the table name for the PickupForm model
func (p *PickupForm) TableName() string {
	return "pickup_forms"
}

// BeforeCreate sets the timestamp before creating a pickup form
func (p *PickupForm) BeforeCreate() {
	p.SubmittedAt = time.Now()
}

// ReceptionReportStatus represents the approval status of a reception report
type ReceptionReportStatus string

// Reception report status constants
const (
	ReceptionReportStatusPendingApproval ReceptionReportStatus = "pending_approval"
	ReceptionReportStatusApproved        ReceptionReportStatus = "approved"
)

// ReceptionReport represents a report submitted by warehouse staff when receiving laptops
// Each reception report is linked to a single laptop
type ReceptionReport struct {
	ID                     int64                 `json:"id" db:"id"`
	LaptopID               int64                 `json:"laptop_id" db:"laptop_id"`
	ShipmentID             *int64                `json:"shipment_id,omitempty" db:"shipment_id"`           // Reference to original shipment
	ClientCompanyID        *int64                `json:"client_company_id,omitempty" db:"client_company_id"` // Reference for tracking
	TrackingNumber         string                `json:"tracking_number,omitempty" db:"tracking_number"`   // Reference for tracking
	WarehouseUserID        int64                 `json:"warehouse_user_id" db:"warehouse_user_id"`
	ReceivedAt             time.Time             `json:"received_at" db:"received_at"`
	Notes                  string                `json:"notes,omitempty" db:"notes"`
	
	// Required photo uploads
	PhotoSerialNumber      string                `json:"photo_serial_number" db:"photo_serial_number"`
	PhotoExternalCondition string                `json:"photo_external_condition" db:"photo_external_condition"`
	PhotoWorkingCondition  string                `json:"photo_working_condition" db:"photo_working_condition"`
	
	// Approval tracking
	Status                 ReceptionReportStatus `json:"status" db:"status"`
	ApprovedBy             *int64                `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt             *time.Time            `json:"approved_at,omitempty" db:"approved_at"`
	
	CreatedAt              time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time             `json:"updated_at" db:"updated_at"`

	// Relations
	Laptop        *Laptop        `json:"laptop,omitempty" db:"-"`
	Shipment      *Shipment      `json:"shipment,omitempty" db:"-"`
	ClientCompany *ClientCompany `json:"client_company,omitempty" db:"-"`
	User          *User          `json:"user,omitempty" db:"-"`
	Approver      *User          `json:"approver,omitempty" db:"-"`
}

// Validate validates the ReceptionReport model
func (r *ReceptionReport) Validate() error {
	if r.LaptopID == 0 {
		return errors.New("laptop ID is required")
	}
	if r.WarehouseUserID == 0 {
		return errors.New("warehouse user ID is required")
	}
	if r.PhotoSerialNumber == "" {
		return errors.New("photo of serial number is required")
	}
	if r.PhotoExternalCondition == "" {
		return errors.New("photo of external condition is required")
	}
	if r.PhotoWorkingCondition == "" {
		return errors.New("photo of working condition is required")
	}
	return nil
}

// TableName returns the table name for the ReceptionReport model
func (r *ReceptionReport) TableName() string {
	return "reception_reports"
}

// BeforeCreate sets the timestamps and default status before creating a reception report
func (r *ReceptionReport) BeforeCreate() {
	now := time.Now()
	r.ReceivedAt = now
	r.CreatedAt = now
	r.UpdatedAt = now
	if r.Status == "" {
		r.Status = ReceptionReportStatusPendingApproval
	}
}

// BeforeUpdate sets the updated_at timestamp before updating a reception report
func (r *ReceptionReport) BeforeUpdate() {
	r.UpdatedAt = time.Now()
}

// IsPendingApproval returns true if the reception report is pending approval
func (r *ReceptionReport) IsPendingApproval() bool {
	return r.Status == ReceptionReportStatusPendingApproval
}

// IsApproved returns true if the reception report is approved
func (r *ReceptionReport) IsApproved() bool {
	return r.Status == ReceptionReportStatusApproved
}

// Approve marks the reception report as approved by a logistics user
func (r *ReceptionReport) Approve(logisticsUserID int64) {
	r.Status = ReceptionReportStatusApproved
	r.ApprovedBy = &logisticsUserID
	now := time.Now()
	r.ApprovedAt = &now
	r.BeforeUpdate()
}

// DeliveryForm represents a form submitted when a laptop is delivered to an engineer
type DeliveryForm struct {
	ID          int64     `json:"id" db:"id"`
	ShipmentID  int64     `json:"shipment_id" db:"shipment_id"`
	EngineerID  int64     `json:"engineer_id" db:"engineer_id"`
	DeliveredAt time.Time `json:"delivered_at" db:"delivered_at"`
	Notes       string    `json:"notes,omitempty" db:"notes"`
	PhotoURLs   []string  `json:"photo_urls,omitempty" db:"photo_urls"`

	// Relations
	Shipment *Shipment         `json:"shipment,omitempty" db:"-"`
	Engineer *SoftwareEngineer `json:"engineer,omitempty" db:"-"`
}

// Validate validates the DeliveryForm model
func (d *DeliveryForm) Validate() error {
	if d.ShipmentID == 0 {
		return errors.New("shipment ID is required")
	}
	if d.EngineerID == 0 {
		return errors.New("engineer ID is required")
	}
	return nil
}

// TableName returns the table name for the DeliveryForm model
func (d *DeliveryForm) TableName() string {
	return "delivery_forms"
}

// BeforeCreate sets the timestamp before creating a delivery form
func (d *DeliveryForm) BeforeCreate() {
	d.DeliveredAt = time.Now()
}

// HasPhotos returns true if the delivery form has photos
func (d *DeliveryForm) HasPhotos() bool {
	return len(d.PhotoURLs) > 0
}

