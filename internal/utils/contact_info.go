package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatContactInfoForForm converts JSON contact info to plain text format for form display.
// If the input is not valid JSON, it returns the input as-is.
// The output format is: "Email: ...\nPhone: ...\nAddress: ..." etc.
func FormatContactInfoForForm(contactInfo string) string {
	if contactInfo == "" {
		return ""
	}

	// Try to parse as JSON
	var contactMap map[string]interface{}
	if err := json.Unmarshal([]byte(contactInfo), &contactMap); err != nil {
		// If not JSON, return as-is
		return contactInfo
	}

	// Build formatted plain text
	var parts []string
	fieldOrder := []string{"email", "phone", "address", "country", "website"}
	fieldLabels := map[string]string{
		"email":   "Email",
		"phone":   "Phone",
		"address": "Address",
		"country": "Country",
		"website": "Website",
	}

	for _, field := range fieldOrder {
		if value, ok := contactMap[field].(string); ok && value != "" {
			label := fieldLabels[field]
			parts = append(parts, fmt.Sprintf("%s: %s", label, value))
		}
	}

	return strings.Join(parts, "\n")
}

