package test

import (
	"UserService/dbclient"
	"UserService/model"
	"UserService/service"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/dbtest"
)

var Server dbtest.DBServer
var MySession *mgo.Session

// TestMain wraps all tests with the needed initialized mock DB and fixtures
func TestMain(m *testing.M) {
	// The tempdir is created so MongoDB has a location to store its files.
	// Contents are wiped once the server stops
	tempDir, _ := ioutil.TempDir("", "testing")
	Server.SetPath(tempDir)

	// My main session var is now set to the temporary MongoDB instance
	MySession = Server.Session()

	// Make sure to insert my fixtures
	insertFixtures()

	// Run the test suite
	retCode := m.Run()

	// Close the session and wipe the database to ensure we're in a consistent state
	MySession.Close()
	Server.Wipe()

	// Stop shuts down the temporary server and removes data on disk.
	Server.Stop()

	// call with result of m.Run()
	os.Exit(retCode)
}

var users = map[string]model.User{
	"test": model.User{Id: "test", UserName: "user", Password: "password"},
}

// insertFixtures just inserts all users (and other types) I've defined above.
func insertFixtures() {
	for _, user := range users {
		if err := MySession.DB("store").C("users").Insert(user); err != nil {
			log.Println(err)
		}
	}
}

// TestGetAllUsers tests that the users returned from the database are what we expect. Nothing more, nothing less
func TestGetAllUsers(t *testing.T) {
	db := dbclient.MongoClient{
		Session: MySession,
	}

	service.MyEnv = service.Env{&db}

	users, _ := db.GetAllUsers()

	expectedUsers := make([]model.User, 1)
	expectedUser := model.User{
		Id: "test", UserName: "user", Password: "password",
	}
	expectedUsers[0] = expectedUser

	if !testEq(users, expectedUsers) {
		t.Errorf("Mongoclient returned unexpected result: got %v want %v",
			users, expectedUsers)
	}

}

// TestUserExists tests that for an existing user, MongoClient returns true
// For a username that doesn't appear in the database, we return false
func TestUserExists(t *testing.T) {
	db := dbclient.MongoClient{
		Session: MySession,
	}

	service.MyEnv = service.Env{&db}

	//-- test user that exists
	existingUser := model.User{
		Id: "test", UserName: "user", Password: "password",
	}
	exists, _ := db.UserExists(existingUser)

	if !exists {
		t.Errorf("Mongoclient returned unexpected result: got %v want %v",
			!exists, exists)
	}

	//-- test user that doesn't exist
	nonExistentUser := model.User{
		Id: "test", UserName: "blank", Password: "password",
	}
	exists, _ = db.UserExists(nonExistentUser)

	if exists {
		t.Errorf("Mongoclient returned unexpected result: got %v want %v",
			exists, !exists)
	}

	//-- TODO: test database failure is handled correctly. This requires mocking out the DB

}

func TestSaveUser(t *testing.T) {
	db := dbclient.MongoClient{
		Session: MySession,
	}

	service.MyEnv = service.Env{&db}

	newUser := model.User{
		UserName: "newuser", Password: "newpassword",
	}

	//-- check that the new user doesn't exist
	exists, _ := db.UserExists(newUser)

	if exists {
		t.Errorf("Unexpected error: user we're trying to save already exists!")
	}

	//-- save the user
	err := db.SaveUser(newUser)

	if err != nil {
		t.Errorf("Unexpected error while saving user")
	}

	//-- check that the new user now does exist
	exists, _ = db.UserExists(newUser)

	if !exists {
		t.Errorf("Unexpected error: user we're trying to save not found!")
	}

	//-- Finally, we need to delete that new user to ensure the database is consistent for the other tests

	error := MySession.DB("store").C("users").Remove(newUser)
	if error != nil {
		t.Errorf("Unexpected error: unable to remove saved user")
	}

	//-- check that the new user now does NOT exist
	exists, _ = db.UserExists(newUser)

	if exists {
		t.Errorf("Unexpected error: database not in a consistent state!")
	}

	//-- TODO: test database failure is handled correctly. This requires mocking out the DB

}

// TestAuthenticateUser tests that we return true for valid username and password combinations
// NOTE: in a production system you shouldn't store passwords in cleartext like we are
func TestAuthenticateUser(t *testing.T) {
	db := dbclient.MongoClient{
		Session: MySession,
	}

	service.MyEnv = service.Env{&db}

	//-- test valid user
	validUser := model.User{
		Id: "test", UserName: "user", Password: "password",
	}

	isValid, _ := db.AuthenticateUser(validUser)

	if !isValid {
		t.Errorf("Mongoclient returned unexpected result: got %v want %v",
			isValid, !isValid)
	}

	//-- test user with incorrect username
	invalidUser := model.User{
		Id: "test", UserName: "unknown", Password: "password",
	}

	isValid, _ = db.AuthenticateUser(invalidUser)

	if isValid {
		t.Errorf("Mongoclient returned unexpected result: got %v want %v",
			!isValid, isValid)
	}

	//-- test user with incorrect password
	invalidPassword := model.User{
		Id: "test", UserName: "user", Password: "spassword",
	}

	isValid, _ = db.AuthenticateUser(invalidPassword)

	if isValid {
		t.Errorf("Mongoclient returned unexpected result: got %v want %v",
			!isValid, isValid)
	}

	//-- TODO: Test DB failure. This requires mocking out the db connection.

}

func testEq(a, b []model.User) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
