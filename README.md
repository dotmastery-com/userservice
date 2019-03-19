# UserService

The UserService microservice provides three endpoints relating to users:

- localhost:7000/auth: which takes a username and password as parameter and returns either a HTTP 200 to indicate success or another error code with associated message to indicate failure. This could be an application error such as user not found or a system error such as database not available

- localhost:7000/getAllUsers: This returns a JSON array of all users, with their usernames and passwords in cleartext.
This is useful for debugging but should be disabled in a production system

- localhost:7000/register: This takes a username and password as a parameter and creates a new user record in the database.


## Prerequisites

1. Install Golang as per the install instructions here: https://golang.org/doc/install

2. Install mongo as per https://docs.mongodb.com/manual/administration/install-community/
Start a mongo daemon by running _mongod_ in a terminal

## Running the code

1. Clone the project and make sure its part of the $GOPATH
2. From the command line run go run main.go

You should see a confirmation message that the service is up and running

