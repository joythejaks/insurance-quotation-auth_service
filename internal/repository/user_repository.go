package repository

import (
	"github.com/jordisetiawan/insurance-auth-service/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id string) (*model.User, error)
	FindAll(page, limit int, search string) ([]model.User, int64, error)
	Update(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(page, limit int, search string) ([]model.User, int64, error) {
	var users []model.User
	var count int64
	query := r.db.Model(&model.User{})
	if search != "" {
		query = query.Where("full_name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query.Count(&count)
	err := query.Offset((page - 1) * limit).Limit(limit).Find(&users).Error
	return users, count, err
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}
