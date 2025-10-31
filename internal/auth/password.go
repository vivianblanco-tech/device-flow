package auth

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const (
	// bcrypt cost factor (10-12 is recommended for production)
	bcryptCost = 12
	// Minimum password length
	minPasswordLength = 8
)

// Password validation regexes
var (
	uppercaseRegex   = regexp.MustCompile(`[A-Z]`)
	lowercaseRegex   = regexp.MustCompile(`[a-z]`)
	digitRegex       = regexp.MustCompile(`[0-9]`)
	specialCharRegex = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
)

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CheckPasswordHash compares a password with a hash
func CheckPasswordHash(password, hash string) bool {
	if password == "" || hash == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePassword validates password strength requirements
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}

	if len(password) < minPasswordLength {
		return errors.New("password must be at least 8 characters")
	}

	if !uppercaseRegex.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !lowercaseRegex.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	if !digitRegex.MatchString(password) {
		return errors.New("password must contain at least one digit")
	}

	if !specialCharRegex.MatchString(password) {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

