package service

import (
	"errors"
	"time"

	"github.com/jordisetiawan/insurance-auth-service/internal/dto"
	"github.com/jordisetiawan/insurance-auth-service/internal/model"
	"github.com/jordisetiawan/insurance-auth-service/internal/repository"
	"github.com/jordisetiawan/insurance-auth-service/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user account is inactive")
)

type AuthService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository // Tambahkan RoleRepository
	secret   string
}

func NewAuthService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository, // Tambahkan RoleRepository
	secret string,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		roleRepo: roleRepo, // Inisialisasi RoleRepository
		secret:   secret,
	}
}

func (s *AuthService) Register(
	req dto.RegisterRequest,
) (*model.User, error) {

	_, err := s.userRepo.FindByEmail(req.Email)

	if err == nil {
		return nil, ErrEmailAlreadyExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := utils.HashPassword(req.Password)

	if err != nil {
		return nil, err
	}

	// Public self-registration always gets the default USER role.
	// Elevating a user's role is only possible via the admin-only
	// AssignRole flow (see AuthService.AssignRole).
	user := model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hash,
		FullName:     req.FullName,
		Role:         "USER",
		IsActive:     true,
	}

	return &user, s.userRepo.Create(&user)
}

func (s *AuthService) Login(req dto.LoginRequest) (*model.User, []string, string, string, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, "", "", ErrInvalidCredentials
		}
		return nil, nil, "", "", err
	}

	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, nil, "", "", ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, nil, "", "", ErrUserInactive
	}

	// Fetch permissions for the role
	permissions, err := s.roleRepo.GetPermissionsByRoleCode(user.Role)
	if err != nil {
		return nil, nil, "", "", err
	}

	// Access Token (15 menit)
	accessToken, err := utils.GenerateToken(user.ID.String(), user.Email, user.FullName, user.Role, permissions, s.secret, 15*time.Minute)
	if err != nil {
		return nil, nil, "", "", err
	}

	// Refresh Token (7 Hari)
	refreshToken, err := utils.GenerateToken(user.ID.String(), user.Email, user.FullName, user.Role, permissions, s.secret, 7*24*time.Hour)
	if err != nil {
		return nil, nil, "", "", err
	}

	return user, permissions, accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(refreshToken string) (string, error) {
	claims, err := utils.ValidateToken(refreshToken, s.secret)
	if err != nil {
		return "", err
	}

	// Fetch permissions for the role again to ensure token stays updated
	permissions, err := s.roleRepo.GetPermissionsByRoleCode(claims.Role)
	if err != nil {
		return "", err
	}

	// Generate new access token
	return utils.GenerateToken(claims.UserID, claims.Email, claims.FullName, claims.Role, permissions, s.secret, 15*time.Minute)
}

func (s *AuthService) GetAllUsers(page, limit int, search string) ([]model.User, int64, error) {
	return s.userRepo.FindAll(page, limit, search)
}

func (s *AuthService) GetUserByID(id string) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *AuthService) UpdateUser(id string, req dto.RegisterRequest) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check if new email is already taken by someone else
	if req.Email != "" && req.Email != user.Email {
		existing, _ := s.userRepo.FindByEmail(req.Email)
		if existing != nil {
			return nil, errors.New("new email is already in use")
		}
		user.Email = req.Email
	}

	user.FullName = req.FullName

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) AssignRole(id string, roleCode string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Verify role exists
	role, err := s.roleRepo.FindByCode(roleCode)
	if err != nil {
		role, err = s.roleRepo.FindByName(roleCode)
		if err != nil {
			return errors.New("role not found")
		}
	}

	user.Role = role.Code
	return s.userRepo.Update(user)
}

func (s *AuthService) AssignPermissionToRole(roleCode, permCode string) error {
	role, err := s.roleRepo.FindByCode(roleCode)
	if err != nil {
		return errors.New("role not found")
	}

	perm, err := s.roleRepo.FindPermissionByCode(permCode)
	if err != nil {
		return errors.New("permission not found")
	}

	return s.roleRepo.AssignPermission(role.ID, perm.ID)
}
