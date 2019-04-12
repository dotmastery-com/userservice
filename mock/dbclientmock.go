package mock

import (
	"UserService/model"

	"github.com/stretchr/testify/mock"
)

type DBClientMock struct {
	mock.Mock
}

// Our mocked GetAllUsers method
func (m *DBClientMock) GetAllUsers() ([]model.User, error) {

	//-- grab the parameters
	args := m.Called()

	//-- draw out the user slice to return
	results := args.Get(0).([]model.User)

	//-- and any errors
	err := args.Error(1)

	//-- return
	return results, err
}

func (m *DBClientMock) Connect() {

}

func (m *DBClientMock) AuthenticateUser(user model.User) (bool, error) {

	//-- grab the parameters
	args := m.Called()

	//-- draw out the user slice to return
	results := args.Bool(0)

	//-- and any errors
	err := args.Error(1)

	//-- return
	return results, err

}

func (m *DBClientMock) UserExists(user model.User) (bool, error) {
	return true, nil

}

func (m *DBClientMock) SaveUser(user model.User) error {
	return nil
}
