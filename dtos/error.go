package dtos

type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
