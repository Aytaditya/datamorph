package handlers

import (
	"data_mapping/models"
	"data_mapping/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/antonmedv/expr"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateMappings(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, err := strconv.Atoi(c.Param("client_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid client ID",
			})
			return
		}

		var rules []models.MappingRule
		if err := c.ShouldBindJSON(&rules); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		// Set ClientID and validate each rule
		for i := range rules {
			rules[i].ClientID = uint(clientID)

			// Special validation for expression type
			if rules[i].TransformType == "expression" {
				if rules[i].TransformLogic == "" {
					c.JSON(http.StatusBadRequest, gin.H{
						"error":   "Validation failed for rule " + strconv.Itoa(i),
						"details": "TransformLogic is required when TransformType is 'expression'",
					})
					return
				}

				// Try to validate expression syntax
				if _, err := expr.Compile(rules[i].TransformLogic); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error":   "Invalid expression syntax in rule " + strconv.Itoa(i),
						"details": err.Error(),
					})
					return
				}
			}

			// Validate the rule after setting required fields using custom validation
			if err := utils.ValidateMappingRule(rules[i]); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Validation failed for rule " + strconv.Itoa(i),
					"details": err.Error(),
				})
				return
			}

			// Validate required fields have appropriate defaults
			if rules[i].Required && rules[i].DefaultValue == "" {
				log.Printf("Warning: Required field mapping without default value: %v -> %v",
					rules[i].SourcePath, rules[i].DestinationPath)
			}
		}

		if result := db.Create(&rules); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create mapping rules",
				"details": result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data":    rules,
		})
	}
}

func GetMappings(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.Param("client_id")
		var rules []models.MappingRule
		result := db.Where("client_id = ?", clientID).Find(&rules)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, rules)
	}
}

func DeleteMappings(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		mappingID := c.Param("mapping_id")
		result := db.Delete(&models.MappingRule{}, mappingID)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Mapping rule not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
