package views

import (
	"testing"

	"github.com/yourusername/laptop-tracking-system/internal/models"
)

func TestGetNavigationLinks(t *testing.T) {
	tests := []struct {
		name          string
		userRole      models.UserRole
		expectedLinks map[string]bool
	}{
		{
			name:     "logistics user has access to all links",
			userRole: models.RoleLogistics,
			expectedLinks: map[string]bool{
				"dashboard":         true,
				"shipments":         true,
				"inventory":         true,
				"calendar":          true,
				"pickup_forms":      true,
				"reception_reports": true,
				"magic_links":       true,
			},
		},
		{
			name:     "project manager has access to dashboards and reports",
			userRole: models.RoleProjectManager,
			expectedLinks: map[string]bool{
				"dashboard":         true,
				"shipments":         true,
				"inventory":         true,
				"calendar":          true,
				"pickup_forms":      false,
				"reception_reports": false,
				"magic_links":       false,
			},
		},
		{
			name:     "warehouse user has access to inventory and reception",
			userRole: models.RoleWarehouse,
			expectedLinks: map[string]bool{
				"dashboard":         false,
				"shipments":         true,
				"inventory":         true,
				"calendar":          true,
				"pickup_forms":      false,
				"reception_reports": true,
				"magic_links":       false,
			},
		},
		{
			name:     "client user has limited access",
			userRole: models.RoleClient,
			expectedLinks: map[string]bool{
				"dashboard":         false,
				"shipments":         true,
				"inventory":         false,
				"calendar":          true,
				"pickup_forms":      true,
				"reception_reports": false,
				"magic_links":       false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nav := GetNavigationLinks(tt.userRole)

			// Check each expected link
			if nav.Dashboard != tt.expectedLinks["dashboard"] {
				t.Errorf("Dashboard visibility = %v, want %v", nav.Dashboard, tt.expectedLinks["dashboard"])
			}
			if nav.Shipments != tt.expectedLinks["shipments"] {
				t.Errorf("Shipments visibility = %v, want %v", nav.Shipments, tt.expectedLinks["shipments"])
			}
			if nav.Inventory != tt.expectedLinks["inventory"] {
				t.Errorf("Inventory visibility = %v, want %v", nav.Inventory, tt.expectedLinks["inventory"])
			}
			if nav.Calendar != tt.expectedLinks["calendar"] {
				t.Errorf("Calendar visibility = %v, want %v", nav.Calendar, tt.expectedLinks["calendar"])
			}
			if nav.PickupForms != tt.expectedLinks["pickup_forms"] {
				t.Errorf("PickupForms visibility = %v, want %v", nav.PickupForms, tt.expectedLinks["pickup_forms"])
			}
			if nav.ReceptionReports != tt.expectedLinks["reception_reports"] {
				t.Errorf("ReceptionReports visibility = %v, want %v", nav.ReceptionReports, tt.expectedLinks["reception_reports"])
			}
			if nav.MagicLinks != tt.expectedLinks["magic_links"] {
				t.Errorf("MagicLinks visibility = %v, want %v", nav.MagicLinks, tt.expectedLinks["magic_links"])
			}
		})
	}
}

func TestNavigationLinks_HasAnyLink(t *testing.T) {
	tests := []struct {
		name     string
		nav      NavigationLinks
		expected bool
	}{
		{
			name: "has at least one link",
			nav: NavigationLinks{
				Dashboard: true,
				Shipments: false,
			},
			expected: true,
		},
		{
			name: "has no links",
			nav: NavigationLinks{
				Dashboard:        false,
				Shipments:        false,
				Inventory:        false,
				Calendar:         false,
				PickupForms:      false,
				ReceptionReports: false,
			},
			expected: false,
		},
		{
			name: "all links enabled",
			nav: NavigationLinks{
				Dashboard:        true,
				Shipments:        true,
				Inventory:        true,
				Calendar:         true,
				PickupForms:      true,
				ReceptionReports: true,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.nav.HasAnyLink()
			if result != tt.expected {
				t.Errorf("HasAnyLink() = %v, want %v", result, tt.expected)
			}
		})
	}
}
