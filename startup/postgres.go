package startup

import (
	"log"

	"github.com/template/go-backend-gin-orm/config"
	"github.com/template/go-backend-gin-orm/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDatabase struct {
	DB *gorm.DB
}

func NewPostgresDatabase() *PostgresDatabase {
	env := config.NewEnv(".env", true)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  env.DBUrl,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable SQL query logging
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		// return nil, err
	}
	log.Println("Connected to the database successfully!")
	return &PostgresDatabase{DB: db}
}

func (db *PostgresDatabase) MigrateDatabase() {
	err := db.DB.AutoMigrate(model.GetAllModels()...)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations completed successfully!")
}
