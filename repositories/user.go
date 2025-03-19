package repositories

import (
	"github.com/template/go-backend-gin-orm/dtos"
	"github.com/template/go-backend-gin-orm/model"
	"github.com/template/go-backend-gin-orm/utils"
	"gorm.io/gorm"
)

type UserRepository struct {
	repo *GenericRepository[model.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{repo: NewGenericRepository[model.User]()}
}

func (r *UserRepository) CreateUser(db *gorm.DB, user *model.User) (*model.User, error) {
	return r.repo.Create(db, user)
}

func (r *UserRepository) FindUserByEmail(db *gorm.DB, email string) (*model.User, error) {
	return r.repo.FindByUniqueField(db, "email", email)
}

func (r *UserRepository) UpdateUserInfo(db *gorm.DB, userId uint, updateInfo *dtos.UpdateUserInfoRequest, avatar *string) (*model.User, error) {
	modifiedUser, err := r.repo.GetById(db, userId)
	if err != nil {
		return nil, err
	}
	// Remove old avatar from local
	if modifiedUser.Avatar != nil {
		err := utils.RemoveFile(*modifiedUser.Avatar)
		if err != nil {
			return nil, err
		}
	}
	// Update
	if avatar != nil {
		modifiedUser.Avatar = avatar
	}
	modifiedUser.Name = *updateInfo.Name
	modifiedUser.Bio = updateInfo.Bio
	modifiedUser.MemberNumber = updateInfo.MemberNumber

	newUser, err := r.repo.Update(db, modifiedUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}
