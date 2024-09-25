package services

import (
	"context"

	"github.com/lordofthemind/EventureGo/internals/newerrors"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/types"
	"github.com/lordofthemind/EventureGo/internals/utils"
	"golang.org/x/crypto/bcrypt"
)

type SuperUserService struct {
	repo repositories.SuperUserRepositoryInterface
}

func NewSuperUserService(repo repositories.SuperUserRepositoryInterface) SuperUserServiceInterface {
	return &SuperUserService{repo: repo}
}

// RegisterSuperUser registers a new superuser and returns the response object
func (s *SuperUserService) RegisterSuperUser(ctx context.Context, req *utils.RegisterSuperuserRequest) (*utils.SuperuserResponse, error) {
	// Validate email and username availability via validation layer
	if err := utils.ValidateUniqueness(ctx, req.Email, req.Username, s.repo); err != nil {
		return nil, err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, newerrors.Wrap(err, "failed to hash password")
	}

	// Create a new SuperUser entity
	superUserEntity := types.NewSuperUser(req.Email, req.FullName, req.Username, string(hashedPassword))

	// Store the SuperUser in the repository
	createdSuperUser, err := s.repo.CreateSuperUser(ctx, superUserEntity)
	if err != nil {
		return nil, newerrors.Wrap(err, "failed to create superuser")
	}

	// Prepare and return the response object
	return utils.CreateSuperuserResponse(createdSuperUser), nil
}
