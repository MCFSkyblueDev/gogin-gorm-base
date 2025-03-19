package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/template/go-backend-gin-orm/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Extract token from Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Missing Authorization header"})
			ctx.Abort()
			return
		}

		// Extract token part from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header format"})
			ctx.Abort()
			return
		}

		token := tokenParts[1]

		// Validate token
		userClaims, err := utils.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token", "error": err.Error()})
			ctx.Abort()
			return
		}

		// Attach user data to context
		ctx.Set("user", userClaims)

		// Continue to next middleware/handler
		ctx.Next()
	}
}
