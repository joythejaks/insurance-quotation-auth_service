package repository

import (
	"github.com/google/uuid"
	"github.com/jordisetiawan/insurance-auth-service/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleRepository interface {
	FindByCode(code string) (*model.Role, error)
	FindByName(name string) (*model.Role, error)
	FindPermissionByCode(code string) (*model.Permission, error)
	AssignPermission(roleID, permissionID uuid.UUID) error
	GetPermissionsByRoleCode(code string) ([]string, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByCode(code string) (*model.Role, error) {
	var role model.Role
	if err := r.db.Where("code = ?", code).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindPermissionByCode(code string) (*model.Permission, error) {
	var perm model.Permission
	if err := r.db.Where("code = ?", code).First(&perm).Error; err != nil {
		return nil, err
	}
	return &perm, nil
}

func (r *roleRepository) AssignPermission(roleID, permissionID uuid.UUID) error {
	return r.db.Table("role_permissions").Clauses(clause.OnConflict{DoNothing: true}).Create(map[string]interface{}{
		"role_id":       roleID,
		"permission_id": permissionID,
	}).Error
}

func (r *roleRepository) GetPermissionsByRoleCode(code string) ([]string, error) {
	var permissions []string
	err := r.db.Table("permissions").
		Select("permissions.code").
		Joins("join role_permissions on role_permissions.permission_id = permissions.id").
		Joins("join roles on roles.id = role_permissions.role_id").
		Where("roles.code = ?", code).
		Pluck("code", &permissions).Error
	return permissions, err
}

func (r *roleRepository) FindByName(name string) (*model.Role, error) {
	var role model.Role
	if err := r.db.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
