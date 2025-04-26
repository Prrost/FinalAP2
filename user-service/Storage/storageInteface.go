package Storage

import "user-service/domain"

type Storage interface {
	CreateUser(user domain.User) (domain.User, error)
	CreateUserAdmin(user domain.User) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	GetUserByID(id int) (domain.User, error)
}
