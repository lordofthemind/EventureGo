package repositories

import (
	"context"

	"github.com/lordofthemind/EventureGo/internals/types"
)

type SuperUserRepositoryInterface interface {
	CreateSuperUser(ctx context.Context, superUser *types.SuperUserType) (*types.SuperUserType, error)
	FindSuperUserByEmail(ctx context.Context, email string) (*types.SuperUserType, error)
	FindSuperUserByUsername(ctx context.Context, username string) (*types.SuperUserType, error)
}
