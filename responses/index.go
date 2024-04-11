package responses

type UserCreatedResponse struct {
	Message string `json:"message"`
	Status  string  `json:"status"`
}

type UserLoginResponse struct {
	Message string `json:"message"`
	Status  string  `json:"status"`
}

type RequestDashboardResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Role    string  `json:"role"`
}
