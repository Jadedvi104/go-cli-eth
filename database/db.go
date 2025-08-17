package database

import (
	"fmt"
	"log"
	"os"

	"go-cli-eth/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB(connectionString string) error {
	var err error

	// If no connection string provided, try to get from environment
	if connectionString == "" {
		connectionString = os.Getenv("DATABASE_URL")
		if connectionString == "" {
			return fmt.Errorf("database connection string not provided")
		}
	}

	// Configure GORM logger
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to PostgreSQL
	DB, err = gorm.Open(postgres.Open(connectionString), config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	err = DB.AutoMigrate(&models.NFT{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
