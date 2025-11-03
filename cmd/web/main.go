package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/yourusername/laptop-tracking-system/internal/config"
	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/email"
	"github.com/yourusername/laptop-tracking-system/internal/handlers"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
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
		"replace": func(old, new, s string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"title": func(s string) string {
			return strings.Title(s)
		},
	}

	templates, err := template.New("").Funcs(funcMap).ParseGlob("templates/pages/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
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
	pickupFormHandler := handlers.NewPickupFormHandler(db, templates, notifier)
	receptionReportHandler := handlers.NewReceptionReportHandler(db, templates, notifier)
	deliveryFormHandler := handlers.NewDeliveryFormHandler(db, templates, notifier)
	shipmentsHandler := handlers.NewShipmentsHandler(db, templates)

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

	// Pickup form routes
	protected.HandleFunc("/pickup-form", pickupFormHandler.PickupFormPage).Methods("GET")
	protected.HandleFunc("/pickup-form", pickupFormHandler.PickupFormSubmit).Methods("POST")

	// Reception report routes
	protected.HandleFunc("/reception-report", receptionReportHandler.ReceptionReportPage).Methods("GET")
	protected.HandleFunc("/reception-report", receptionReportHandler.ReceptionReportSubmit).Methods("POST")

	// Delivery form routes
	protected.HandleFunc("/delivery-form", deliveryFormHandler.DeliveryFormPage).Methods("GET")
	protected.HandleFunc("/delivery-form", deliveryFormHandler.DeliveryFormSubmit).Methods("POST")

	// Shipment routes
	protected.HandleFunc("/shipments", shipmentsHandler.ShipmentsList).Methods("GET")
	protected.HandleFunc("/shipments/{id:[0-9]+}", shipmentsHandler.ShipmentDetail).Methods("GET")
	protected.HandleFunc("/shipments/{id:[0-9]+}/status", shipmentsHandler.UpdateShipmentStatus).Methods("POST")
	protected.HandleFunc("/shipments/{id:[0-9]+}/assign-engineer", shipmentsHandler.AssignEngineer).Methods("POST")

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
