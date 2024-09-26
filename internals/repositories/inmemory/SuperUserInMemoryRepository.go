package inmemory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/types"
)

type inMemorySuperUserRepository struct {
	mu         sync.RWMutex
	superUsers map[uuid.UUID]*types.SuperUserType
}

func NewInMemorySuperUserRepository() repositories.SuperUserRepositoryInterface {
	return &inMemorySuperUserRepository{
		superUsers: make(map[uuid.UUID]*types.SuperUserType),
	}
}

func (r *inMemorySuperUserRepository) CreateSuperUser(ctx context.Context, superUser *types.SuperUserType) (*types.SuperUserType, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	superUser.ID = uuid.New()
	superUser.CreatedAt = time.Now()
	superUser.UpdatedAt = time.Now()

	r.superUsers[superUser.ID] = superUser
	return superUser, nil
}

func (r *inMemorySuperUserRepository) FindSuperUserByEmail(ctx context.Context, email string) (*types.SuperUserType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, su := range r.superUsers {
		if su.Email == email {
			return su, nil
		}
	}
	return nil, nil // Return nil instead of an error for not found
}

func (r *inMemorySuperUserRepository) FindSuperUserByUsername(ctx context.Context, username string) (*types.SuperUserType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, su := range r.superUsers {
		if su.Username == username {
			return su, nil
		}
	}
	return nil, nil // Return nil instead of an error for not found
}

func (r *inMemorySuperUserRepository) FindSuperUserByResetToken(ctx context.Context, token string) (*types.SuperUserType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, su := range r.superUsers {
		if su.ResetToken != nil && *su.ResetToken == token {
			return su, nil
		}
	}
	return nil, errors.New("superuser not found")
}

func (r *inMemorySuperUserRepository) UpdateResetToken(ctx context.Context, superUserID uuid.UUID, resetToken string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	superUser, exists := r.superUsers[superUserID]
	if !exists {
		return errors.New("superuser not found")
	}

	superUser.ResetToken = &resetToken
	superUser.UpdatedAt = time.Now()
	return nil
}

func (r *inMemorySuperUserRepository) UpdateSuperUser(ctx context.Context, superUser *types.SuperUserType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existingSuperUser, exists := r.superUsers[superUser.ID]
	if !exists {
		return errors.New("superuser not found")
	}

	existingSuperUser = superUser
	existingSuperUser.UpdatedAt = time.Now()
	return nil
}
