package startup

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/template/go-backend-gin-orm/config"
	"github.com/template/go-backend-gin-orm/docs"
	"github.com/template/go-backend-gin-orm/repositories"
	"github.com/template/go-backend-gin-orm/routes"
	"github.com/template/go-backend-gin-orm/services"
)

type Server struct {
	app *gin.Engine
}

func NewServer() *Server {
	dbInstance := NewPostgresDatabase()
	dbInstance.MigrateDatabase()

	app := gin.Default()
	app.Static("/upload", "./uploads")

	// Initialize Repositories
	// userRepo := repositories.NewUserRepository(dbInstance.DB)
	userRepo := repositories.NewUserRepository()
	loginRepo := repositories.NewLoginRepository()

	// Initialize Services
	authService := services.NewAuthService(userRepo, loginRepo, dbInstance.DB)
	userService := services.NewUserService(userRepo, dbInstance.DB)

	// Register Routes
	api := app.Group("/api")
	routes.AuthRoutes(api, authService)
	routes.UserRoutes(api, userService)

	// Swagger
	docs.SwaggerInfo.BasePath = "/api"
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Server{app: app}
}

func (s *Server) Start() error {
	env := config.NewEnv(".env", true)
	log.Println("Starting server on port " + env.Port)
	return s.app.Run(fmt.Sprintf(":%s", env.Port))
}
