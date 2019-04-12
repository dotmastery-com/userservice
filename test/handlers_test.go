package test

import (
	"UserService/mock"
	"UserService/model"
	"UserService/service"
	"bytes"
	"encoding/json"

	"errors"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.HealthCheckHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `"{alive: true}"`

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetAllUsersWithMultipleUsersReturnsCorrectResult(t *testing.T) {

	user1 := model.User{
		Id:       "TestUser",
		UserName: "UserName",
		Password: "Password",
	}

	user2 := model.User{
		Id:       "TestUser2",
		UserName: "UserName2",
		Password: "Password2",
	}

	users := make([]model.User, 2)
	users[0] = user1
	users[1] = user2

	expected := `[{"id":"TestUser","username":"UserName","password":"Password"},{"id":"TestUser2","username":"UserName2","password":"Password2"}]`

	testGetAllUsers(t, users, expected, nil, 200)

}

func TestGetAllUsersWithSingleUserReturnsCorrectResult(t *testing.T) {

	user := model.User{
		Id:       "TestUser",
		UserName: "UserName",
		Password: "Password",
	}

	users := make([]model.User, 1)
	users[0] = user

	expected := `[{"id":"TestUser","username":"UserName","password":"Password"}]`

	testGetAllUsers(t, users, expected, nil, 200)

}

func TestGetAllUsersWithNoUsersReturnsCorrectResult(t *testing.T) {

	expected := "null"
	testGetAllUsers(t, nil, expected, nil, 200)

}

func TestGetAllUsersWithDBErrorFailsCorrectly(t *testing.T) {

	expected := "{message: \"Database error\"}"
	err := errors.New("Unable to access database")
	testGetAllUsers(t, nil, expected, err, 500)

}

func testGetAllUsers(t *testing.T, input []model.User, expected string, err error, errcode int) {
	mockDatabase := new(mock.DBClientMock)
	service.MyEnv.DB = mockDatabase

	mockDatabase.On("GetAllUsers").Return(input, err)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/getAllUsers", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.GetAllUsers)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != errcode {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, errcode)
	}

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAuthUsersTableTest(t *testing.T) {

	var tests = []struct {
		input    bool
		err      error
		expected string
		errcode  int
	}{
		//-- user exists in database
		{true, nil, "", 200},

		//-- user doesn't exist in database
		{false, nil, "{message: \"Invalid username or password\"}", 401},

		//-- database connection error
		{false, errors.New("Database issue"), "{message: \"Database error\"}", 500},
	}

	for _, test := range tests {
		testAuthUser(t, test.input, test.err, test.expected, test.errcode)
	}
}

func testAuthUser(t *testing.T, input bool, err error, expected string, errcode int) {
	mockDatabase := new(mock.DBClientMock)
	service.MyEnv.DB = mockDatabase

	mockDatabase.On("AuthenticateUser").Return(input, err)

	user := &model.User{
		Id:       "TestUser",
		UserName: "UserName",
		Password: "Password",
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(user)
	req, err := http.NewRequest("POST", "/authUser", buf)

	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.AuthUser)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != errcode {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, errcode)
	}

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
