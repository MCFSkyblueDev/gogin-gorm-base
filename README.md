# Go Backend with Gin + GORM + Migration Guide

Complete guide for setting up a Go backend project using Gin framework, GORM ORM, and database migrations.

## 📦 Project Setup

### 1. Initialize Project

```bash
# Create project directory
mkdir go-backend-gin-orm
cd go-backend-gin-orm

# Initialize Go module
go mod init github.com/template/go-backend-gin-orm
```

### 2. Install Dependencies

```bash
# Core dependencies
go get -u github.com/spf13/viper           # Configuration management
go get -u gorm.io/gorm                     # ORM
go get -u gorm.io/driver/postgres          # PostgreSQL driver
go get -u github.com/gin-gonic/gin         # Web framework

# Authentication
go get github.com/golang-jwt/jwt/v5        # JWT authentication

# API Documentation
go get -u github.com/swaggo/gin-swagger    # Swagger for Gin
go get -u github.com/swaggo/files          # Static files for Swagger
```

### 3. Install Development Tools

```bash
# Database migration tool
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Swagger documentation generator
go install github.com/swaggo/swag/cmd/swag@latest

# Optional: Atlas for schema management
go install ariga.io/atlas/cmd/atlas@latest
```

## 🗂️ Project Structure

```
go-backend-gin-orm/
├── cmd/
│   └── main.go                 # Application entry point
├── config/
│   ├── config.go              # Configuration loader
│   └── config.yaml            # Configuration file
├── internal/
│   ├── models/                # Database models
│   │   └── user.go
│   ├── repositories/          # Data access layer
│   │   └── user_repository.go
│   ├── services/              # Business logic
│   │   └── user_service.go
│   ├── handlers/              # HTTP handlers
│   │   └── user_handler.go
│   ├── middleware/            # Middleware functions
│   │   └── auth.go
│   └── routes/                # Route definitions
│       └── routes.go
├── migrations/                # SQL migration files
│   ├── 000001_create_users_table.up.sql
│   └── 000001_create_users_table.down.sql
├── docs/                      # Swagger documentation (auto-generated)
├── .env                       # Environment variables
├── go.mod
└── go.sum
```

## ⚙️ Configuration

### config/config.yaml

```yaml
server:
  port: 5050
  mode: debug  # debug, release

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: myapp_db
  sslmode: disable

jwt:
  secret: your-secret-key-change-in-production
  expiration: 24h
```

### config/config.go

```go
package config

import (
    "github.com/spf13/viper"
    "log"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
}

type ServerConfig struct {
    Port string
    Mode string
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type JWTConfig struct {
    Secret     string
    Expiration string
}

func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")
    viper.AddConfigPath(".")

    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        log.Printf("Error reading config file: %v", err)
        return nil, err
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}
```

## 🗄️ Database Migrations

### Using golang-migrate

#### Install with PostgreSQL support

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

#### Create Migration Files

```bash
# Create a new migration
migrate create -ext sql -dir migrations -seq create_users_table

# This creates two files:
# migrations/000001_create_users_table.up.sql
# migrations/000001_create_users_table.down.sql
```

#### Example Migration Files

**migrations/000001_create_users_table.up.sql**

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
```

**migrations/000001_create_users_table.down.sql**

```sql
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

#### Run Migrations

```bash
# Set database connection string
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/myapp_db?sslmode=disable"

# Run all pending migrations
migrate -database $DATABASE_URL -path migrations up

# Rollback last migration
migrate -database $DATABASE_URL -path migrations down 1

# Rollback all migrations
migrate -database $DATABASE_URL -path migrations down

# Check migration version
migrate -database $DATABASE_URL -path migrations version

# Force version (use with caution)
migrate -database $DATABASE_URL -path migrations force 1
```

### Using Atlas (Recommend)

Atlas is a modern schema migration tool that can automatically generate migrations from GORM models.

#### Install Atlas Provider for GORM

```bash
# Install Atlas
go install ariga.io/atlas/cmd/atlas@latest

# The atlas-provider-gorm will be used by Atlas automatically
```

#### atlas.hcl Configuration

Create `atlas.hcl` in your project root:

```hcl
# Load GORM schema using external provider
data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/models",      # Path to your GORM models
    "--dialect", "postgres",             # Database dialect
  ]
}

# Environment configuration for Atlas
env "gorm" {
  # Source schema from GORM models
  src = data.external_schema.gorm.url

  # DEV DATABASE - Must be CLEAN - Only used for diff/comparison
  # Atlas uses this to calculate differences between your models and actual schema
  dev = "postgresql://postgres:postgres@localhost:5432/myapp_dev?sslmode=disable"

  # PRODUCTION DATABASE - Your actual database where migrations will be applied
  url = "postgresql://postgres:postgres@localhost:5432/myapp_db?sslmode=disable"

  # Migration directory
  migration {
    dir = "file://migrations"
  }

  # Format settings for generated SQL
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"  # Indent with 2 spaces
    }
  }
}
```

#### Understanding Atlas Databases

**DEV Database (`dev`):**
- Must be a **clean, empty database**
- Used only for schema comparison
- Atlas uses it as a temporary workspace
- Never manually modify this database
- Can be dropped and recreated anytime

**Production Database (`url`):**
- Your actual application database
- Where migrations are applied
- Contains real data
- Should be backed up regularly

#### Atlas Workflow

```bash
# 1. Ensure the dev database is completely clean (critical for accurate diff)
psql -U postgres -c "DROP DATABASE IF EXISTS myapp_dev WITH (FORCE);"
psql -U postgres -c "CREATE DATABASE myapp_dev;"

# 2. Generate a new migration from GORM model changes
atlas migrate diff add_user_profile --env gorm
# → Creates: migrations/20251121103022_add_user_profile.sql

# 3. Review & edit the generated migration (add Down section if needed)
code migrations/20251121103022_add_user_profile.sql   # or vim, nano, etc.

# 4. RECALCULATE CHECKSUM AFTER EDITING THE FILE (MANDATORY!)
atlas migrate hash --env gorm

# 5. Apply all pending migrations to the real database
atlas migrate apply --env gorm

# 6. Apply only a specific number of migrations
atlas migrate apply 1 --env gorm   # apply only the latest one
atlas migrate apply 3 --env gorm   # apply the latest 3

# 7. Check current migration status
atlas migrate status --env gorm

# 8. Show detailed migration history
atlas migrate status --env gorm --verbose

# 9. Lint – detect dangerous migrations (data loss, breaking changes…)
atlas migrate lint --env gorm
atlas migrate lint --env gorm --latest 3   # lint only the last 3 migrations

# 10. Preview SQL that will be executed (dry-run) – 100% safe
atlas migrate apply --env gorm --dry-run

# 11. Revert / Downgrade (Atlas’s superpower)
atlas migrate set --env gorm -1          # rollback last migration (most used)
atlas migrate set --env gorm -3          # rollback last 3 migrations
atlas migrate set --env gorm 20251121100000   # go exactly to this version
atlas migrate set --env gorm 0           # wipe everything (dev/staging only)

# 12. Validate entire migration directory (checksum + syntax)
atlas migrate validate --env gorm

# 13. Create an empty migration for manual SQL / data migration
atlas migrate new seed_initial_data

# 14. Push migrations to Atlas Cloud (perfect for CI/CD + approval flow)
atlas migrate push my-project-name --env gorm
# Then apply from Cloud without cloning repo:
atlas migrate apply --url atlas://my-project-name
```

#### Example: Complete Migration Flow

**Step 1: Create GORM Models**

```go
// internal/models/product.go
package models

import (
    "gorm.io/gorm"
    "time"
)

type Product struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Name        string         `gorm:"not null;index" json:"name"`
    Description string         `json:"description"`
    Price       float64        `gorm:"not null" json:"price"`
    Stock       int            `gorm:"default:0" json:"stock"`
    CategoryID  uint           `gorm:"not null" json:"category_id"`
    Category    Category       `gorm:"foreignKey:CategoryID" json:"category"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Category struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Name      string         `gorm:"uniqueIndex;not null" json:"name"`
    Products  []Product      `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**Step 2: Generate Migration**

```bash
# Generate migration from models
atlas migrate diff create_products_and_categories --env gorm

# This creates: migrations/20240115120000_create_products_and_categories.sql
```

**Step 3: Generated SQL Example**

```sql
-- migrations/20240115120000_create_products_and_categories.sql
-- Create "categories" table
CREATE TABLE "categories" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);

-- Create index "idx_categories_deleted_at" to table: "categories"
CREATE INDEX "idx_categories_deleted_at" ON "categories" ("deleted_at");

-- Create index "idx_categories_name" to table: "categories"
CREATE UNIQUE INDEX "idx_categories_name" ON "categories" ("name");

-- Create "products" table
CREATE TABLE "products" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "price" numeric NOT NULL,
  "stock" bigint NULL DEFAULT 0,
  "category_id" bigint NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_category" FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Create index "idx_products_deleted_at" to table: "products"
CREATE INDEX "idx_products_deleted_at" ON "products" ("deleted_at");

-- Create index "idx_products_name" to table: "products"
CREATE INDEX "idx_products_name" ON "products" ("name");
```

**Step 4: Apply Migration**

```bash
# Apply to production
atlas migrate apply --env gorm

# Output:
# Migrating to version 20240115120000 (1 migrations in total):
#   -- migrating version 20240115120000
#     -> CREATE TABLE "categories" ...
#     -> CREATE TABLE "products" ...
#   -- ok (25.5ms)
```

#### Atlas vs golang-migrate Comparison

| Feature | Atlas | golang-migrate |
|---------|-------|----------------|
| **Auto-generate from models** | ✅ Yes | ❌ No (manual SQL) |
| **Schema awareness** | ✅ Yes | ❌ No |
| **Diff calculation** | ✅ Automatic | ❌ Manual |
| **Rollback** | ✅ Built-in | ✅ Down migrations |
| **Migration validation** | ✅ Advanced | ⚠️ Basic |
| **Learning curve** | ⚠️ Medium | ✅ Simple |
| **Flexibility** | ⚠️ Less (auto-generated) | ✅ Full control |

#### Common Atlas Commands Reference

```bash
# Schema inspection
atlas schema inspect --env gorm --url "postgres://..."

# Schema comparison
atlas schema diff \
  --from "postgres://user:pass@localhost/db1" \
  --to "postgres://user:pass@localhost/db2"

# Apply schema directly (without migrations)
atlas schema apply --env gorm --auto-approve

# Create a new migration
atlas migrate diff add_email_to_users --env gorm

# Apply all pending migrations
atlas migrate apply --env gorm

# Apply specific version
atlas migrate apply --env gorm --to-version 20240115120000

# Rollback last migration
atlas migrate down --env gorm

# Check current database version
atlas migrate status --env gorm

# Validate migration files
atlas migrate validate --env gorm

# Lint migrations for potential issues
atlas migrate lint --env gorm

# Hash migration files (for integrity)
atlas migrate hash --env gorm
```

#### Best Practices for Atlas

1. **Always Keep Dev DB Clean**
   ```bash
   # Reset dev database before generating migrations
   psql -U postgres -c "DROP DATABASE myapp_dev;"
   psql -U postgres -c "CREATE DATABASE myapp_dev;"
   ```

2. **Review Generated Migrations**
   - Always review SQL before applying
   - Atlas generates safe migrations, but verify logic

3. **Use Version Control**
   ```bash
   # Commit migration files
   git add migrations/
   git commit -m "feat: add products and categories tables"
   ```

4. **Test Migrations**
   ```bash
   # Test on staging first
   atlas migrate apply --env staging
   
   # Then production
   atlas migrate apply --env gorm
   ```

5. **Backup Before Migration**
   ```bash
   # PostgreSQL backup
   pg_dump -U postgres myapp_db > backup_$(date +%Y%m%d_%H%M%S).sql
   
   # Apply migration
   atlas migrate apply --env gorm
   ```

#### Troubleshooting Atlas

**Problem: Dev database not clean**
```bash
Error: dev database is not clean

Solution:
psql -U postgres -c "DROP DATABASE myapp_dev; CREATE DATABASE myapp_dev;"
```

**Problem: Migration conflict**
```bash
Error: migration version conflict

Solution:
atlas migrate hash --env gorm  # Rehash migrations
```

**Problem: Cannot connect to database**
```bash
Error: failed to connect to database

Solution:
# Check connection string in atlas.hcl
# Verify database is running
psql -U postgres -d myapp_db -c "SELECT 1"
```

## 📝 Models

### internal/models/user.go

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID           uint           `gorm:"primaryKey" json:"id"`
    Email        string         `gorm:"uniqueIndex;not null" json:"email"`
    Username     string         `gorm:"uniqueIndex;not null" json:"username"`
    PasswordHash string         `gorm:"not null" json:"-"`
    FullName     string         `json:"full_name"`
    IsActive     bool           `gorm:"default:true" json:"is_active"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateUserRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Username string `json:"username" binding:"required,min=3,max=50"`
    Password string `json:"password" binding:"required,min=6"`
    FullName string `json:"full_name"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type UserResponse struct {
    ID        uint      `json:"id"`
    Email     string    `json:"email"`
    Username  string    `json:"username"`
    FullName  string    `json:"full_name"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
}
```

## 🔧 Repository Layer

### internal/repositories/user_repository.go

```go
package repositories

import (
    "github.com/template/go-backend-gin-orm/internal/models"
    "gorm.io/gorm"
)

type UserRepository interface {
    Create(user *models.User) error
    FindByID(id uint) (*models.User, error)
    FindByEmail(email string) (*models.User, error)
    FindByUsername(username string) (*models.User, error)
    Update(user *models.User) error
    Delete(id uint) error
    List(page, pageSize int) ([]models.User, int64, error)
}

type userRepository struct {
    // Note: We don't store db here because transactions need to pass db
}

func NewUserRepository() UserRepository {
    return &userRepository{}
}

func (r *userRepository) Create(user *models.User) error {
    return user.DB.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
    var user models.User
    err := user.DB.First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := user.DB.Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
    var user models.User
    err := user.DB.Where("username = ?", username).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
    return user.DB.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
    return user.DB.Delete(&models.User{}, id).Error
}

func (r *userRepository) List(page, pageSize int) ([]models.User, int64, error) {
    var users []models.User
    var total int64
    
    offset := (page - 1) * pageSize
    
    if err := db.Model(&models.User{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    err := db.Offset(offset).Limit(pageSize).Find(&users).Error
    return users, total, err
}
```

**Note**: We don't pass `db` to repository constructor because transactions need to pass `db` dynamically to each method.

## 💼 Service Layer

### internal/services/user_service.go

```go
package services

import (
    "errors"
    "github.com/template/go-backend-gin-orm/internal/models"
    "github.com/template/go-backend-gin-orm/internal/repositories"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type UserService interface {
    Register(req *models.CreateUserRequest) (*models.UserResponse, error)
    Login(req *models.LoginRequest) (string, error)
    GetUserByID(id uint) (*models.UserResponse, error)
    UpdateUser(id uint, updates map[string]interface{}) error
    DeleteUser(id uint) error
}

type userService struct {
    repo repositories.UserRepository
    db   *gorm.DB
}

func NewUserService(repo repositories.UserRepository, db *gorm.DB) UserService {
    return &userService{
        repo: repo,
        db:   db,
    }
}

func (s *userService) Register(req *models.CreateUserRequest) (*models.UserResponse, error) {
    // Check if email exists
    if _, err := s.repo.FindByEmail(req.Email); err == nil {
        return nil, errors.New("email already exists")
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    // Create user
    user := &models.User{
        Email:        req.Email,
        Username:     req.Username,
        PasswordHash: string(hashedPassword),
        FullName:     req.FullName,
        IsActive:     true,
    }
    
    if err := s.repo.Create(user); err != nil {
        return nil, err
    }
    
    return s.userToResponse(user), nil
}

func (s *userService) Login(req *models.LoginRequest) (string, error) {
    user, err := s.repo.FindByEmail(req.Email)
    if err != nil {
        return "", errors.New("invalid credentials")
    }
    
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return "", errors.New("invalid credentials")
    }
    
    // Generate JWT token (implement this)
    token, err := generateJWTToken(user)
    if err != nil {
        return "", err
    }
    
    return token, nil
}

func (s *userService) userToResponse(user *models.User) *models.UserResponse {
    return &models.UserResponse{
        ID:        user.ID,
        Email:     user.Email,
        Username:  user.Username,
        FullName:  user.FullName,
        IsActive:  user.IsActive,
        CreatedAt: user.CreatedAt,
    }
}
```

## 🎯 Handlers

### internal/handlers/user_handler.go

```go
package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/template/go-backend-gin-orm/internal/models"
    "github.com/template/go-backend-gin-orm/internal/services"
)

type UserHandler struct {
    service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
    return &UserHandler{service: service}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User registration data"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
    var req models.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user, err := h.service.Register(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    token, err := h.service.Login(&req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get user details by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }
    
    user, err := h.service.GetUserByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    
    c.JSON(http.StatusOK, user)
}
```

## 🛣️ Routes

### internal/routes/routes.go

```go
package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/template/go-backend-gin-orm/internal/handlers"
    "github.com/template/go-backend-gin-orm/internal/middleware"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine, userHandler *handlers.UserHandler) {
    // Swagger documentation
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // API v1 routes
    v1 := r.Group("/api/v1")
    {
        // Public routes
        auth := v1.Group("/auth")
        {
            auth.POST("/login", userHandler.Login)
            auth.POST("/register", userHandler.Register)
        }
        
        // Protected routes
        users := v1.Group("/users")
        users.Use(middleware.AuthMiddleware())
        {
            users.GET("/:id", userHandler.GetUser)
            users.PUT("/:id", userHandler.UpdateUser)
            users.DELETE("/:id", userHandler.DeleteUser)
            users.GET("", userHandler.ListUsers)
        }
    }
}
```

## 🚀 Main Application

### cmd/main.go

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/gin-gonic/gin"
    "github.com/template/go-backend-gin-orm/config"
    "github.com/template/go-backend-gin-orm/internal/handlers"
    "github.com/template/go-backend-gin-orm/internal/repositories"
    "github.com/template/go-backend-gin-orm/internal/routes"
    "github.com/template/go-backend-gin-orm/internal/services"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    
    _ "github.com/template/go-backend-gin-orm/docs"
)

// @title Go Backend API
// @version 1.0
// @description API documentation for Go Backend with Gin and GORM
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:5050
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // Connect to database
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        cfg.Database.Host,
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.DBName,
        cfg.Database.Port,
        cfg.Database.SSLMode,
    )
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    // Initialize layers
    userRepo := repositories.NewUserRepository()
    userService := services.NewUserService(userRepo, db)
    userHandler := handlers.NewUserHandler(userService)
    
    // Setup Gin
    gin.SetMode(cfg.Server.Mode)
    r := gin.Default()
    
    // Setup routes
    routes.SetupRoutes(r, userHandler)
    
    // Start server
    addr := ":" + cfg.Server.Port
    log.Printf("Server starting on %s", addr)
    if err := r.Run(addr); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

## 📚 Swagger Documentation

### Generate Swagger Docs

```bash
# Initialize swagger (run from project root)
swag init -g cmd/main.go

# This generates:
# - docs/docs.go
# - docs/swagger.json
# - docs/swagger.yaml
```

### Access Swagger UI

```
http://localhost:5050/swagger/index.html
```

### Swagger Annotations Examples

```go
// @Summary Short description
// @Description Detailed description
// @Tags tag-name
// @Accept json
// @Produce json
// @Param paramName path/query/body type required "description"
// @Success 200 {object} ResponseType
// @Failure 400 {object} ErrorType
// @Security BearerAuth
// @Router /path [method]
```

## 🏃 Running the Application

### Development Mode

```bash
# Simple run
go run cmd/main.go

# With hot reload (using nodemon)
nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run cmd/main.go

# Or using air (alternative)
go install github.com/cosmtrek/air@latest
air
```

### Production Build

```bash
# Build binary
go build -o bin/app cmd/main.go

# Run binary
./bin/app
```

## 🐳 Docker Support

### Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

EXPOSE 5050

CMD ["./main"]
```

### docker-compose.yml

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: myapp_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  app:
    build: .
    ports:
      - "5050:5050"
    depends_on:
      - postgres
    environment:
      DATABASE_HOST: postgres
      DATABASE_PORT: 5432
      DATABASE_USER: postgres
      DATABASE_PASSWORD: postgres
      DATABASE_NAME: myapp_db

volumes:
  postgres_data:
```

## 📋 Common Commands Cheat Sheet

```bash
# Project Setup
go mod init github.com/yourname/project
go mod tidy

# Run Application
go run cmd/main.go

# Build
go build -o bin/app cmd/main.go

# Migrations
migrate create -ext sql -dir migrations -seq migration_name
migrate -database $DATABASE_URL -path migrations up
migrate -database $DATABASE_URL -path migrations down 1

# Swagger
swag init -g cmd/main.go

# Testing
go test ./...
go test -v ./internal/services

# Format Code
go fmt ./...
gofmt -s -w .

# Linting
golangci-lint run
```

## 🔐 Environment Variables

### .env file

```env
SERVER_PORT=5050
SERVER_MODE=debug

DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=myapp_db
DATABASE_SSLMODE=disable

JWT_SECRET=your-super-secret-key-change-in-production
JWT_EXPIRATION=24h
```

## 📝 Best Practices

1. **Transaction Management**: Pass `db` to repository methods for proper transaction support
2. **Error Handling**: Always handle errors appropriately and return meaningful messages
3. **Validation**: Use Gin's binding validation for request validation
4. **Security**: 
   - Hash passwords with bcrypt
   - Use JWT for authentication
   - Implement rate limiting
   - Validate all inputs
5. **Documentation**: Keep Swagger docs updated with code changes
6. **Migrations**: Always create both up and down migrations
7. **Testing**: Write unit tests for services and integration tests for handlers

## 🎓 Additional Resources

- [Gin Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [Swagger Documentation](https://swagger.io/docs/)
- [Atlas Documentation](https://atlasgo.io/docs)

---
