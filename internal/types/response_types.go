package types

type SuccessResponse struct {
	Message string `json:"message" example:"OK"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
