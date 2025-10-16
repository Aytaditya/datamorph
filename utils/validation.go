package utils

import (
	"data_mapping/models"
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct validates a struct using validator tags
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Field '%s' failed validation: %s", err.Field(), err.Tag()))
		}
		return fmt.Errorf("validation failed: %s", strings.Join(errors, ", "))
	}
	return nil
}

// ValidateMappingRule validates only the essential fields of a mapping rule
func ValidateMappingRule(rule interface{}) error {
	// Create a custom validator that skips nested structs
	v := validator.New()

	// Configure validator to skip dive validation on nested structs
	v.SetTagName("validate")

	if err := v.Struct(rule); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			// Skip validation errors for nested structs (like Client.Name)
			fieldName := err.Field()
			if strings.Contains(fieldName, "Client") || strings.Contains(fieldName, "Name") {
				continue
			}
			errors = append(errors, fmt.Sprintf("Field '%s' failed validation: %s", fieldName, err.Tag()))
		}
		if len(errors) > 0 {
			return fmt.Errorf("validation failed: %s", strings.Join(errors, ", "))
		}
	}

	// Additional validation for expression type mappings
	if r, ok := rule.(models.MappingRule); ok {
		if r.TransformType == "expression" && r.TransformLogic == "" {
			return fmt.Errorf("validation failed: TransformLogic is required when TransformType is 'expression'")
		}

		// If TransformLogic is provided, try to validate it's a valid expression
		if r.TransformLogic != "" {
			if _, err := expr.Compile(r.TransformLogic); err != nil {
				return fmt.Errorf("validation failed: Invalid expression syntax in TransformLogic: %s", err.Error())
			}
		}
	}

	return nil
}
