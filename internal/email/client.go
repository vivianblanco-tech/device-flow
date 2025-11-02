package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

// Config holds the SMTP configuration for the email client
type Config struct {
	Host     string // SMTP server hostname (e.g., "smtp.gmail.com")
	Port     int    // SMTP server port (e.g., 587 for TLS, 465 for SSL)
	Username string // SMTP authentication username (optional for some servers)
	Password string // SMTP authentication password (optional for some servers)
	From     string // Default sender email address
}

// Message represents an email message to be sent
type Message struct {
	To       []string // List of recipient email addresses
	Subject  string   // Email subject line
	Body     string   // Plain text body
	HTMLBody string   // HTML body (optional)
}

// Client represents an email client for sending messages via SMTP
type Client struct {
	config Config
}

// NewClient creates a new email client with the given configuration
// Returns an error if the configuration is invalid
func NewClient(config Config) (*Client, error) {
	// Validate required configuration
	if config.Host == "" {
		return nil, fmt.Errorf("SMTP host is required")
	}

	if config.Port < 1 || config.Port > 65535 {
		return nil, fmt.Errorf("SMTP port must be between 1 and 65535")
	}

	if config.From == "" {
		return nil, fmt.Errorf("from address is required")
	}

	return &Client{
		config: config,
	}, nil
}

// BuildMessage validates a message structure
// Returns an error if the message is invalid
func (c *Client) BuildMessage(msg Message) error {
	// Validate recipients
	if len(msg.To) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}

	// Validate subject
	if msg.Subject == "" {
		return fmt.Errorf("subject is required")
	}

	// Validate body
	if msg.Body == "" {
		return fmt.Errorf("body is required")
	}

	return nil
}

// Send sends an email message via SMTP
// Returns an error if the message fails to send
func (c *Client) Send(msg Message) error {
	// Validate the message
	if err := c.BuildMessage(msg); err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}

	// Build the email content
	body := c.buildEmailBody(msg)

	// Set up authentication if credentials are provided
	var auth smtp.Auth
	if c.config.Username != "" && c.config.Password != "" {
		auth = smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)
	}

	// Build server address
	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)

	// For port 587 (TLS), use STARTTLS
	if c.config.Port == 587 {
		return c.sendWithTLS(addr, auth, msg.To, body)
	}

	// For other ports (25, 465, etc.), use standard SMTP
	return smtp.SendMail(addr, auth, c.config.From, msg.To, body)
}

// sendWithTLS sends email using STARTTLS (for port 587)
func (c *Client) sendWithTLS(addr string, auth smtp.Auth, to []string, body []byte) error {
	// Connect to the server
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	// Start TLS
	tlsConfig := &tls.Config{
		ServerName: c.config.Host,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// Authenticate if credentials provided
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Set sender
	if err = client.Mail(c.config.From); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to add recipient %s: %w", recipient, err)
		}
	}

	// Send the email body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to send DATA command: %w", err)
	}

	_, err = w.Write(body)
	if err != nil {
		return fmt.Errorf("failed to write message body: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close DATA writer: %w", err)
	}

	return client.Quit()
}

// buildEmailBody constructs the email body with headers
func (c *Client) buildEmailBody(msg Message) []byte {
	var body strings.Builder

	// Add headers
	body.WriteString(fmt.Sprintf("From: %s\r\n", c.config.From))
	body.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(msg.To, ", ")))
	body.WriteString(fmt.Sprintf("Subject: %s\r\n", msg.Subject))
	body.WriteString("MIME-Version: 1.0\r\n")

	// If HTML body is provided, create multipart message
	if msg.HTMLBody != "" {
		boundary := "boundary-string-12345"
		body.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", boundary))
		body.WriteString("\r\n")

		// Plain text part
		body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		body.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
		body.WriteString("\r\n")
		body.WriteString(msg.Body)
		body.WriteString("\r\n")

		// HTML part
		body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		body.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
		body.WriteString("\r\n")
		body.WriteString(msg.HTMLBody)
		body.WriteString("\r\n")

		// End boundary
		body.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else {
		// Plain text only
		body.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
		body.WriteString("\r\n")
		body.WriteString(msg.Body)
	}

	return []byte(body.String())
}

