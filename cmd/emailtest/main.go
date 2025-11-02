package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/yourusername/laptop-tracking-system/internal/email"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	// Configure email client for Mailhog
	emailClient, err := email.NewClient(email.Config{
		Host:     getEnv("SMTP_HOST", "localhost"),
		Port:     getEnvInt("SMTP_PORT", 1025),
		Username: getEnv("SMTP_USERNAME", ""),
		Password: getEnv("SMTP_PASSWORD", ""),
		From:     getEnv("SMTP_FROM", "noreply@bairesdev.com"),
	})
	if err != nil {
		log.Fatalf("❌ Failed to create email client: %v", err)
	}

	fmt.Println("📧 Email Testing Tool")
	fmt.Println("=====================")
	fmt.Println()
	fmt.Printf("SMTP Host: %s:%d\n", getEnv("SMTP_HOST", "localhost"), getEnvInt("SMTP_PORT", 1025))
	fmt.Printf("From: %s\n", getEnv("SMTP_FROM", "noreply@bairesdev.com"))
	fmt.Println()

	// Initialize templates
	templates := email.NewEmailTemplates()

	// Test recipient
	recipient := getEnv("TEST_EMAIL", "test@example.com")

	fmt.Println("Sending test emails to:", recipient)
	fmt.Println()

	// 1. Test Magic Link Email
	fmt.Print("1️⃣  Sending Magic Link email... ")
	if err := sendMagicLinkTest(emailClient, templates, recipient); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Sent")
	}

	time.Sleep(500 * time.Millisecond)

	// 2. Test Address Confirmation Email
	fmt.Print("2️⃣  Sending Address Confirmation email... ")
	if err := sendAddressConfirmationTest(emailClient, templates, recipient); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Sent")
	}

	time.Sleep(500 * time.Millisecond)

	// 3. Test Pickup Confirmation Email
	fmt.Print("3️⃣  Sending Pickup Confirmation email... ")
	if err := sendPickupConfirmationTest(emailClient, templates, recipient); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Sent")
	}

	time.Sleep(500 * time.Millisecond)

	// 4. Test Warehouse Pre-Alert Email
	fmt.Print("4️⃣  Sending Warehouse Pre-Alert email... ")
	if err := sendWarehousePreAlertTest(emailClient, templates, recipient); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Sent")
	}

	time.Sleep(500 * time.Millisecond)

	// 5. Test Release Notification Email
	fmt.Print("5️⃣  Sending Release Notification email... ")
	if err := sendReleaseNotificationTest(emailClient, templates, recipient); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Sent")
	}

	time.Sleep(500 * time.Millisecond)

	// 6. Test Delivery Confirmation Email
	fmt.Print("6️⃣  Sending Delivery Confirmation email... ")
	if err := sendDeliveryConfirmationTest(emailClient, templates, recipient); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Println("✅ Sent")
	}

	fmt.Println()
	fmt.Println("🎉 All test emails sent successfully!")
	fmt.Println()
	fmt.Println("📬 Check Mailhog web UI at: http://localhost:8025")
	fmt.Println()
}

func sendMagicLinkTest(client *email.Client, templates *email.EmailTemplates, recipient string) error {
	data := email.MagicLinkData{
		RecipientName: "John Doe",
		MagicLink:     "http://localhost:8080/forms/pickup?token=abc123def456ghi789",
		ExpiresAt:     time.Now().Add(24 * time.Hour),
		FormType:      "Pickup Form",
	}

	html, err := templates.RenderTemplate("magic_link", data)
	if err != nil {
		return err
	}

	return client.Send(email.Message{
		To:       []string{recipient},
		Subject:  templates.GetSubject("magic_link", data),
		Body:     "Please view this email in an HTML-capable email client.",
		HTMLBody: html,
	})
}

func sendAddressConfirmationTest(client *email.Client, templates *email.EmailTemplates, recipient string) error {
	data := email.AddressConfirmationData{
		EngineerName:    "Sarah Martinez",
		CompanyName:     "BairesDev",
		ProjectName:     "E-Commerce Platform Redesign",
		ExpectedDate:    "December 15, 2024",
		ConfirmationURL: "http://localhost:8080/confirm-address?id=xyz789",
	}

	html, err := templates.RenderTemplate("address_confirmation", data)
	if err != nil {
		return err
	}

	return client.Send(email.Message{
		To:       []string{recipient},
		Subject:  templates.GetSubject("address_confirmation", data),
		Body:     "Please view this email in an HTML-capable email client.",
		HTMLBody: html,
	})
}

func sendPickupConfirmationTest(client *email.Client, templates *email.EmailTemplates, recipient string) error {
	data := email.PickupConfirmationData{
		ClientName:       "Alice Johnson",
		ClientCompany:    "TechCorp Solutions",
		TrackingNumber:   "1Z999AA10123456784",
		PickupDate:       "Monday, December 9, 2024",
		PickupTimeSlot:   "Morning (9:00 AM - 12:00 PM)",
		NumberOfDevices:  3,
		ConfirmationCode: "CONF-2024-12345",
	}

	html, err := templates.RenderTemplate("pickup_confirmation", data)
	if err != nil {
		return err
	}

	return client.Send(email.Message{
		To:       []string{recipient},
		Subject:  templates.GetSubject("pickup_confirmation", data),
		Body:     "Please view this email in an HTML-capable email client.",
		HTMLBody: html,
	})
}

func sendWarehousePreAlertTest(client *email.Client, templates *email.EmailTemplates, recipient string) error {
	data := email.WarehousePreAlertData{
		TrackingNumber:    "1Z999AA10123456784",
		ExpectedDate:      "Wednesday, December 11, 2024",
		ShipperName:       "Alice Johnson",
		ShipperCompany:    "TechCorp Solutions",
		DeviceDescription: "3 configured Dell Latitude laptops",
		ProjectName:       "E-Commerce Platform Redesign",
		TrackingURL:       "https://www.ups.com/track?tracknum=1Z999AA10123456784",
	}

	html, err := templates.RenderTemplate("warehouse_pre_alert", data)
	if err != nil {
		return err
	}

	return client.Send(email.Message{
		To:       []string{recipient},
		Subject:  templates.GetSubject("warehouse_pre_alert", data),
		Body:     "Please view this email in an HTML-capable email client.",
		HTMLBody: html,
	})
}

func sendReleaseNotificationTest(client *email.Client, templates *email.EmailTemplates, recipient string) error {
	data := email.ReleaseNotificationData{
		CourierName:        "FedEx Express",
		CourierCompany:     "FedEx",
		PickupDate:         "Friday, December 13, 2024",
		PickupTimeSlot:     "Afternoon (1:00 PM - 5:00 PM)",
		WarehouseAddress:   "456 Warehouse Drive, San Francisco, CA 94102",
		ContactPerson:      "Mike Anderson",
		ContactPhone:       "+1 (415) 555-0199",
		DeviceSerialNumber: "DELL-LAT-2024-001234",
		EngineerName:       "Sarah Martinez",
		TrackingNumber:     "1Z999AA10123456784",
	}

	html, err := templates.RenderTemplate("release_notification", data)
	if err != nil {
		return err
	}

	return client.Send(email.Message{
		To:       []string{recipient},
		Subject:  templates.GetSubject("release_notification", data),
		Body:     "Please view this email in an HTML-capable email client.",
		HTMLBody: html,
	})
}

func sendDeliveryConfirmationTest(client *email.Client, templates *email.EmailTemplates, recipient string) error {
	data := email.DeliveryConfirmationData{
		EngineerName:       "Sarah Martinez",
		DeviceSerialNumber: "DELL-LAT-2024-001234",
		DeviceModel:        "Dell Latitude 7430 (Intel i7, 16GB RAM, 512GB SSD)",
		DeliveryDate:       "Monday, December 16, 2024",
		TrackingNumber:     "1Z999AA10123456784",
		ProjectName:        "E-Commerce Platform Redesign",
	}

	html, err := templates.RenderTemplate("delivery_confirmation", data)
	if err != nil {
		return err
	}

	return client.Send(email.Message{
		To:       []string{recipient},
		Subject:  templates.GetSubject("delivery_confirmation", data),
		Body:     "Please view this email in an HTML-capable email client.",
		HTMLBody: html,
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}


