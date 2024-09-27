package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/types"
)

// LogInSuperuserRequest represents the request structure for logging in a superuser.
type LogInSuperuserRequest struct {
	Email    string `json:"email" binding:"omitempty,email" validate:"omitempty,email"`    // Email of the superuser
	Username string `json:"username" binding:"omitempty,min=3" validate:"omitempty,min=3"` // Username of the superuser
	Password string `json:"password" binding:"required,min=8" validate:"required,min=8"`   // Password for authentication
}

// LoginSuperuserResponse represents the response structure for a successful login.
type LoginSuperuserResponse struct {
	ID           uuid.UUID `json:"id"`             // Unique identifier of the superuser
	Email        string    `json:"email"`          // Email of the superuser
	Username     string    `json:"username"`       // Username of the superuser
	FullName     string    `json:"full_name"`      // Full name of the superuser
	Token        string    `json:"token"`          // Authentication token for the session
	Is2FAEnabled bool      `json:"is_2fa_enabled"` // Indicates if two-factor authentication is enabled
}

// RegisterSuperuserRequest represents the request structure for registering a new superuser.
type RegisterSuperuserRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`                                    // Email of the new superuser
	FullName string `json:"full_name" binding:"required,min=3,max=32" validate:"required,min=3,max=32"`                  // Full name of the new superuser
	Username string `json:"username" binding:"required,min=3,max=32,alphanum" validate:"required,min=3,max=32,alphanum"` // Username of the new superuser
	Password string `json:"password" binding:"required,min=8" validate:"required,min=8"`                                 // Password for the new superuser
}

// RegisterSuperuserResponse represents the response structure for a successful registration.
type RegisterSuperuserResponse struct {
	ID           uuid.UUID `json:"id"`             // Unique identifier of the superuser
	Email        string    `json:"email"`          // Email of the new superuser
	FullName     string    `json:"full_name"`      // Full name of the new superuser
	Username     string    `json:"username"`       // Username of the new superuser
	Role         string    `json:"role"`           // Role assigned to the new superuser
	CreatedAt    time.Time `json:"created_at"`     // Timestamp when the superuser was created
	UpdatedAt    time.Time `json:"updated_at"`     // Timestamp when the superuser was last updated
	Is2FAEnabled bool      `json:"is_2fa_enabled"` // Indicates if two-factor authentication is enabled
}

// CreateSuperuserResponse creates a RegisterSuperuserResponse from a SuperUserType.
func CreateSuperuserResponse(superUser *types.SuperUserType) *RegisterSuperuserResponse {
	return &RegisterSuperuserResponse{
		ID:           superUser.ID,
		Email:        superUser.Email,
		FullName:     superUser.FullName,
		Username:     superUser.Username,
		Role:         superUser.Role,
		CreatedAt:    superUser.CreatedAt,
		UpdatedAt:    superUser.UpdatedAt,
		Is2FAEnabled: superUser.Is2FAEnabled,
	}
}

// ValidateUniqueness checks if the provided email and username are unique within the repository.
func ValidateUniqueness(ctx context.Context, email, username string, repo repositories.SuperUserRepositoryInterface) error {
	if _, err := repo.FindSuperUserByEmail(ctx, email); err == nil {
		return errors.New("email already in use") // Return error if email is already in use
	}

	if _, err := repo.FindSuperUserByUsername(ctx, username); err == nil {
		return errors.New("username already in use") // Return error if username is already in use
	}

	return nil // Return nil if both email and username are unique
}

// GenerateResetToken generates a unique token for password resets.
// The token consists of a UUID and a random hex string.
func GenerateResetToken() string {
	// Generate a UUID for the token
	uuidPart := uuid.New().String()

	// Generate a random hex string
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		panic("failed to generate random bytes for token") // Panic if random bytes cannot be generated
	}
	randomPart := hex.EncodeToString(randomBytes)

	// Combine both parts to create the final reset token
	resetToken := uuidPart + "-" + randomPart
	return resetToken
}
