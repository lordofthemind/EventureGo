package utils

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/types"
)

type RegisterSuperuserRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	FullName string `json:"full_name" binding:"required,min=3,max=32" validate:"required,min=3,max=32"`
	Username string `json:"username" binding:"required,min=3,max=32,alphanum" validate:"required,min=3,max=32,alphanum"`
	Password string `json:"password" binding:"required,min=8" validate:"required,min=8"`
}

type SuperuserResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Username     string    `json:"username"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Is2FAEnabled bool      `json:"is_2fa_enabled"`
}

// CreateSuperuserResponse creates a response object for SuperUserType
func CreateSuperuserResponse(superUser *types.SuperUserType) *SuperuserResponse {
	return &SuperuserResponse{
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

// ValidateUniqueness checks if email and username are unique
func ValidateUniqueness(ctx context.Context, email, username string, repo repositories.SuperUserRepositoryInterface) error {
	if _, err := repo.FindSuperUserByEmail(ctx, email); err == nil {
		return errors.New("email already in use")
	}

	if _, err := repo.FindSuperUserByUsername(ctx, username); err == nil {
		return errors.New("username already in use")
	}

	return nil
}
