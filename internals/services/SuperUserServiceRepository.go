package services

import (
	"context"

	"github.com/lordofthemind/EventureGo/internals/utils"
)

type SuperUserServiceInterface interface {
	RegisterSuperUser(ctx context.Context, superUser *utils.RegisterSuperuserRequest) (*utils.RegisterSuperuserResponse, error)
	LogInSuperuser(ctx context.Context, loginRequest *utils.LogInSuperuserRequest) (*utils.LoginSuperuserResponse, error)
	SendPasswordResetEmail(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}
