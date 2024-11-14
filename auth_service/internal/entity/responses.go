package entity

type AuthResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
