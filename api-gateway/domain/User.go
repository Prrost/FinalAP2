package domain

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}
