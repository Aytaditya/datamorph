package handlers

import (
	"data_mapping/models"
	"data_mapping/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UnifiedTransformHandler handles both standard and large payloads for transformation.
func UnifiedTransformHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.Param("client_id")
		var rules []models.MappingRule
		result := db.Where("client_id = ?", clientID).Find(&rules)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to load mapping rules",
				"details": result.Error.Error(),
			})
			return
		}

		if len(rules) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "No mapping rules found for this client",
			})
			return
		}

		// Limit payload size for security (e.g., 10MB)
		if c.Request.ContentLength > 10*1024*1024 {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "Payload too large. Max 10MB allowed.",
			})
			return
		}

		// Handle streaming for large payloads
		stream := c.GetHeader("X-Stream-Transform") == "true"
		if stream || (c.Request.ContentLength > 5*1024*1024) {
			c.Writer.Header().Set("Content-Type", "application/json")
			if err := utils.StreamTransformJSONWithRules(c.Request.Body, c.Writer, rules); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Streaming transformation failed",
					"details": err.Error(),
				})
			}
			return
		}

		// Standard transformation for smaller payloads
		var request models.TransformationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid JSON input",
				"details": err.Error(),
			})
			return
		}

		// Debug: Log the number of rules and input structure
		log.Printf("Transform Debug - Client ID: %s, Rules count: %d", clientID, len(rules))
		inputKeys := make([]string, 0, len(request.InputData))
		for k := range request.InputData {
			inputKeys = append(inputKeys, k)
		}
		log.Printf("Transform Debug - Input data keys: %v", inputKeys)
		for i, rule := range rules {
			log.Printf("Rule %d: %v -> %v (%s)", i, rule.SourcePath, rule.DestinationPath, rule.TransformType)
		}

		output, err := utils.Transform(request.InputData, rules)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Transformation failed",
				"details": err.Error(),
			})
			return
		}

		// Validate that all required fields are present
		var missingFields []string
		for _, rule := range rules {
			if rule.Required {
				// Check if the output has applicants array
				if applicants, ok := output["applicants"].([]interface{}); ok {
					// Check the first applicant (assuming all applicants have the same structure)
					if len(applicants) > 0 {
						if applicant, ok := applicants[0].(map[string]interface{}); ok {
							if _, exists := utils.GetNestedValue(applicant, rule.DestinationPath); !exists {
								path := strings.Join(rule.DestinationPath, ".")
								missingFields = append(missingFields, path)
							}
						}
					}
				} else {
					// Single object output
					if _, exists := utils.GetNestedValue(output, rule.DestinationPath); !exists {
						path := strings.Join(rule.DestinationPath, ".")
						missingFields = append(missingFields, path)
					}
				}
			}
		}

		// Remove duplicate entries from missingFields
		seen := make(map[string]bool)
		unique := make([]string, 0, len(missingFields))
		for _, field := range missingFields {
			if !seen[field] {
				seen[field] = true
				unique = append(unique, field)
			}
		}
		missingFields = unique

		response := gin.H{
			"success": true,
			"data":    output,
		}

		if len(missingFields) > 0 {
			response["warnings"] = gin.H{
				"missingRequiredFields": missingFields,
			}
		}

		c.JSON(http.StatusOK, response)
	}
}
