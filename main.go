package main

import (
	"fmt"
	"statistical-analysis/config"
	"statistical-analysis/server"
	"statistical-analysis/service"
)

func main() {
	dbUrl := config.Configuration.MongoURL
	dbName := config.Configuration.Database
	port := config.Configuration.Port
	fmt.Println("Starting service")
	service.Init(dbUrl, dbName)
	fmt.Println("Started service")

	fmt.Println("Starting server")
	server.Start(port)
}
