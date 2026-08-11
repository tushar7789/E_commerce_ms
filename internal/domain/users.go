package domain

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
