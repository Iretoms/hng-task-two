package responses

type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type Data struct {
	AccessToken string `json:"accessToken"`
	User        User   `json:"user"`
}

type User struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

