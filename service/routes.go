package service

import "net/http"

// Defines a single route, e.g. a human readable name, HTTP method and the
// pattern the function that will execute when the route is called.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Defines the type Routes which is just an array (slice) of Route structs.
type Routes []Route

// Initialize our routes
var routes = Routes{

	Route{
		"AuthUser", // Name
		"POST",     // HTTP method
		"/auth",
		AuthUser, // Route pattern

	},
	Route{
		"RegisterUser", // Name
		"POST",     // HTTP method
		"/register",
		RegisterUser, // Route pattern

	},
	Route{
		"AllUsers", // Name
		"GET",     // HTTP method
		"/getAllUsers",
		GetAllUsers, // Route pattern

	},
}
