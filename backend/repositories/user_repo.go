package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

// UserRepository interface
type UserRepository interface {
	Create(user *models.User) error
	GetByUsername(username string) (*models.User, error)
}

// GormUserRepository implementation
type GormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
