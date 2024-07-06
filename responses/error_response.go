package responses

type ErrorResponse struct {
	Status     string  `json:"status,omitempty"`
	Message    string  `json:"message,omitempty"`
	StatusCode int     `json:"statusCode,omitempty"`
	Errors     []Error `json:"errors,omitempty"`
}

type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
