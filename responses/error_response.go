package responses

type ErrorResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
	StatusCode int `json:"statusCode"`
}