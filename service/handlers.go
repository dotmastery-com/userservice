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

var DBClient dbclient.MongoClient


func errorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}


/*
   Authenticates a user by checking whether a corresponding entry exists in the database

   Returns HTTP 200 if a user with the given username and password exists
   Returns HTTP 401 if no such user exists
   Returns
 */
func AuthUser(w http.ResponseWriter, r *http.Request) {

	var user model.User
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&user)
	if err != nil {
		errorWithJSON(w, "Invalid body", http.StatusBadRequest)
		return
	}

	found, err := DBClient.AuthenticateUser(user)
	if err != nil {
		errorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed to retrieve uer: ", err)
		return
	}
	
	if !found {
		errorWithJSON(w, "Invalid username or password", http.StatusUnauthorized)
		log.Println("Failed get all messages: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

/*
   Registers a new user in the database if the username doesn't already exist
 */
func RegisterUser(w http.ResponseWriter, r *http.Request)  {

	var user model.User
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&user)
	if err != nil {
		errorWithJSON(w, "Incorrect body", http.StatusBadRequest)
		return
	}

	found, err := DBClient.UserExists(user)
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
		err = DBClient.SaveUser(user);
		if err != nil {
			errorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Unable to add user to the database ", err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	users, err := DBClient.GetAllUsers()

	// If err, return a 404
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// If found, marshal into JSON, write headers and content
	data, _ := json.Marshal(users)
	fmt.Println(data)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
