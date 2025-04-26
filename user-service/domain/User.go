package domain

type User struct {
	ID       uint
	Email    string `validate:"required,email"`
	Password string
	IsAdmin  bool
}
