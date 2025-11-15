package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
	"github.com/yourusername/laptop-tracking-system/internal/views"
)

// InventoryHandler handles inventory-related requests
type InventoryHandler struct {
	DB        *sql.DB
	Templates *template.Template
}

// NewInventoryHandler creates a new InventoryHandler
func NewInventoryHandler(db *sql.DB, templates *template.Template) *InventoryHandler {
	return &InventoryHandler{
		DB:        db,
		Templates: templates,
	}
}

// InventoryList displays the inventory list with search and filter options
func (h *InventoryHandler) InventoryList(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse query parameters
	searchQuery := r.URL.Query().Get("search")
	statusFilter := r.URL.Query().Get("status")

	// Build filter
	filter := &models.LaptopFilter{
		Search:   searchQuery,
		UserRole: user.Role, // Apply role-based filtering
	}

	if statusFilter != "" {
		filter.Status = models.LaptopStatus(statusFilter)
	}

	// Get laptops
	laptops, err := models.GetAllLaptops(h.DB, filter)
	if err != nil {
		log.Printf("Error getting laptops: %v", err)
		http.Error(w, "Failed to load inventory", http.StatusInternalServerError)
		return
	}

	// Prepare template data
	data := map[string]interface{}{
		"User":         user,
		"Nav":          views.GetNavigationLinks(user.Role),
		"CurrentPage":  "inventory",
		"Laptops":      laptops,
		"SearchQuery":  searchQuery,
		"StatusFilter": statusFilter,
		"Statuses":     models.GetAllowedStatusesForRole(user.Role), // Filter statuses by user role
	}

	// Execute template using pre-parsed global templates
	if err := h.Templates.ExecuteTemplate(w, "inventory-list.html", data); err != nil {
		log.Printf("Error executing inventory template: %v", err)
		http.Error(w, "Failed to render inventory", http.StatusInternalServerError)
		return
	}
}

// LaptopDetail displays details of a specific laptop
func (h *InventoryHandler) LaptopDetail(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get laptop ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid laptop ID", http.StatusBadRequest)
		return
	}

	// Get laptop
	laptop, err := models.GetLaptopByID(h.DB, id)
	if err != nil {
		log.Printf("Error getting laptop: %v", err)
		http.Error(w, "Laptop not found", http.StatusNotFound)
		return
	}

	// Get reception report if laptop is at warehouse
	var receptionReport *models.ReceptionReport
	if laptop.Status == models.LaptopStatusAtWarehouse {
		receptionReport, err = models.GetLaptopReceptionReport(r.Context(), h.DB, id)
		if err != nil {
			log.Printf("Error getting reception report: %v", err)
			// Don't fail the request, just log the error
		}
	}

	// Prepare template data
	data := map[string]interface{}{
		"User":            user,
		"Laptop":          laptop,
		"ReceptionReport": receptionReport,
	}

	// Execute template using pre-parsed global templates
	if err := h.Templates.ExecuteTemplate(w, "laptop-detail.html", data); err != nil {
		log.Printf("Error executing laptop detail template: %v", err)
		http.Error(w, "Failed to render laptop details", http.StatusInternalServerError)
		return
	}
}

// AddLaptopPage displays the form to add a new laptop
func (h *InventoryHandler) AddLaptopPage(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only logistics role can add laptops
	if user.Role != models.RoleLogistics && user.Role != models.RoleWarehouse {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Get all client companies
	companies, err := models.GetAllClientCompanies(h.DB)
	if err != nil {
		log.Printf("Error getting client companies: %v", err)
		http.Error(w, "Failed to load client companies", http.StatusInternalServerError)
		return
	}

	// Get all software engineers
	engineers, err := models.GetAllSoftwareEngineers(h.DB)
	if err != nil {
		log.Printf("Error getting software engineers: %v", err)
		http.Error(w, "Failed to load software engineers", http.StatusInternalServerError)
		return
	}

	// Prepare template data
	data := map[string]interface{}{
		"User":      user,
		"Statuses":  models.GetLaptopStatusesForNewLaptop(),
		"Companies": companies,
		"Engineers": engineers,
	}

	// Execute template using pre-parsed global templates
	if err := h.Templates.ExecuteTemplate(w, "laptop-form.html", data); err != nil {
		log.Printf("Error executing laptop form template: %v", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}
}

// AddLaptopSubmit handles the submission of a new laptop
func (h *InventoryHandler) AddLaptopSubmit(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only logistics and warehouse roles can add laptops
	if user.Role != models.RoleLogistics && user.Role != models.RoleWarehouse {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Create laptop from form data
	laptop := &models.Laptop{
		SerialNumber: r.FormValue("serial_number"),
		Brand:        r.FormValue("brand"),
		Model:        r.FormValue("model"),
		RAMGB:        r.FormValue("ram_gb"),
		SSDGB:        r.FormValue("ssd_gb"),
		Status:       models.LaptopStatus(r.FormValue("status")),
	}

	// Parse client company ID (required field)
	clientCompanyIDStr := r.FormValue("client_company_id")
	if clientCompanyIDStr != "" {
		clientCompanyID, err := strconv.ParseInt(clientCompanyIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid client company ID", http.StatusBadRequest)
			return
		}
		laptop.ClientCompanyID = &clientCompanyID
	}

	// Parse software engineer ID
	softwareEngineerIDStr := r.FormValue("software_engineer_id")
	if softwareEngineerIDStr != "" {
		softwareEngineerID, err := strconv.ParseInt(softwareEngineerIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid software engineer ID", http.StatusBadRequest)
			return
		}
		laptop.SoftwareEngineerID = &softwareEngineerID
	}

	// Create laptop
	if err := models.CreateLaptop(h.DB, laptop); err != nil {
		log.Printf("Error creating laptop: %v", err)
		http.Error(w, "Failed to create laptop: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to inventory list
	http.Redirect(w, r, "/inventory", http.StatusSeeOther)
}

// EditLaptopPage displays the form to edit an existing laptop
func (h *InventoryHandler) EditLaptopPage(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only logistics and warehouse roles can edit laptops
	if user.Role != models.RoleLogistics && user.Role != models.RoleWarehouse {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Get laptop ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid laptop ID", http.StatusBadRequest)
		return
	}

	// Get laptop
	laptop, err := models.GetLaptopByID(h.DB, id)
	if err != nil {
		log.Printf("Error getting laptop: %v", err)
		http.Error(w, "Laptop not found", http.StatusNotFound)
		return
	}

	// Get all client companies
	companies, err := models.GetAllClientCompanies(h.DB)
	if err != nil {
		log.Printf("Error getting client companies: %v", err)
		http.Error(w, "Failed to load client companies", http.StatusInternalServerError)
		return
	}

	// Get all software engineers
	engineers, err := models.GetAllSoftwareEngineers(h.DB)
	if err != nil {
		log.Printf("Error getting software engineers: %v", err)
		http.Error(w, "Failed to load software engineers", http.StatusInternalServerError)
		return
	}

	// Prepare template data
	data := map[string]interface{}{
		"User":      user,
		"Laptop":    laptop,
		"Statuses":  models.GetLaptopStatusesInOrder(),
		"Companies": companies,
		"Engineers": engineers,
		"IsEdit":    true,
	}

	// Execute template using pre-parsed global templates
	if err := h.Templates.ExecuteTemplate(w, "laptop-form.html", data); err != nil {
		log.Printf("Error executing laptop form template: %v", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}
}

// UpdateLaptopSubmit handles the submission of laptop updates
func (h *InventoryHandler) UpdateLaptopSubmit(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only logistics and warehouse roles can update laptops
	if user.Role != models.RoleLogistics && user.Role != models.RoleWarehouse {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Get laptop ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid laptop ID", http.StatusBadRequest)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get existing laptop
	laptop, err := models.GetLaptopByID(h.DB, id)
	if err != nil {
		log.Printf("Error getting laptop: %v", err)
		http.Error(w, "Laptop not found", http.StatusNotFound)
		return
	}

	// Update laptop fields
	laptop.SerialNumber = r.FormValue("serial_number")
	laptop.Brand = r.FormValue("brand")
	laptop.Model = r.FormValue("model")
	laptop.RAMGB = r.FormValue("ram_gb")
	laptop.SSDGB = r.FormValue("ssd_gb")
	laptop.Status = models.LaptopStatus(r.FormValue("status"))

	// Parse client company ID
	clientCompanyIDStr := r.FormValue("client_company_id")
	if clientCompanyIDStr != "" {
		clientCompanyID, err := strconv.ParseInt(clientCompanyIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid client company ID", http.StatusBadRequest)
			return
		}
		laptop.ClientCompanyID = &clientCompanyID
	} else {
		laptop.ClientCompanyID = nil
	}

	// Parse software engineer ID
	softwareEngineerIDStr := r.FormValue("software_engineer_id")
	if softwareEngineerIDStr != "" {
		softwareEngineerID, err := strconv.ParseInt(softwareEngineerIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid software engineer ID", http.StatusBadRequest)
			return
		}
		laptop.SoftwareEngineerID = &softwareEngineerID
	} else {
		laptop.SoftwareEngineerID = nil
	}

	// Update laptop
	if err := models.UpdateLaptop(h.DB, laptop); err != nil {
		log.Printf("Error updating laptop: %v", err)
		http.Error(w, "Failed to update laptop: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to laptop detail
	http.Redirect(w, r, "/inventory/"+idStr, http.StatusSeeOther)
}

// DeleteLaptop handles the deletion of a laptop
func (h *InventoryHandler) DeleteLaptop(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only logistics role can delete laptops
	if user.Role != models.RoleLogistics {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Get laptop ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid laptop ID", http.StatusBadRequest)
		return
	}

	// Delete laptop
	if err := models.DeleteLaptop(h.DB, id); err != nil {
		log.Printf("Error deleting laptop: %v", err)
		http.Error(w, "Failed to delete laptop: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to inventory list
	http.Redirect(w, r, "/inventory", http.StatusSeeOther)
}

