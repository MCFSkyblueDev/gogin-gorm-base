package repositories

import (
	"github.com/template/go-backend-gin-orm/model"
	"gorm.io/gorm"
)

type LoginRepository struct {
	repo *GenericRepository[model.Login]
}

func NewLoginRepository() *LoginRepository {
	return &LoginRepository{repo: NewGenericRepository[model.Login]()}
}

func (r *LoginRepository) CreateLogin(db *gorm.DB, login *model.Login) (*model.Login, error) {
	return r.repo.Create(db, login)
}

func (r *LoginRepository) FindLoginByEmail(db *gorm.DB, email string) (*model.Login, error) {
	return r.repo.FindByUniqueField(db, "email", email)
}
