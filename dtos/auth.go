package dtos

import "github.com/golang-jwt/jwt/v5"

type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Name            string `json:"name" binding:"required,min=2,max=20"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type UserClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
