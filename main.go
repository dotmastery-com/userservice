package main

import (
	"UserService/dbclient"
	"UserService/service"
	"fmt"


)

var appName = "User Service"

func main() {
	fmt.Printf("Starting %v\n", appName)
	initDBClient()
	service.StartWebServer("7000")
}

func initDBClient() {
	service.DBClient = dbclient.MongoClient{}
	service.DBClient.Connect()
}


