package utils

import "github.com/google/uuid"

type LogInSuperuserRequest struct {
	Email    string `json:"email" binding:"omitempty,email" validate:"omitempty,email"`
	Username string `json:"username" binding:"omitempty,min=3" validate:"omitempty,min=3"`
	Password string `json:"password" binding:"required,min=8" validate:"required,min=8"`
}

type LoginSuperuserResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	FullName     string    `json:"full_name"`
	Token        string    `json:"token"`
	Is2FAEnabled bool      `json:"is_2fa_enabled"`
}
