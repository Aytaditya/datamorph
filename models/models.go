package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Client struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name" validate:"required,min=1,max=100"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MappingRule struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ClientID        uint           `gorm:"not null" json:"client_id"`
	Client          Client         `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-" validate:"-"`
	SourcePath      JSONStringList `gorm:"type:jsonb;not null" json:"source_path" validate:"required,min=1"`
	DestinationPath JSONStringList `gorm:"type:jsonb;not null" json:"destination_path" validate:"required,min=1"`
	TransformType   string         `gorm:"not null" json:"transform_type" validate:"required,oneof=copy toString mapGender toBool formatDate toUpperCase toLowerCase capitalize expression"`
	TransformLogic  string         `gorm:"type:text" json:"transform_logic"`
	Required        bool           `gorm:"default:false" json:"required"`
	DefaultValue    string         `gorm:"type:text" json:"default_value"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type JSONStringList []string

func (j *JSONStringList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value")
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONStringList) Value() (driver.Value, error) {
	return json.Marshal(j)
}
