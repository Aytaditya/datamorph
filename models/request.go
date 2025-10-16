package models

type TransformationRequest struct {
	InputData map[string]interface{} `json:"input_data" binding:"required" validate:"required"`
}

type CreateClientRequest struct {
	Name string `json:"name" binding:"required" validate:"required,min=1,max=100"`
}
