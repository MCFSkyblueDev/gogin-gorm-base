package services

import (
	"errors"

	"github.com/template/go-backend-gin-orm/dtos"
	"github.com/template/go-backend-gin-orm/model"
	"github.com/template/go-backend-gin-orm/repositories"
	"github.com/template/go-backend-gin-orm/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo  *repositories.UserRepository
	loginRepo *repositories.LoginRepository
	db        *gorm.DB
}

func NewAuthService(userRepo *repositories.UserRepository, loginRepo *repositories.LoginRepository, db *gorm.DB) *AuthService {
	return &AuthService{userRepo: userRepo, loginRepo: loginRepo, db: db}
}

func (s *AuthService) RegisterUser(registerRequest dtos.RegisterRequest) (*model.User, error) {
	if registerRequest.Password != registerRequest.ConfirmPassword {
		return nil, errors.New("passwords do not match")
	}
	// Check if user already exists
	existingUser, _ := s.userRepo.FindUserByEmail(s.db, registerRequest.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword := utils.HashPassword(registerRequest.Password)
	newUser := &model.User{
		Name:  registerRequest.Name,
		Email: registerRequest.Email,
	}
	// Start Transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create User
	createdUser, err := s.userRepo.CreateUser(tx, newUser)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create Login'
	newLogin := &model.Login{
		Email:        registerRequest.Email,
		HashPassword: hashedPassword,
		UserID:       createdUser.ID,
	}
	_, err = s.loginRepo.CreateLogin(tx, newLogin)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit Transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *AuthService) Login(loginRequest dtos.LoginRequest) (*dtos.AuthResponse, error) {
	existingUser, _ := s.userRepo.FindUserByEmail(s.db, loginRequest.Email)
	if existingUser == nil {
		return nil, errors.New("email or password is incorrect")
	}
	// Check login record
	existingLogin, _ := s.loginRepo.FindLoginByEmail(s.db, existingUser.Email)
	if existingLogin == nil {
		return nil, errors.New("email or password is incorrect")
	}
	// Check password
	isMatch := utils.CheckPasswordHash(loginRequest.Password, existingLogin.HashPassword)
	if !isMatch {
		return nil, errors.New("email or password is incorrect")
	}
	token, err := utils.GenerateJWT(existingUser.ID, existingLogin.Email)
	if err != nil {
		return nil, errors.New("cannot get token")
	}

	return &dtos.AuthResponse{Token: token}, nil
}
