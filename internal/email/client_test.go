package email

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid configuration",
			config: Config{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "test@example.com",
				Password: "password123",
				From:     "noreply@example.com",
			},
			wantErr: false,
		},
		{
			name: "missing host",
			config: Config{
				Host:     "",
				Port:     587,
				Username: "test@example.com",
				Password: "password123",
				From:     "noreply@example.com",
			},
			wantErr: true,
			errMsg:  "SMTP host is required",
		},
		{
			name: "invalid port - zero",
			config: Config{
				Host:     "smtp.example.com",
				Port:     0,
				Username: "test@example.com",
				Password: "password123",
				From:     "noreply@example.com",
			},
			wantErr: true,
			errMsg:  "SMTP port must be between 1 and 65535",
		},
		{
			name: "invalid port - too high",
			config: Config{
				Host:     "smtp.example.com",
				Port:     70000,
				Username: "test@example.com",
				Password: "password123",
				From:     "noreply@example.com",
			},
			wantErr: true,
			errMsg:  "SMTP port must be between 1 and 65535",
		},
		{
			name: "missing from address",
			config: Config{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "test@example.com",
				Password: "password123",
				From:     "",
			},
			wantErr: true,
			errMsg:  "from address is required",
		},
		{
			name: "optional username and password",
			config: Config{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "",
				Password: "",
				From:     "noreply@example.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err.Error() != tt.errMsg {
					t.Errorf("NewClient() error message = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if client == nil {
				t.Error("NewClient() returned nil client")
				return
			}

			// Verify client has the correct configuration
			if client.config.Host != tt.config.Host {
				t.Errorf("Client host = %v, want %v", client.config.Host, tt.config.Host)
			}
			if client.config.Port != tt.config.Port {
				t.Errorf("Client port = %v, want %v", client.config.Port, tt.config.Port)
			}
			if client.config.From != tt.config.From {
				t.Errorf("Client from = %v, want %v", client.config.From, tt.config.From)
			}
		})
	}
}

func TestClient_BuildMessage(t *testing.T) {
	client, err := NewClient(Config{
		Host: "smtp.example.com",
		Port: 587,
		From: "noreply@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name    string
		message Message
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid message",
			message: Message{
				To:      []string{"recipient@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr: false,
		},
		{
			name: "multiple recipients",
			message: Message{
				To:      []string{"recipient1@example.com", "recipient2@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr: false,
		},
		{
			name: "missing recipients",
			message: Message{
				To:      []string{},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr: true,
			errMsg:  "at least one recipient is required",
		},
		{
			name: "missing subject",
			message: Message{
				To:      []string{"recipient@example.com"},
				Subject: "",
				Body:    "Test Body",
			},
			wantErr: true,
			errMsg:  "subject is required",
		},
		{
			name: "missing body",
			message: Message{
				To:      []string{"recipient@example.com"},
				Subject: "Test Subject",
				Body:    "",
			},
			wantErr: true,
			errMsg:  "body is required",
		},
		{
			name: "with HTML body",
			message: Message{
				To:       []string{"recipient@example.com"},
				Subject:  "Test Subject",
				Body:     "Plain text body",
				HTMLBody: "<html><body>HTML body</body></html>",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.BuildMessage(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("BuildMessage() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
