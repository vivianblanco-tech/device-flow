package validator

import (
	"errors"
)

// DeliveryFormInput represents the input data for a delivery form
type DeliveryFormInput struct {
	ShipmentID int64    `json:"shipment_id"`
	EngineerID int64    `json:"engineer_id"`
	Notes      string   `json:"notes"`
	PhotoURLs  []string `json:"photo_urls"`
}

const (
	// MaxDeliveryFormNotes is the maximum length for delivery form notes
	MaxDeliveryFormNotes = 1000
	// MaxDeliveryPhotos is the maximum number of photos that can be uploaded
	MaxDeliveryPhotos = 10
)

// ValidateDeliveryForm validates the delivery form input
func ValidateDeliveryForm(input DeliveryFormInput) error {
	// Validate shipment ID
	if input.ShipmentID == 0 {
		return errors.New("shipment ID is required")
	}

	// Validate engineer ID
	if input.EngineerID == 0 {
		return errors.New("engineer ID is required")
	}

	// Validate notes length
	if len(input.Notes) > MaxDeliveryFormNotes {
		return errors.New("notes must not exceed 1000 characters")
	}

	// Validate photo count
	if len(input.PhotoURLs) > MaxDeliveryPhotos {
		return errors.New("cannot upload more than 10 photos")
	}

	// Validate each photo URL
	for _, photoURL := range input.PhotoURLs {
		if !isValidURL(photoURL) {
			return errors.New("invalid photo URL format")
		}
	}

	return nil
}

