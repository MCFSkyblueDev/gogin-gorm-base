package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/template/go-backend-gin-orm/dtos"
	"github.com/template/go-backend-gin-orm/services"
)

// @title Go-Gin-ORM API
// @version 1.0
// @description This is a sample API using Gin and GORM.
// @BasePath /api
// AuthRoutes defines authentication-related routes
func authRoutes(authGroup *gin.RouterGroup, authService *services.AuthService) {
	authGroup.POST("/register", registerUserHandler(authService))
	authGroup.POST("/login", login(authService))
}

// @Summary Register a new user
// @Description Creates a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dtos.RegisterRequest true "User registration data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} dtos.ErrorResponse "Bad Request"
// @Router /auth/register [post]
func registerUserHandler(authService *services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var registerRequest dtos.RegisterRequest
		//Validate
		if err := ctx.ShouldBindJSON(&registerRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: err.Error(), Status: http.StatusBadRequest})
			return
		}
		newUser, err := authService.RegisterUser(registerRequest)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: err.Error(), Status: http.StatusBadRequest})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": newUser})
	}
}

// @Summary Login account
// @Description Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dtos.LoginRequest true "User login data"
// @Success 201 {object} dtos.AuthResponse ư
// @Failure 400 {object} dtos.ErrorResponse "Bad Request"
// @Router /auth/login [post]
func login(authService *services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginRequest dtos.LoginRequest
		// Validate
		if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: err.Error(), Status: http.StatusBadRequest})
			return
		}
		token, err := authService.Login(loginRequest)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: err.Error(), Status: http.StatusBadRequest})
			return
		}
		ctx.JSON(http.StatusCreated, token)
	}
}

// RegisterRoutes registers all route groups
func AuthRoutes(router *gin.RouterGroup, authService *services.AuthService) {
	authRoutes(router.Group("auth"), authService)
}
