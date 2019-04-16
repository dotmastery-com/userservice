package service

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func StartWebServer(port string) {

	r := router()
	http.Handle("/", r)

	handler := cors.Default().Handler(r)

	log.Println("Starting Http service at " + port)
	err := http.ListenAndServe(":"+port, handler) // goroutine will block here

	if err != nil {
		log.Println("An error occured starting Http listener at port " + port)
		log.Println("Error: " + err.Error())
	}

}

func router() *mux.Router {

	// create an instance of the gorilla router
	router := mux.NewRouter().StrictSlash(true)

	// iterate over the routes we declared in routes.go an attach them to the router interface

	for _, route := range routes {

		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router

}
