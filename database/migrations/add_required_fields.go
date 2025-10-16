package migrations

import (
	"gorm.io/gorm"
)

// AddRequiredFieldsToMappingRules adds the 'required' and 'default_value' columns to the mapping_rules table
func AddRequiredFieldsToMappingRules(db *gorm.DB) error {
	// Add required field with default value false
	if err := db.Exec("ALTER TABLE mapping_rules ADD COLUMN IF NOT EXISTS required BOOLEAN DEFAULT false").Error; err != nil {
		return err
	}

	// Add default_value field
	if err := db.Exec("ALTER TABLE mapping_rules ADD COLUMN IF NOT EXISTS default_value TEXT").Error; err != nil {
		return err
	}

	return nil
}
