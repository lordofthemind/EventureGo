package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
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
	// Check if the email or username already exists
	if _, err := s.repo.FindSuperUserByEmail(ctx, superUser.Email); err == nil {
		return nil, errors.New("email already in use")
	}

	if _, err := s.repo.FindSuperUserByUsername(ctx, superUser.Username); err == nil {
		return nil, errors.New("username already in use")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(superUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create a new SuperUser entity
	superUserEntity := &types.SuperUserType{
		ID:             uuid.New(),
		Email:          superUser.Email,
		FullName:       superUser.FullName,
		Username:       superUser.Username,
		HashedPassword: string(hashedPassword),
		Role:           "guest", // default role
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Store the SuperUser in the repository
	createdSuperUser, err := s.repo.CreateSuperUser(ctx, superUserEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to create superuser: %w", err)
	}

	// Prepare and return the response object
	response := &utils.SuperuserResponse{
		ID:           createdSuperUser.ID,
		Email:        createdSuperUser.Email,
		FullName:     createdSuperUser.FullName,
		Username:     createdSuperUser.Username,
		Role:         createdSuperUser.Role,
		CreatedAt:    createdSuperUser.CreatedAt,
		UpdatedAt:    createdSuperUser.UpdatedAt,
		Is2FAEnabled: createdSuperUser.Is2FAEnabled,
	}

	return response, nil
}
