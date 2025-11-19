package utils

import (
	"testing"
)

func TestFormatContactInfoForForm(t *testing.T) {
	tests := []struct {
		name        string
		contactInfo string
		expected    string
		description string
	}{
		{
			name:        "empty contact info",
			contactInfo: "",
			expected:    "",
			description: "Empty string should return empty string",
		},
		{
			name:        "plain text contact info",
			contactInfo: "Email: contact@example.com\nPhone: +1-555-0100",
			expected:    "Email: contact@example.com\nPhone: +1-555-0100",
			description: "Plain text should be returned as-is",
		},
		{
			name:        "JSON contact info with all fields",
			contactInfo: `{"email":"contact@example.com","phone":"+1-555-0100","address":"123 Main St","country":"USA","website":"https://example.com"}`,
			expected:    "Email: contact@example.com\nPhone: +1-555-0100\nAddress: 123 Main St\nCountry: USA\nWebsite: https://example.com",
			description: "JSON should be converted to readable plain text format",
		},
		{
			name:        "JSON contact info with partial fields",
			contactInfo: `{"email":"contact@example.com","phone":"+1-555-0100"}`,
			expected:    "Email: contact@example.com\nPhone: +1-555-0100",
			description: "JSON with only some fields should format correctly",
		},
		{
			name:        "JSON contact info with empty fields",
			contactInfo: `{"email":"contact@example.com","phone":"","address":"123 Main St"}`,
			expected:    "Email: contact@example.com\nAddress: 123 Main St",
			description: "Empty fields in JSON should be skipped",
		},
		{
			name:        "invalid JSON should return as-is",
			contactInfo: `{"email":"contact@example.com","phone":}`,
			expected:    `{"email":"contact@example.com","phone":}`,
			description: "Invalid JSON should be treated as plain text",
		},
		{
			name:        "JSON with only email",
			contactInfo: `{"email":"contact@example.com"}`,
			expected:    "Email: contact@example.com",
			description: "JSON with single field should format correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatContactInfoForForm(tt.contactInfo)
			if result != tt.expected {
				t.Errorf("FormatContactInfoForForm() = %q, want %q. %s", result, tt.expected, tt.description)
			}
		})
	}
}
