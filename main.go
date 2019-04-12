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
	db := dbclient.MongoClient{}
	service.MyEnv = service.Env{&db}
	service.MyEnv.DB.Connect()
}
