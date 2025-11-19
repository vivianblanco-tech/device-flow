package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
	"github.com/yourusername/laptop-tracking-system/internal/views"
	"golang.org/x/crypto/bcrypt"
)

// FormsHandler handles forms management requests
type FormsHandler struct {
	DB        *sql.DB
	Templates *template.Template
}

// NewFormsHandler creates a new FormsHandler
func NewFormsHandler(db *sql.DB, templates *template.Template) *FormsHandler {
	return &FormsHandler{
		DB:        db,
		Templates: templates,
	}
}

// FormsPage displays the forms management page
func (h *FormsHandler) FormsPage(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only logistics users can access forms page
	if user.Role != models.RoleLogistics {
		http.Error(w, "Forbidden: Only logistics users can access this page", http.StatusForbidden)
		return
	}

	// Prepare template data
	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "forms",
	}

	// Execute template using pre-parsed global templates
	if err := h.Templates.ExecuteTemplate(w, "forms.html", data); err != nil {
		log.Printf("Error executing forms template: %v", err)
		http.Error(w, "Failed to render forms page", http.StatusInternalServerError)
		return
	}
}

// requireLogisticsRole checks if the user is a logistics user
func (h *FormsHandler) requireLogisticsRole(w http.ResponseWriter, r *http.Request) bool {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return false
	}
	if user.Role != models.RoleLogistics {
		http.Error(w, "Forbidden: Only logistics users can access this page", http.StatusForbidden)
		return false
	}
	return true
}

// ========== USER HANDLERS ==========

// UsersList displays a list of all users
func (h *FormsHandler) UsersList(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	users, err := models.GetAllUsers(h.DB)
	if err != nil {
		log.Printf("Error getting users: %v", err)
		http.Error(w, "Failed to load users", http.StatusInternalServerError)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "forms",
		"Users":       users,
		"Roles":       []models.UserRole{models.RoleLogistics, models.RoleClient, models.RoleWarehouse, models.RoleProjectManager},
	}

	if err := h.Templates.ExecuteTemplate(w, "users-list.html", data); err != nil {
		log.Printf("Error executing users list template: %v", err)
		http.Error(w, "Failed to render users list", http.StatusInternalServerError)
		return
	}
}

// UserAddPage displays the form to add a new user
func (h *FormsHandler) UserAddPage(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	companies, err := models.GetAllClientCompanies(h.DB)
	if err != nil {
		log.Printf("Error getting client companies: %v", err)
		http.Error(w, "Failed to load client companies", http.StatusInternalServerError)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "forms",
		"Companies":   companies,
		"Roles":       []models.UserRole{models.RoleLogistics, models.RoleClient, models.RoleWarehouse, models.RoleProjectManager},
	}

	if err := h.Templates.ExecuteTemplate(w, "user-form.html", data); err != nil {
		log.Printf("Error executing user form template: %v", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}
}

// UserAddSubmit handles the submission of a new user
func (h *FormsHandler) UserAddSubmit(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Hash password - required for new users
	password := r.FormValue("password")
	if password == "" {
		http.Error(w, "Password is required for new users", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Email:        r.FormValue("email"),
		PasswordHash: string(hashedPassword),
		Role:         models.UserRole(r.FormValue("role")),
	}

	// Parse client company ID if provided
	if clientCompanyIDStr := r.FormValue("client_company_id"); clientCompanyIDStr != "" {
		clientCompanyID, err := strconv.ParseInt(clientCompanyIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid client company ID", http.StatusBadRequest)
			return
		}
		user.ClientCompanyID = &clientCompanyID
	}

	if err := models.CreateUser(h.DB, user); err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/forms/users?success="+url.QueryEscape("User created successfully"), http.StatusSeeOther)
}

// UserEditPage displays the form to edit an existing user
func (h *FormsHandler) UserEditPage(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByID(h.DB, id)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	companies, err := models.GetAllClientCompanies(h.DB)
	if err != nil {
		log.Printf("Error getting client companies: %v", err)
		http.Error(w, "Failed to load client companies", http.StatusInternalServerError)
		return
	}

	currentUser := middleware.GetUserFromContext(r.Context())
	
	// Get current company ID for pre-selection (0 if nil)
	var currentCompanyID int64
	if user.ClientCompanyID != nil {
		currentCompanyID = *user.ClientCompanyID
	}
	
	data := map[string]interface{}{
		"User":             currentUser,
		"Nav":              views.GetNavigationLinks(currentUser.Role),
		"CurrentPage":      "forms",
		"EditUser":         user,
		"Companies":        companies,
		"Roles":            []models.UserRole{models.RoleLogistics, models.RoleClient, models.RoleWarehouse, models.RoleProjectManager},
		"IsEdit":           true,
		"CurrentCompanyID": currentCompanyID,
	}

	if err := h.Templates.ExecuteTemplate(w, "user-form.html", data); err != nil {
		log.Printf("Error executing user form template: %v", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}
}

// UserEditSubmit handles the submission of user updates
func (h *FormsHandler) UserEditSubmit(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByID(h.DB, id)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Update fields
	user.Email = r.FormValue("email")
	user.Role = models.UserRole(r.FormValue("role"))

	// Update password if provided
	if password := r.FormValue("password"); password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			http.Error(w, "Failed to process password", http.StatusInternalServerError)
			return
		}
		user.PasswordHash = string(hashedPassword)
	}
	// If no password provided, existing password hash from GetUserByID is preserved

	// Parse client company ID
	if clientCompanyIDStr := r.FormValue("client_company_id"); clientCompanyIDStr != "" {
		clientCompanyID, err := strconv.ParseInt(clientCompanyIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid client company ID", http.StatusBadRequest)
			return
		}
		user.ClientCompanyID = &clientCompanyID
	} else {
		user.ClientCompanyID = nil
	}

	if err := models.UpdateUser(h.DB, user); err != nil {
		log.Printf("Error updating user: %v", err)
		http.Redirect(w, r, "/forms/users/"+idStr+"/edit?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/forms/users?success="+url.QueryEscape("User updated successfully"), http.StatusSeeOther)
}

// ========== CLIENT COMPANY HANDLERS ==========

// ClientCompaniesList displays a list of all client companies
func (h *FormsHandler) ClientCompaniesList(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	companies, err := models.GetAllClientCompanies(h.DB)
	if err != nil {
		log.Printf("Error getting client companies: %v", err)
		http.Error(w, "Failed to load client companies", http.StatusInternalServerError)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "forms",
		"Companies":   companies,
	}

	if err := h.Templates.ExecuteTemplate(w, "client-companies-list.html", data); err != nil {
		log.Printf("Error executing client companies list template: %v", err)
		http.Error(w, "Failed to render client companies list", http.StatusInternalServerError)
		return
	}
}

// ClientCompanyAddPage displays the form to add a new client company
func (h *FormsHandler) ClientCompanyAddPage(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "forms",
	}

	if err := h.Templates.ExecuteTemplate(w, "client-company-form.html", data); err != nil {
		log.Printf("Error executing client company form template: %v", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}
}

// ClientCompanyAddSubmit handles the submission of a new client company
func (h *FormsHandler) ClientCompanyAddSubmit(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	company := &models.ClientCompany{
		Name:        r.FormValue("name"),
		ContactInfo: r.FormValue("contact_info"),
	}

	if err := models.CreateClientCompany(h.DB, company); err != nil {
		log.Printf("Error creating client company: %v", err)
		http.Error(w, "Failed to create client company: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/forms/client-companies?success="+url.QueryEscape("Client company created successfully"), http.StatusSeeOther)
}

// ClientCompanyEditPage displays the form to edit an existing client company
func (h *FormsHandler) ClientCompanyEditPage(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid client company ID", http.StatusBadRequest)
		return
	}

	company, err := models.GetClientCompanyByID(h.DB, id)
	if err != nil {
		log.Printf("Error getting client company: %v", err)
		http.Error(w, "Client company not found", http.StatusNotFound)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "forms",
		"Company":     company,
		"IsEdit":      true,
	}

	if err := h.Templates.ExecuteTemplate(w, "client-company-form.html", data); err != nil {
		log.Printf("Error executing client company form template: %v", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}
}

// ClientCompanyEditSubmit handles the submission of client company updates
func (h *FormsHandler) ClientCompanyEditSubmit(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid client company ID", http.StatusBadRequest)
		return
	}

	company, err := models.GetClientCompanyByID(h.DB, id)
	if err != nil {
		log.Printf("Error getting client company: %v", err)
		http.Error(w, "Client company not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	company.Name = r.FormValue("name")
	company.ContactInfo = r.FormValue("contact_info")

	if err := models.UpdateClientCompany(h.DB, company); err != nil {
		log.Printf("Error updating client company: %v", err)
		http.Redirect(w, r, "/forms/client-companies/"+idStr+"/edit?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/forms/client-companies?success="+url.QueryEscape("Client company updated successfully"), http.StatusSeeOther)
}

// ========== SOFTWARE ENGINEER HANDLERS ==========

// SoftwareEngineersList displays a list of all software engineers
func (h *FormsHandler) SoftwareEngineersList(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	// Parse query parameters
	searchQuery := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort")
	sortOrder := r.URL.Query().Get("order")

	// Build filter
	filter := &models.SoftwareEngineerFilter{
		Search:    searchQuery,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	engineers, err := models.GetAllSoftwareEngineers(h.DB, filter)
	if err != nil {
		log.Printf("Error getting software engineers: %v", err)
		http.Error(w, "Failed to load software engineers", http.StatusInternalServerError)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "forms",
		"Engineers":   engineers,
		"SearchQuery": searchQuery,
		"SortBy":      sortBy,
		"SortOrder":   sortOrder,
	}

	if err := h.Templates.ExecuteTemplate(w, "software-engineers-list.html", data); err != nil {
		log.Printf("Error executing software engineers list template: %v", err)
		http.Error(w, "Failed to render software engineers list", http.StatusInternalServerError)
		return
	}
}

// SoftwareEngineerAddPage displays the form to add a new software engineer
func (h *FormsHandler) SoftwareEngineerAddPage(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "forms",
	}

	if err := h.Templates.ExecuteTemplate(w, "software-engineer-form.html", data); err != nil {
		log.Printf("Error executing software engineer form template: %v", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}
}

// SoftwareEngineerAddSubmit handles the submission of a new software engineer
func (h *FormsHandler) SoftwareEngineerAddSubmit(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	addressConfirmed := r.FormValue("address_confirmed") == "on"
	engineer := &models.SoftwareEngineer{
		Name:           r.FormValue("name"),
		Email:          r.FormValue("email"),
		Address:        r.FormValue("address"),
		Phone:          r.FormValue("phone"),
		EmployeeNumber: r.FormValue("employee_number"),
		AddressConfirmed: addressConfirmed,
	}

	// Set address confirmation timestamp when checkbox is checked
	if addressConfirmed {
		engineer.ConfirmAddress()
	}

	if err := models.CreateSoftwareEngineer(h.DB, engineer); err != nil {
		log.Printf("Error creating software engineer: %v", err)
		http.Error(w, "Failed to create software engineer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/forms/software-engineers?success="+url.QueryEscape("Software engineer created successfully"), http.StatusSeeOther)
}

// SoftwareEngineerEditPage displays the form to edit an existing software engineer
func (h *FormsHandler) SoftwareEngineerEditPage(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid software engineer ID", http.StatusBadRequest)
		return
	}

	engineer, err := models.GetSoftwareEngineerByID(h.DB, id)
	if err != nil {
		log.Printf("Error getting software engineer: %v", err)
		http.Error(w, "Software engineer not found", http.StatusNotFound)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "forms",
		"Engineer":    engineer,
		"IsEdit":      true,
	}

	if err := h.Templates.ExecuteTemplate(w, "software-engineer-form.html", data); err != nil {
		log.Printf("Error executing software engineer form template: %v", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}
}

// SoftwareEngineerEditSubmit handles the submission of software engineer updates
func (h *FormsHandler) SoftwareEngineerEditSubmit(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid software engineer ID", http.StatusBadRequest)
		return
	}

	engineer, err := models.GetSoftwareEngineerByID(h.DB, id)
	if err != nil {
		log.Printf("Error getting software engineer: %v", err)
		http.Error(w, "Software engineer not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	engineer.Name = r.FormValue("name")
	engineer.Email = r.FormValue("email")
	engineer.Address = r.FormValue("address")
	engineer.Phone = r.FormValue("phone")
	engineer.EmployeeNumber = r.FormValue("employee_number")
	
	addressConfirmed := r.FormValue("address_confirmed") == "on"
	wasConfirmed := engineer.AddressConfirmed
	engineer.AddressConfirmed = addressConfirmed

	// Handle address confirmation timestamp
	if addressConfirmed {
		// If changing from unchecked to checked, set timestamp
		if !wasConfirmed {
			engineer.ConfirmAddress()
		}
		// If already checked, preserve existing timestamp (do nothing)
	} else {
		// If unchecked, clear timestamp
		engineer.AddressConfirmationAt = nil
	}

	if err := models.UpdateSoftwareEngineer(h.DB, engineer); err != nil {
		log.Printf("Error updating software engineer: %v", err)
		http.Redirect(w, r, "/forms/software-engineers/"+idStr+"/edit?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/forms/software-engineers?success="+url.QueryEscape("Software engineer updated successfully"), http.StatusSeeOther)
}

// ========== COURIER HANDLERS ==========

// CouriersList displays a list of all couriers
func (h *FormsHandler) CouriersList(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	couriers, err := models.GetAllCouriers(h.DB)
	if err != nil {
		log.Printf("Error getting couriers: %v", err)
		http.Error(w, "Failed to load couriers", http.StatusInternalServerError)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "forms",
		"Couriers":    couriers,
	}

	if err := h.Templates.ExecuteTemplate(w, "couriers-list.html", data); err != nil {
		log.Printf("Error executing couriers list template: %v", err)
		http.Error(w, "Failed to render couriers list", http.StatusInternalServerError)
		return
	}
}

// CourierAddPage displays the form to add a new courier
func (h *FormsHandler) CourierAddPage(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "forms",
	}

	if err := h.Templates.ExecuteTemplate(w, "courier-form.html", data); err != nil {
		log.Printf("Error executing courier form template: %v", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}
}

// CourierAddSubmit handles the submission of a new courier
func (h *FormsHandler) CourierAddSubmit(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	courier := &models.Courier{
		Name:        r.FormValue("name"),
		ContactInfo: r.FormValue("contact_info"),
	}

	if err := models.CreateCourier(h.DB, courier); err != nil {
		log.Printf("Error creating courier: %v", err)
		http.Error(w, "Failed to create courier: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/forms/couriers?success="+url.QueryEscape("Courier created successfully"), http.StatusSeeOther)
}

// CourierEditPage displays the form to edit an existing courier
func (h *FormsHandler) CourierEditPage(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid courier ID", http.StatusBadRequest)
		return
	}

	courier, err := models.GetCourierByID(h.DB, id)
	if err != nil {
		log.Printf("Error getting courier: %v", err)
		http.Error(w, "Courier not found", http.StatusNotFound)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "forms",
		"Courier":     courier,
		"IsEdit":      true,
	}

	if err := h.Templates.ExecuteTemplate(w, "courier-form.html", data); err != nil {
		log.Printf("Error executing courier form template: %v", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}
}

// CourierEditSubmit handles the submission of courier updates
func (h *FormsHandler) CourierEditSubmit(w http.ResponseWriter, r *http.Request) {
	if !h.requireLogisticsRole(w, r) {
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid courier ID", http.StatusBadRequest)
		return
	}

	courier, err := models.GetCourierByID(h.DB, id)
	if err != nil {
		log.Printf("Error getting courier: %v", err)
		http.Error(w, "Courier not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	courier.Name = r.FormValue("name")
	courier.ContactInfo = r.FormValue("contact_info")

	if err := models.UpdateCourier(h.DB, courier); err != nil {
		log.Printf("Error updating courier: %v", err)
		http.Redirect(w, r, "/forms/couriers/"+idStr+"/edit?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/forms/couriers?success="+url.QueryEscape("Courier updated successfully"), http.StatusSeeOther)
}

