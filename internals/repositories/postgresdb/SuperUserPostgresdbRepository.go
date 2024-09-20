package postgresdb

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/types"
	"gorm.io/gorm"
)

type postgresSuperUserRepository struct {
	db *gorm.DB
}

func NewPostgresSuperUserRepository(db *gorm.DB) repositories.SuperUserRepositoryInterface {
	return &postgresSuperUserRepository{
		db: db,
	}
}

func (r *postgresSuperUserRepository) CreateSuperUser(ctx context.Context, superUser *types.SuperUserType) (*types.SuperUserType, error) {
	superUser.ID = uuid.New()
	superUser.CreatedAt = time.Now()
	superUser.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(superUser).Error; err != nil {
		return nil, err
	}
	return superUser, nil
}

func (r *postgresSuperUserRepository) FindSuperUserByEmail(ctx context.Context, email string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&superUser).Error
	return &superUser, err
}

func (r *postgresSuperUserRepository) FindSuperUserByUsername(ctx context.Context, username string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&superUser).Error
	return &superUser, err
}
