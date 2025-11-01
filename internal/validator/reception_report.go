package validator

import (
	"errors"
	"net/url"
	"strings"
)

// ReceptionReportInput represents the input data for a warehouse reception report
type ReceptionReportInput struct {
	ShipmentID      int64    `json:"shipment_id"`
	WarehouseUserID int64    `json:"warehouse_user_id"`
	Notes           string   `json:"notes"`
	PhotoURLs       []string `json:"photo_urls"`
}

const (
	// MaxReceptionReportNotes is the maximum length for reception report notes
	MaxReceptionReportNotes = 1000
	// MaxReceptionPhotos is the maximum number of photos that can be uploaded
	MaxReceptionPhotos = 10
)

// ValidateReceptionReport validates the reception report input
func ValidateReceptionReport(input ReceptionReportInput) error {
	// Validate shipment ID
	if input.ShipmentID == 0 {
		return errors.New("shipment ID is required")
	}

	// Validate warehouse user ID
	if input.WarehouseUserID == 0 {
		return errors.New("warehouse user ID is required")
	}

	// Validate notes length
	if len(input.Notes) > MaxReceptionReportNotes {
		return errors.New("notes must not exceed 1000 characters")
	}

	// Validate photo count
	if len(input.PhotoURLs) > MaxReceptionPhotos {
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

// isValidURL validates URL format
func isValidURL(rawURL string) bool {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return false
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	// Must have a scheme (http or https) and host
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	if u.Host == "" {
		return false
	}

	// Check for spaces (invalid in URLs)
	if strings.Contains(rawURL, " ") {
		return false
	}

	return true
}

