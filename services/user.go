package services

import (
	"mime/multipart"

	"github.com/template/go-backend-gin-orm/dtos"
	"github.com/template/go-backend-gin-orm/model"
	"github.com/template/go-backend-gin-orm/repositories"
	"github.com/template/go-backend-gin-orm/utils"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repositories.UserRepository
	db       *gorm.DB
}

func NewUserService(userRepo *repositories.UserRepository, db *gorm.DB) *UserService {
	return &UserService{userRepo: userRepo, db: db}
}

func (s *UserService) UpdateUserInfo(userId uint, userInfo dtos.UpdateUserInfoRequest, avatar *multipart.FileHeader) (*model.User, error) {
	var filePath *string
	if avatar != nil {
		uploadedFilePath, err := utils.UploadFile(avatar, "uploads")
		if err != nil {
			return nil, err
		}
		filePath = &uploadedFilePath
	}
	updatedUser, err := s.userRepo.UpdateUserInfo(s.db, userId, &userInfo, filePath)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
