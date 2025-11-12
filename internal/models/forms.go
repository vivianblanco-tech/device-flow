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

// ReceptionReport represents a report submitted by warehouse staff when receiving laptops
type ReceptionReport struct {
	ID              int64     `json:"id" db:"id"`
	ShipmentID      int64     `json:"shipment_id" db:"shipment_id"`
	WarehouseUserID int64     `json:"warehouse_user_id" db:"warehouse_user_id"`
	ReceivedAt      time.Time `json:"received_at" db:"received_at"`
	Notes           string    `json:"notes,omitempty" db:"notes"`
	PhotoURLs       []string  `json:"photo_urls,omitempty" db:"photo_urls"`
	
	// Serial number correction tracking
	ExpectedSerialNumber  string  `json:"expected_serial_number,omitempty" db:"expected_serial_number"`
	ActualSerialNumber    string  `json:"actual_serial_number,omitempty" db:"actual_serial_number"`
	SerialNumberCorrected bool    `json:"serial_number_corrected" db:"serial_number_corrected"`
	CorrectionNote        string  `json:"correction_note,omitempty" db:"correction_note"`
	CorrectionApprovedBy  *int64  `json:"correction_approved_by,omitempty" db:"correction_approved_by"`

	// Relations
	Shipment *Shipment `json:"shipment,omitempty" db:"-"`
	User     *User     `json:"user,omitempty" db:"-"`
}

// Validate validates the ReceptionReport model
func (r *ReceptionReport) Validate() error {
	if r.ShipmentID == 0 {
		return errors.New("shipment ID is required")
	}
	if r.WarehouseUserID == 0 {
		return errors.New("warehouse user ID is required")
	}
	return nil
}

// TableName returns the table name for the ReceptionReport model
func (r *ReceptionReport) TableName() string {
	return "reception_reports"
}

// BeforeCreate sets the timestamp before creating a reception report
func (r *ReceptionReport) BeforeCreate() {
	r.ReceivedAt = time.Now()
}

// HasPhotos returns true if the reception report has photos
func (r *ReceptionReport) HasPhotos() bool {
	return len(r.PhotoURLs) > 0
}

// HasSerialNumberCorrection returns true if serial number was corrected
func (r *ReceptionReport) HasSerialNumberCorrection() bool {
	return r.SerialNumberCorrected
}

// SerialNumberCorrectionNote returns the correction note
func (r *ReceptionReport) SerialNumberCorrectionNote() string {
	return r.CorrectionNote
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

