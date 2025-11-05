package views

import "github.com/yourusername/laptop-tracking-system/internal/models"

// NavigationLinks represents which navigation links should be visible to a user
type NavigationLinks struct {
	Dashboard        bool
	Shipments        bool
	Inventory        bool
	Calendar         bool
	PickupForms      bool
	ReceptionReports bool
}

// GetNavigationLinks returns the navigation links visible to a user based on their role
func GetNavigationLinks(role models.UserRole) NavigationLinks {
	nav := NavigationLinks{}

	switch role {
	case models.RoleLogistics:
		// Logistics has full access to all navigation links
		nav.Dashboard = true
		nav.Shipments = true
		nav.Inventory = true
		nav.Calendar = true
		nav.PickupForms = true
		nav.ReceptionReports = true

	case models.RoleProjectManager:
		// Project Manager has access to dashboards and reports
		nav.Dashboard = true
		nav.Shipments = true
		nav.Inventory = true
		nav.Calendar = true
		nav.PickupForms = false
		nav.ReceptionReports = false

	case models.RoleWarehouse:
		// Warehouse has access to inventory and reception
		nav.Dashboard = false
		nav.Shipments = true
		nav.Inventory = true
		nav.Calendar = true
		nav.PickupForms = false
		nav.ReceptionReports = true

	case models.RoleClient:
		// Client has limited access
		nav.Dashboard = false
		nav.Shipments = true
		nav.Inventory = false
		nav.Calendar = true
		nav.PickupForms = true
		nav.ReceptionReports = false
	}

	return nav
}

// HasAnyLink returns true if at least one navigation link is visible
func (n NavigationLinks) HasAnyLink() bool {
	return n.Dashboard || n.Shipments || n.Inventory || n.Calendar || n.PickupForms || n.ReceptionReports
}

