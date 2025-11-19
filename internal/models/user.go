package models

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"time"
)

// UserRole represents the role of a user in the system
type UserRole string

// User role constants
const (
	RoleLogistics      UserRole = "logistics"
	RoleClient         UserRole = "client"
	RoleWarehouse      UserRole = "warehouse"
	RoleProjectManager UserRole = "project_manager"
)

// User represents a user in the system
type User struct {
	ID                int64     `json:"id" db:"id"`
	Email             string    `json:"email" db:"email"`
	PasswordHash      string    `json:"-" db:"password_hash"`
	Role              UserRole  `json:"role" db:"role"`
	ClientCompanyID   *int64    `json:"client_company_id,omitempty" db:"client_company_id"`
	ClientCompanyName string    `json:"client_company_name,omitempty" db:"-"` // Populated via JOIN queries
	GoogleID          *string   `json:"google_id,omitempty" db:"google_id"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Validate validates the User model
func (u *User) Validate() error {
	// Email validation
	if u.Email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}

	// Role validation
	if u.Role == "" {
		return errors.New("role is required")
	}
	if !IsValidRole(u.Role) {
		return errors.New("invalid role")
	}

	// Either password_hash or google_id must be provided
	if u.PasswordHash == "" && u.GoogleID == nil {
		return errors.New("either password_hash or google_id must be provided")
	}

	return nil
}

// IsValidRole checks if a given role is valid
func IsValidRole(role UserRole) bool {
	switch role {
	case RoleLogistics, RoleClient, RoleWarehouse, RoleProjectManager:
		return true
	}
	return false
}

// HasRole checks if the user has the specified role
func (u *User) HasRole(role UserRole) bool {
	return u.Role == role
}

// IsGoogleUser checks if the user authenticated via Google OAuth
func (u *User) IsGoogleUser() bool {
	return u.GoogleID != nil && *u.GoogleID != ""
}

// TableName returns the table name for the User model
func (u *User) TableName() string {
	return "users"
}

// BeforeCreate sets the timestamps before creating a user
func (u *User) BeforeCreate() {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
}

// BeforeUpdate sets the updated_at timestamp before updating a user
func (u *User) BeforeUpdate() {
	u.UpdatedAt = time.Now()
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(db *sql.DB) ([]User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, u.role, u.client_company_id, u.google_id, u.created_at, u.updated_at,
		       cc.name as client_company_name
		FROM users u
		LEFT JOIN client_companies cc ON cc.id = u.client_company_id
		ORDER BY u.email ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var clientCompanyName sql.NullString
		var googleID sql.NullString

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
			&user.ClientCompanyID,
			&googleID,
			&user.CreatedAt,
			&user.UpdatedAt,
			&clientCompanyName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if googleID.Valid {
			user.GoogleID = &googleID.String
		}
		if clientCompanyName.Valid {
			user.ClientCompanyName = clientCompanyName.String
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// GetUserByID retrieves a user by its ID
func GetUserByID(db *sql.DB, id int64) (*User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, u.role, u.client_company_id, u.google_id, u.created_at, u.updated_at,
		       cc.name as client_company_name
		FROM users u
		LEFT JOIN client_companies cc ON cc.id = u.client_company_id
		WHERE u.id = $1
	`

	var user User
	var clientCompanyName sql.NullString
	var googleID sql.NullString

	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.ClientCompanyID,
		&googleID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&clientCompanyName,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found with id %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if clientCompanyName.Valid {
		user.ClientCompanyName = clientCompanyName.String
	}

	return &user, nil
}

// CreateUser creates a new user in the database
func CreateUser(db *sql.DB, user *User) error {
	// Validate user
	if err := user.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Set timestamps
	user.BeforeCreate()

	query := `
		INSERT INTO users (email, password_hash, role, client_company_id, google_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := db.QueryRow(
		query,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.ClientCompanyID,
		user.GoogleID,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// UpdateUser updates an existing user in the database
func UpdateUser(db *sql.DB, user *User) error {
	// Validate user
	if err := user.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Update timestamp
	user.BeforeUpdate()

	query := `
		UPDATE users
		SET email = $1, password_hash = $2, role = $3, client_company_id = $4, google_id = $5, updated_at = $6
		WHERE id = $7
	`

	result, err := db.Exec(
		query,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.ClientCompanyID,
		user.GoogleID,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found with id %d", user.ID)
	}

	return nil
}

// DeleteUser deletes a user from the database
func DeleteUser(db *sql.DB, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found with id %d", id)
	}

	return nil
}

