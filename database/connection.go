package database

import (
	"data_mapping/config"
	"data_mapping/database/migrations"
	"data_mapping/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	// Try to connect with retries
	var db *gorm.DB
	var err error
	
	maxRetries := 5
	retryDelay := time.Second * 3
	
	for i := 0; i < maxRetries; i++ {
		db, err = connectDB()
		if err == nil {
			break
		}
		
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			log.Printf("Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
		}
	}
	
	if err != nil {
		log.Fatal("Failed to connect to database after multiple attempts:", err)
	}
	
	log.Println("Successfully connected to database")
	DB = db
	
	// Run migrations
	log.Println("Running auto migrations...")
	err = DB.AutoMigrate(&models.Log{}, &models.Client{}, &models.MappingRule{})
	if err != nil {
		log.Printf("Warning: Failed to run auto migrations: %v", err)
	}

	// Run custom migrations
	if err := migrations.AddRequiredFieldsToMappingRules(DB); err != nil {
		log.Printf("Warning: Failed to run custom migrations: %v", err)
	}
	
	log.Println("Database initialization complete")
}

func connectDB() (*gorm.DB, error) {
	var dsn string
	
	// Print connection details for debugging (remove sensitive info in production)
	log.Printf("Connecting to database at %s:%s as %s", 
		config.AppConfig.DBHost, 
		config.AppConfig.DBPort, 
		config.AppConfig.DBUser)
	
	if config.AppConfig.DatabaseURL != "" {
		dsn = config.AppConfig.DatabaseURL
	} else {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
			config.AppConfig.DBHost,
			config.AppConfig.DBUser,
			config.AppConfig.DBPassword,
			config.AppConfig.DBName,
			config.AppConfig.DBPort,
		)
	}

	// Configure GORM with connection pooling and timeouts
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC() // Use UTC for all timestamps
		},
	})
}