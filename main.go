package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"statistical-analysis/server"
	"statistical-analysis/service"
)

var config map[string]interface{}

func init() {
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}
}

func main() {
	dbUrl, ok := config["mongoUrl"].(string)
	if !ok {
		panic(`Expected "mongoUrl" to be "string"`)
	}
	dbName, ok := config["dbName"].(string)
	if !ok {
		panic(`Expected "dbName" to be "string"`)
	}
	port, ok := config["statPort"].(string)
	if !ok {
		panic(`Expected "statPort" to be "string"`)
	}
	fmt.Println("Starting service")
	service.Init(dbUrl, dbName)
	fmt.Println("Started service")

	fmt.Println("Starting server")
	server.Start(port)
}
