package responses

type SuccessResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    Data   `json:"data,omitempty"`
}

type Data struct {
	AccessToken   string            `json:"accessToken,omitempty"`
	UserID        string            `json:"userId,omitempty"`
	FirstName     string            `json:"firstName,omitempty"`
	LastName      string            `json:"lastName,omitempty"`
	Email         string            `json:"email,omitempty"`
	Phone         string            `json:"phone,omitempty"`
	OrgID         string            `json:"orgId,omitempty"`
	Name          string            `json:"name,omitempty"`
	Description   string            `json:"description,omitempty"`
	User          *UserRes          `json:"user,omitempty"`
	Organisations []OrganisationRes `json:"organisations,omitempty"`
}

type UserRes struct {
	UserID        string            `json:"userId,omitempty"`
	FirstName     string            `json:"firstName,omitempty"`
	LastName      string            `json:"lastName,omitempty"`
	Email         string            `json:"email,omitempty"`
	Phone         string            `json:"phone,omitempty"`
	Organisations []OrganisationRes `json:"organisations,omitempty"`
}

type OrganisationRes struct {
	OrgID       string `json:"orgId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
