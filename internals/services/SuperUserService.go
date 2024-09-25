package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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
func (s *SuperUserService) RegisterSuperUser(ctx context.Context, superUser *utils.RegisterSuperuserRequest) (*utils.SuperuserResponse, error) {
	// Validate email and username availability
	if err := s.validateSuperUserUniqueness(ctx, superUser.Email, superUser.Username); err != nil {
		return nil, err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(superUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, newerrors.Wrap(err, "failed to hash password")
	}

	// Create a new SuperUser entity
	superUserEntity := s.createSuperUserEntity(superUser, hashedPassword)

	// Store the SuperUser in the repository
	createdSuperUser, err := s.repo.CreateSuperUser(ctx, superUserEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to create superuser: %w", err)
	}

	// Prepare and return the response object
	return s.createSuperuserResponse(createdSuperUser), nil
}

func (s *SuperUserService) validateSuperUserUniqueness(ctx context.Context, email, username string) error {
	if _, err := s.repo.FindSuperUserByEmail(ctx, email); err == nil {
		return newerrors.NewValidationError("email already in use")
	}

	if _, err := s.repo.FindSuperUserByUsername(ctx, username); err == nil {
		return newerrors.NewValidationError("username already in use")
	}
	return nil
}

func (s *SuperUserService) createSuperUserEntity(superUser *utils.RegisterSuperuserRequest, hashedPassword []byte) *types.SuperUserType {
	return &types.SuperUserType{
		ID:             uuid.New(),
		Email:          superUser.Email,
		FullName:       superUser.FullName,
		Username:       superUser.Username,
		HashedPassword: string(hashedPassword),
		Role:           "guest", // default role
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func (s *SuperUserService) createSuperuserResponse(superUser *types.SuperUserType) *utils.SuperuserResponse {
	return &utils.SuperuserResponse{
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
