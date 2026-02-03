package auth_domain

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
