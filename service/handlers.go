package service

import (
	"UserService/dbclient"
	"UserService/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Env struct {
	DB dbclient.Datastore
}

var MyEnv = Env{}

func errorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

/*
   Authenticates a user by checking whether a corresponding entry exists in the database

   Returns HTTP 200 if a user with the given username and password exists
   Returns HTTP 401 if no such user exists
   Returns HTTP 500 if internal error
*/
func AuthUser(w http.ResponseWriter, r *http.Request) {

	var user model.User
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&user)

	if err != nil {
		errorWithJSON(w, "Invalid body", http.StatusBadRequest)
		return
	}

	found, err := MyEnv.DB.AuthenticateUser(user)
	if err != nil {
		errorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed to retrieve user: ", err)
		return
	}

	if !found {
		errorWithJSON(w, "Invalid username or password", http.StatusUnauthorized)
		log.Println("Failed to auth users: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return
}

/*
   Registers a new user in the database if the username doesn't already exist

   Returns 200 if successful
   Returns 500 if database unavailbale
   Returns 403 if user already exists
*/
func RegisterUser(w http.ResponseWriter, r *http.Request) {

	var user model.User
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&user)
	if err != nil {
		errorWithJSON(w, "Incorrect body", http.StatusBadRequest)
		return
	}

	found, err := MyEnv.DB.UserExists(user)
	if err != nil {
		errorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed to look up user: ", err)
		return
	}

	if found {
		errorWithJSON(w, "User exists", http.StatusForbidden)
		return
	}

	if !found {
		err = MyEnv.DB.SaveUser(user)
		if err != nil {
			errorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Unable to add user to the database ", err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return

}

/*
	Returns a JSON array of all users, or a 500 error if unable to connect to the database
*/
func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	users, err := MyEnv.DB.GetAllUsers()

	// If err, return a 503, service unavailable with error message
	if err != nil {
		errorWithJSON(w, "Database error", http.StatusInternalServerError)
		return
	}

	// If found, marshal into JSON, write headers and content
	data, _ := json.Marshal(users)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// e.g. http.HandleFunc("/health-check", HealthCheckHandler)
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache and include them in the response.
	response := `{alive: true}`
	data, _ := json.Marshal(response)
	w.Write(data)
}
