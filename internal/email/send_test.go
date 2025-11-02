package email

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

// mockSMTPServer creates a simple mock SMTP server for testing
type mockSMTPServer struct {
	listener    net.Listener
	messages    []string
	shouldError bool
}

func newMockSMTPServer(port int) (*mockSMTPServer, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, err
	}

	server := &mockSMTPServer{
		listener: listener,
		messages: []string{},
	}

	go server.serve()

	return server, nil
}

func (s *mockSMTPServer) serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handleConnection(conn)
	}
}

func (s *mockSMTPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Send greeting
	conn.Write([]byte("220 mock.smtp.server ESMTP\r\n"))

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		command := string(buf[:n])
		s.messages = append(s.messages, command)

		// Handle SMTP commands
		if strings.HasPrefix(command, "EHLO") || strings.HasPrefix(command, "HELO") {
			conn.Write([]byte("250 Hello\r\n"))
		} else if strings.HasPrefix(command, "MAIL FROM") {
			conn.Write([]byte("250 OK\r\n"))
		} else if strings.HasPrefix(command, "RCPT TO") {
			conn.Write([]byte("250 OK\r\n"))
		} else if strings.HasPrefix(command, "DATA") {
			conn.Write([]byte("354 Start mail input\r\n"))
		} else if strings.HasPrefix(command, ".") {
			conn.Write([]byte("250 OK\r\n"))
		} else if strings.HasPrefix(command, "QUIT") {
			conn.Write([]byte("221 Bye\r\n"))
			return
		} else {
			conn.Write([]byte("250 OK\r\n"))
		}
	}
}

func (s *mockSMTPServer) close() {
	s.listener.Close()
}

func TestClient_Send(t *testing.T) {
	// Skip in short mode as this requires network setup
	if testing.Short() {
		t.Skip("Skipping email send test in short mode")
	}

	// Start a mock SMTP server
	mockServer, err := newMockSMTPServer(2525)
	if err != nil {
		t.Fatalf("Failed to start mock SMTP server: %v", err)
	}
	defer mockServer.close()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Create a client pointing to the mock server
	client, err := NewClient(Config{
		Host: "127.0.0.1",
		Port: 2525,
		From: "sender@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name    string
		message Message
		wantErr bool
	}{
		{
			name: "send valid email",
			message: Message{
				To:      []string{"recipient@example.com"},
				Subject: "Test Subject",
				Body:    "Test body content",
			},
			wantErr: false,
		},
		{
			name: "send email with multiple recipients",
			message: Message{
				To:      []string{"recipient1@example.com", "recipient2@example.com"},
				Subject: "Test Subject",
				Body:    "Test body content",
			},
			wantErr: false,
		},
		{
			name: "send email with HTML body",
			message: Message{
				To:       []string{"recipient@example.com"},
				Subject:  "Test Subject",
				Body:     "Plain text body",
				HTMLBody: "<html><body><h1>HTML Body</h1></body></html>",
			},
			wantErr: false,
		},
		{
			name: "fail to send invalid message",
			message: Message{
				To:      []string{},
				Subject: "Test Subject",
				Body:    "Test body",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Send(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_buildEmailBody(t *testing.T) {
	client, err := NewClient(Config{
		Host: "smtp.example.com",
		Port: 587,
		From: "sender@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name           string
		message        Message
		expectedParts  []string
		notExpected    []string
	}{
		{
			name: "plain text email",
			message: Message{
				To:      []string{"recipient@example.com"},
				Subject: "Test Subject",
				Body:    "This is a plain text body",
			},
			expectedParts: []string{
				"From: sender@example.com",
				"To: recipient@example.com",
				"Subject: Test Subject",
				"Content-Type: text/plain",
				"This is a plain text body",
			},
			notExpected: []string{
				"multipart",
				"text/html",
			},
		},
		{
			name: "multipart email with HTML",
			message: Message{
				To:       []string{"recipient@example.com"},
				Subject:  "HTML Email",
				Body:     "Plain text version",
				HTMLBody: "<html><body>HTML version</body></html>",
			},
			expectedParts: []string{
				"From: sender@example.com",
				"To: recipient@example.com",
				"Subject: HTML Email",
				"Content-Type: multipart/alternative",
				"Content-Type: text/plain",
				"Plain text version",
				"Content-Type: text/html",
				"<html><body>HTML version</body></html>",
			},
		},
		{
			name: "multiple recipients",
			message: Message{
				To:      []string{"recipient1@example.com", "recipient2@example.com", "recipient3@example.com"},
				Subject: "Multiple Recipients",
				Body:    "Body content",
			},
			expectedParts: []string{
				"To: recipient1@example.com, recipient2@example.com, recipient3@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := client.buildEmailBody(tt.message)
			bodyStr := string(body)

			// Check that all expected parts are present
			for _, expected := range tt.expectedParts {
				if !strings.Contains(bodyStr, expected) {
					t.Errorf("buildEmailBody() missing expected part: %s", expected)
				}
			}

			// Check that not-expected parts are absent
			for _, notExpected := range tt.notExpected {
				if strings.Contains(bodyStr, notExpected) {
					t.Errorf("buildEmailBody() contains unexpected part: %s", notExpected)
				}
			}
		})
	}
}

