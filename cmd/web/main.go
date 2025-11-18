package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/yourusername/laptop-tracking-system/internal/config"
	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/email"
	"github.com/yourusername/laptop-tracking-system/internal/handlers"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
	"github.com/yourusername/laptop-tracking-system/internal/views"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connected successfully")

	// Load templates with custom functions
	funcMap := template.FuncMap{
		"replace": func(old, new string, v interface{}) string {
			// Convert interface{} to string first
			var s string
			switch val := v.(type) {
			case string:
				s = val
			case models.UserRole:
				s = string(val)
			case models.LaptopStatus:
				s = string(val)
			default:
				s = fmt.Sprintf("%v", val)
			}
			return strings.ReplaceAll(s, old, new)
		},
		"title": func(v interface{}) string {
			// Convert interface{} to string
			var s string
			switch val := v.(type) {
			case string:
				s = val
			case models.UserRole:
				s = string(val)
			case models.LaptopStatus:
				s = string(val)
			default:
				s = fmt.Sprintf("%v", val)
			}
			return strings.Title(s)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"len": func(v interface{}) int {
			switch val := v.(type) {
			case []models.TimelineItem:
				return len(val)
			case []interface{}:
				return len(val)
			default:
				return 0
			}
		},
		// Navigation helper function
		"getNav": func(role models.UserRole) views.NavigationLinks {
			return views.GetNavigationLinks(role)
		},
		// Calendar template functions
		"formatDate": func(t time.Time) string {
			return t.Format("Jan 2, 2006")
		},
		"formatTime": func(t time.Time) string {
			return t.Format("3:04 PM")
		},
		"formatDateShort": func(t time.Time) string {
			return t.Format("Jan 2")
		},
		"daysInMonth": func(year int, month time.Month) int {
			return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
		},
		"firstWeekday": func(year int, month time.Month) time.Weekday {
			return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).Weekday()
		},
		// Dashboard template functions
		"statusColor": func(status models.ShipmentStatus) string {
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
		},
		"laptopStatusColor": func(status models.LaptopStatus) string {
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
		},
		// Inventory template specific statusColor (with text color)
		"inventoryStatusColor": func(status models.LaptopStatus) string {
			switch status {
			case models.LaptopStatusAvailable:
				return "bg-green-100 text-green-800"
			case models.LaptopStatusInTransitToWarehouse:
				return "bg-purple-100 text-purple-800"
			case models.LaptopStatusAtWarehouse:
				return "bg-indigo-100 text-indigo-800"
			case models.LaptopStatusInTransitToEngineer:
				return "bg-cyan-100 text-cyan-800"
			case models.LaptopStatusDelivered:
				return "bg-blue-100 text-blue-800"
			case models.LaptopStatusRetired:
				return "bg-gray-100 text-gray-800"
			default:
				return "bg-gray-100 text-gray-800"
			}
		},
		"laptopStatusDisplayName": func(status models.LaptopStatus) string {
			return models.GetLaptopStatusDisplayName(status)
		},
		"receptionReportStatusColor": func(status string) string {
			switch models.ReceptionReportStatus(status) {
			case models.ReceptionReportStatusPendingApproval:
				return "bg-yellow-100 text-yellow-800"
			case models.ReceptionReportStatusApproved:
				return "bg-green-100 text-green-800"
			default:
				return "bg-gray-100 text-gray-800"
			}
		},
		"receptionReportStatusDisplayName": func(status string) string {
			switch models.ReceptionReportStatus(status) {
			case models.ReceptionReportStatusPendingApproval:
				return "Pending Approval"
			case models.ReceptionReportStatusApproved:
				return "Approved"
			default:
				return "Unknown"
			}
		},
	}

	templates, err := template.New("").Funcs(funcMap).ParseGlob("templates/pages/*.html")
	if err != nil {
		log.Fatalf("Failed to parse page templates: %v", err)
	}

	// Parse component templates (navbar, etc.)
	templates, err = templates.ParseGlob("templates/components/*.html")
	if err != nil {
		log.Fatalf("Failed to parse component templates: %v", err)
	}

	// Set up Google OAuth config
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.Google.ClientID,
		ClientSecret: cfg.Google.ClientSecret,
		RedirectURL:  cfg.Google.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Initialize email client and notifier
	smtpPort, err := strconv.Atoi(cfg.SMTP.Port)
	if err != nil {
		log.Printf("Warning: Invalid SMTP port '%s', defaulting to 1025", cfg.SMTP.Port)
		smtpPort = 1025
	}

	emailClient, err := email.NewClient(email.Config{
		Host:     cfg.SMTP.Host,
		Port:     smtpPort,
		Username: cfg.SMTP.User,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
	})
	if err != nil {
		log.Printf("Warning: Failed to initialize email client: %v", err)
		log.Println("Email notifications will be disabled")
		emailClient = nil
	}

	var notifier *email.Notifier
	if emailClient != nil {
		notifier = email.NewNotifier(emailClient, db)
		log.Println("Email notifications enabled")
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, templates)
	authHandler.OAuthConfig = oauthConfig
	authHandler.OAuthDomain = cfg.Google.AllowedDomain

	dashboardHandler := handlers.NewDashboardHandler(db, templates)
	chartsHandler := handlers.NewChartsHandler(db)
	calendarHandler := handlers.NewCalendarHandler(db, templates)
	inventoryHandler := handlers.NewInventoryHandler(db, templates)
	pickupFormHandler := handlers.NewPickupFormHandler(db, templates, notifier)
	receptionReportHandler := handlers.NewReceptionReportHandler(db, templates, notifier)
	deliveryFormHandler := handlers.NewDeliveryFormHandler(db, templates, notifier)
	shipmentsHandler := handlers.NewShipmentsHandler(db, templates, notifier)

	// Initialize router
	router := mux.NewRouter()

	// Apply auth middleware globally
	router.Use(middleware.AuthMiddleware(db))

	// Public routes
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		user := middleware.GetUserFromContext(r.Context())
		if user != nil {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}).Methods("GET")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Authentication routes (public)
	router.HandleFunc("/login", authHandler.LoginPage).Methods("GET")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/logout", authHandler.Logout).Methods("POST", "GET")
	router.HandleFunc("/auth/google", authHandler.GoogleLogin).Methods("GET")
	router.HandleFunc("/auth/google/callback", authHandler.GoogleCallback).Methods("GET")
	router.HandleFunc("/auth/magic-link", authHandler.MagicLinkLogin).Methods("GET")

	// Protected routes (require authentication)
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.RequireAuth)

	// Dashboard
	protected.HandleFunc("/dashboard", dashboardHandler.Dashboard).Methods("GET")

	// Calendar
	protected.HandleFunc("/calendar", calendarHandler.Calendar).Methods("GET")

	// Magic Links (logistics only)
	protected.HandleFunc("/magic-links", authHandler.MagicLinksList).Methods("GET")
	protected.HandleFunc("/auth/send-magic-link", authHandler.SendMagicLink).Methods("POST")

	// Inventory routes
	protected.HandleFunc("/inventory", inventoryHandler.InventoryList).Methods("GET")
	protected.HandleFunc("/inventory/add", inventoryHandler.AddLaptopPage).Methods("GET")
	protected.HandleFunc("/inventory/add", inventoryHandler.AddLaptopSubmit).Methods("POST")
	protected.HandleFunc("/inventory/{id:[0-9]+}", inventoryHandler.LaptopDetail).Methods("GET")
	protected.HandleFunc("/inventory/{id:[0-9]+}/edit", inventoryHandler.EditLaptopPage).Methods("GET")
	protected.HandleFunc("/inventory/{id:[0-9]+}/update", inventoryHandler.UpdateLaptopSubmit).Methods("POST")
	protected.HandleFunc("/inventory/{id:[0-9]+}/delete", inventoryHandler.DeleteLaptop).Methods("POST")
	
	// Laptop reception report routes (laptop-based)
	protected.HandleFunc("/laptops/{id:[0-9]+}/reception-report", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		laptopID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid laptop ID", http.StatusBadRequest)
			return
		}
		if r.Method == "GET" {
			receptionReportHandler.LaptopReceptionReportPage(w, r, laptopID)
		} else {
			receptionReportHandler.LaptopReceptionReportSubmit(w, r, laptopID)
		}
	}).Methods("GET", "POST")

	// Chart API endpoints
	protected.HandleFunc("/api/charts/shipments-over-time", chartsHandler.ShipmentsOverTimeAPI).Methods("GET")
	protected.HandleFunc("/api/charts/status-distribution", chartsHandler.StatusDistributionAPI).Methods("GET")
	protected.HandleFunc("/api/charts/delivery-time-trends", chartsHandler.DeliveryTimeTrendsAPI).Methods("GET")

	// Pickup forms landing page
	protected.HandleFunc("/pickup-forms", pickupFormHandler.PickupFormsLandingPage).Methods("GET")
	
	// Pickup form routes (legacy)
	protected.HandleFunc("/pickup-form", pickupFormHandler.PickupFormPage).Methods("GET")
	protected.HandleFunc("/pickup-form", pickupFormHandler.PickupFormSubmit).Methods("POST")
	
	// Three shipment type form routes (Phase 5)
	protected.HandleFunc("/shipments/create/single", pickupFormHandler.SingleShipmentFormPage).Methods("GET")
	protected.HandleFunc("/shipments/create/single-minimal", pickupFormHandler.CreateMinimalSingleShipment).Methods("POST")
	protected.HandleFunc("/shipments/create/bulk", pickupFormHandler.BulkShipmentFormPage).Methods("GET")
	protected.HandleFunc("/shipments/create/warehouse-to-engineer", pickupFormHandler.WarehouseToEngineerFormPage).Methods("GET")

	// Reception report routes (laptop-based)
	protected.HandleFunc("/reception-reports", receptionReportHandler.LaptopBasedReceptionReportsList).Methods("GET")
	protected.HandleFunc("/reception-reports/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid report ID", http.StatusBadRequest)
			return
		}
		receptionReportHandler.ReceptionReportDetail(w, r, id)
	}).Methods("GET")
	protected.HandleFunc("/reception-reports/{id:[0-9]+}/approve", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		reportID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid report ID", http.StatusBadRequest)
			return
		}
		receptionReportHandler.ApproveReceptionReport(w, r, reportID)
	}).Methods("POST")
	
	// Legacy reception report routes (shipment-based - deprecated)
	protected.HandleFunc("/reception-report", receptionReportHandler.ReceptionReportPage).Methods("GET")
	protected.HandleFunc("/reception-report", receptionReportHandler.ReceptionReportSubmit).Methods("POST")

	// Delivery form routes
	protected.HandleFunc("/delivery-form", deliveryFormHandler.DeliveryFormPage).Methods("GET")
	protected.HandleFunc("/delivery-form", deliveryFormHandler.DeliveryFormSubmit).Methods("POST")

	// Shipment routes
	protected.HandleFunc("/shipments", shipmentsHandler.ShipmentsList).Methods("GET")
	protected.HandleFunc("/shipments/create", shipmentsHandler.CreateShipment).Methods("GET", "POST")
	protected.HandleFunc("/shipments/{id:[0-9]+}", shipmentsHandler.ShipmentDetail).Methods("GET")
	protected.HandleFunc("/shipments/{id:[0-9]+}/status", shipmentsHandler.UpdateShipmentStatus).Methods("POST")
	protected.HandleFunc("/shipments/{id:[0-9]+}/assign-engineer", shipmentsHandler.AssignEngineer).Methods("POST")
	protected.HandleFunc("/shipments/{id:[0-9]+}/edit", shipmentsHandler.EditShipmentGET).Methods("GET")
	protected.HandleFunc("/shipments/{id:[0-9]+}/edit", shipmentsHandler.EditShipmentPOST).Methods("POST")
	protected.HandleFunc("/shipments/{id:[0-9]+}/form", shipmentsHandler.ShipmentPickupFormPage).Methods("GET")
	protected.HandleFunc("/shipments/{id:[0-9]+}/form", shipmentsHandler.ShipmentPickupFormSubmit).Methods("POST")
	protected.HandleFunc("/shipments/{id:[0-9]+}/complete-details", pickupFormHandler.CompleteShipmentDetails).Methods("POST")
	protected.HandleFunc("/shipments/{id:[0-9]+}/edit-details", pickupFormHandler.EditShipmentDetails).Methods("POST")
	protected.HandleFunc("/shipments/{id:[0-9]+}/laptops/add", shipmentsHandler.AddLaptopToBulkShipment).Methods("POST")

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static"))))

	// Serve uploads (photos)
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/",
		http.FileServer(http.Dir("./uploads"))))

	// Start server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	log.Printf("Environment: %s", cfg.App.Environment)
	log.Printf("Login at: http://%s/login", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
