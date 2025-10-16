package handlers

import (
	"data_mapping/models"
	"data_mapping/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateClient(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateClientRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		// Validate the request
		if err := utils.ValidateStruct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return
		}

		client := models.Client{
			Name: req.Name,
		}

		if result := db.Create(&client); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create client",
				"details": result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data":    client,
		})
	}
}

func ListClients(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var clients []models.Client
		if result := db.Find(&clients); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, clients)
	}
}

func DeleteClient(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}
		if result := db.Where("client_id = ?", id).Delete(&models.MappingRule{}); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		if result := db.Delete(&models.Client{}, id); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
