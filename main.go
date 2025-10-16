package main

import (
	"data_mapping/config"
	"data_mapping/database"
	"data_mapping/handlers"
	"data_mapping/middleware"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	if config.AppConfig.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	
	router.Use(gin.Logger())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.SecurityMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggingMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "Data Mapping API",
			"version": "1.0.0",
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the Data Mapping API",
			"version": "1.0.0",
			"docs":    "/docs",
		})
	})

	router.POST("/login", handlers.LoginHandler())

	// Protected routes
	auth := router.Group("/")
	auth.Use(handlers.JWTAuthMiddleware())
	{
		// Client management
		auth.POST("/clients", handlers.CreateClient(database.DB))
		auth.GET("/clients", handlers.ListClients(database.DB))
		auth.DELETE("/clients/:id", handlers.DeleteClient(database.DB))
		auth.POST("/clients/:client_id/mappings", handlers.CreateMappings(database.DB))
		auth.GET("/clients/:client_id/mappings", handlers.GetMappings(database.DB))
		auth.DELETE("/mappings/:mapping_id", handlers.DeleteMappings(database.DB))

		auth.POST("/clients/:client_id/transform", handlers.UnifiedTransformHandler(database.DB))
	}

	serverAddr := ":" + config.AppConfig.ServerPort

	if config.AppConfig.DevelopmentMode {
		fmt.Printf("Starting server in DEVELOPMENT mode on http://localhost%s\n", serverAddr)
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	} else {
		fmt.Printf("Starting server on https://localhost%s\n", serverAddr)
		if err := router.RunTLS(serverAddr, config.AppConfig.CertFilePath, config.AppConfig.KeyFilePath); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}
