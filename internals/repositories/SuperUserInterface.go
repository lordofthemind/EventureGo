package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/types"
)

type SuperUserRepositoryInterface interface {
	CreateSuperUser(ctx context.Context, superUser *types.SuperUserType) (*types.SuperUserType, error)
	FindSuperUserByEmail(ctx context.Context, email string) (*types.SuperUserType, error)
	FindSuperUserByUsername(ctx context.Context, username string) (*types.SuperUserType, error)
	FindSuperUserByResetToken(ctx context.Context, token string) (*types.SuperUserType, error)
	UpdateResetToken(ctx context.Context, superUserID uuid.UUID, resetToken string) error
	UpdateSuperUser(ctx context.Context, superUser *types.SuperUserType) error
	FindSuperUserByOTP(ctx context.Context, otp string) (*types.SuperUserType, error)
	VerifySuperUserOTP(ctx context.Context, superUser *types.SuperUserType) error
}
