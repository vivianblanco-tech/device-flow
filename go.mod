module github.com/yourusername/laptop-tracking-system

go 1.24.0

toolchain go1.24.6

require (
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.43.0
	golang.org/x/oauth2 v0.32.0
)

require cloud.google.com/go/compute/metadata v0.3.0 // indirect
