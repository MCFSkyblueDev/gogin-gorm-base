package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/template/go-backend-gin-orm/dtos"
	"github.com/template/go-backend-gin-orm/middlewares"
	"github.com/template/go-backend-gin-orm/services"
)

func userRoutes(userGroup *gin.RouterGroup, userService *services.UserService) {
	// userGroup.Use(middlewares.AuthMiddleware())
	userGroup.POST("update-info", middlewares.AuthMiddleware(), updateUserInfo(userService))
}

// @Summary Update user information
// @Description Updates the user's profile information, including name, bio, and avatar
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Authorization"
// @Param name formData string false "User's name (min: 2, max: 20)"
// @Param bio formData string false "User's bio (min: 12 characters)"
// @Param member_number formData string false "User's member number"
// @Param avatar formData file false "User's avatar image"
// @Success 200 {object} map[string]interface{} "User info updated successfully"
// @Failure 400 {object} dtos.ErrorResponse "Bad Request"
// @Failure 401 {object} dtos.ErrorResponse "Unauthorized"
// @Router /user/update-info [post]
func updateUserInfo(userService *services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get User from Token
		claims, exists := ctx.Get("user")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userClaims, ok := claims.(*dtos.UserClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user data"})
			return
		}

		var updateInfo dtos.UpdateUserInfoRequest
		if err := ctx.ShouldBind(&updateInfo); err != nil {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: err.Error(), Status: http.StatusBadRequest})
			return
		}

		// Handle file upload
		file, err := ctx.FormFile("avatar")
		if err != nil && err != http.ErrMissingFile {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
			return
		}
		log.Println(file)

		userId := userClaims.UserID
		updatedUser, err := userService.UpdateUserInfo(uint(userId), updateInfo, file)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "User info updated successfully", "user": updatedUser})

	}
}

func UserRoutes(router *gin.RouterGroup, userService *services.UserService) {
	userRoutes(router.Group("user"), userService)
}
