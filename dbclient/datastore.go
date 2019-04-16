package dbclient

import "UserService/model"

type Datastore interface {
	Connect()
	AuthenticateUser(user model.User) (bool, error)
	UserExists(user model.User) (bool, error)
	SaveUser(user model.User) error
	GetAllUsers() ([]model.User, error)
}
