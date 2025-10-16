package middleware

import (
	"data_mapping/database"
	"data_mapping/models"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		user, _ := c.Get("user")
		userStr, _ := user.(string)

		logEntry := models.Log{
			Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
			Method:      c.Request.Method,
			Path:        c.Request.URL.Path,
			ClientIP:    c.ClientIP(),
			StatusCode:  c.Writer.Status(),
			User:        userStr,
			QueryParams: c.Request.URL.RawQuery,
			Error:       c.Errors.ByType(gin.ErrorTypePrivate).String(),
		}

		if err := database.DB.Create(&logEntry).Error; err != nil {
			log.Println("Failed to save log to database:", err)
		} else {
			log.Println("Log saved to database.")
		}

		logJSON, _ := json.MarshalIndent(logEntry, "", "  ")
		fmt.Println(string(logJSON))
	}
}
