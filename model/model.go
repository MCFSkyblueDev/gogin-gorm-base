package model

import (
	"time"
)

type User struct {
	ID    uint   `gorm:"primaryKey"` // Standard field for the primary key
	Name  string // A regular string field
	Email string `gorm:"unique"` // A pointer to a string, allowing for null values
	Bio   *string
	// MemberNumber sql.NullString // Uses sql.NullString to handle nullable strings
	MemberNumber *string // Uses sql.NullString to handle nullable strings
	Avatar       *string
	CreatedAt    time.Time // Automatically managed by GORM for creation time
	UpdatedAt    time.Time // Automatically managed by GORM for update time

	Books []Book `gorm:"foreignKey:UserID"`
	Login *Login `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type Login struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex"`
	UserID       uint   `gorm:"unique"`
	HashPassword string

	User User `gorm:"constraint:OnDelete:CASCADE;"`
}

type Book struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Title     string
	Author    string
	Publisher string
	CreatedAt time.Time
	UpdatedAt time.Time

	User User `gorm:"constraint:OnDelete:CASCADE"`
}

func GetAllModels() []interface{} {
	return []interface{}{
		&User{},
		&Book{},
		&Login{},
	}
}
