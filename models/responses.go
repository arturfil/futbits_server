package models

type JsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type TokenResponse struct {
	Token string `json:"token"`
	User  *User
}
