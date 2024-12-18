package models

// SuccessResponse - структура для успешного ответа
type SuccessResponse struct {
	Message string `json:"message" example:"Event added successfully"`
}

// ErrorResponse - структура для ошибок
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid input"`
	Details string `json:"details,omitempty" example:"Error details here"`
}
