package repositories

import "gorm.io/gorm"

type IGenericRepository[T any] interface {
	Create(db *gorm.DB, entity *T) (*T, error)
	GetAll(db *gorm.DB) ([]T, error)
	GetById(db *gorm.DB, id uint) (*T, error)
	Update(db *gorm.DB, entity *T) (*T, error)
	Delete(db *gorm.DB, id uint) (bool, error)
	FindByUniqueField(db *gorm.DB, field string, value interface{}) (*T, error)
}

type GenericRepository[T any] struct{}

func NewGenericRepository[T any]() *GenericRepository[T] {
	return &GenericRepository[T]{}
}

func (r *GenericRepository[T]) Create(db *gorm.DB, entity *T) (*T, error) {
	if err := db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *GenericRepository[T]) GetAll(db *gorm.DB) ([]T, error) {
	var entities []T
	if err := db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *GenericRepository[T]) GetById(db *gorm.DB, id uint) (*T, error) {
	var entity T
	if err := db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GenericRepository[T]) Update(db *gorm.DB, entity *T) (*T, error) {
	if err := db.Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *GenericRepository[T]) Delete(db *gorm.DB, id uint) (bool, error) {
	result := db.Delete(new(T), id)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *GenericRepository[T]) FindByUniqueField(db *gorm.DB, field string, value interface{}) (*T, error) {
	var entity T
	if err := db.Where(field+"= ?", value).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}
