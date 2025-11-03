package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// DashboardHandler handles dashboard-related requests
type DashboardHandler struct {
	DB        *sql.DB
	Templates *template.Template
}

// NewDashboardHandler creates a new DashboardHandler
func NewDashboardHandler(db *sql.DB, templates *template.Template) *DashboardHandler {
	return &DashboardHandler{
		DB:        db,
		Templates: templates,
	}
}

// Dashboard displays the main dashboard with statistics
func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only logistics and project manager users can access the dashboard
	if user.Role != models.RoleLogistics && user.Role != models.RoleProjectManager {
		http.Error(w, "Forbidden: Only logistics and project manager users can access this page", http.StatusForbidden)
		return
	}

	// Get dashboard statistics
	stats, err := models.GetDashboardStats(h.DB)
	if err != nil {
		log.Printf("Error getting dashboard stats: %v", err)
		http.Error(w, "Failed to load dashboard statistics", http.StatusInternalServerError)
		return
	}

	// Prepare template data
	data := map[string]interface{}{
		"User":  user,
		"Stats": stats,
	}

	// Execute template using pre-parsed global templates
	if err := h.Templates.ExecuteTemplate(w, "dashboard-with-charts.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render dashboard", http.StatusInternalServerError)
		return
	}
}

// getStatusColor returns a Tailwind CSS class for the given shipment status
func getStatusColor(status models.ShipmentStatus) string {
	switch status {
	case models.ShipmentStatusPendingPickup:
		return "bg-yellow-400"
	case models.ShipmentStatusPickedUpFromClient:
		return "bg-orange-400"
	case models.ShipmentStatusInTransitToWarehouse:
		return "bg-purple-400"
	case models.ShipmentStatusAtWarehouse:
		return "bg-indigo-400"
	case models.ShipmentStatusReleasedFromWarehouse:
		return "bg-blue-400"
	case models.ShipmentStatusInTransitToEngineer:
		return "bg-cyan-400"
	case models.ShipmentStatusDelivered:
		return "bg-green-400"
	default:
		return "bg-gray-400"
	}
}

// getLaptopStatusColor returns a Tailwind CSS class for the given laptop status
func getLaptopStatusColor(status models.LaptopStatus) string {
	switch status {
	case models.LaptopStatusAvailable:
		return "bg-green-400"
	case models.LaptopStatusInTransitToWarehouse:
		return "bg-purple-400"
	case models.LaptopStatusAtWarehouse:
		return "bg-indigo-400"
	case models.LaptopStatusInTransitToEngineer:
		return "bg-cyan-400"
	case models.LaptopStatusDelivered:
		return "bg-blue-400"
	case models.LaptopStatusRetired:
		return "bg-gray-400"
	default:
		return "bg-gray-400"
	}
}

