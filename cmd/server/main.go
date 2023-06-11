package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexgaudon/budgie/config"
	"github.com/alexgaudon/budgie/server"
	"github.com/alexgaudon/budgie/storage"
)

func main() {
	config.LoadConfig()

	db, err := storage.ConnectDatabase("./migrations")

	if err != nil {
		panic(err)
	}

	err = db.Initialize()

	if err != nil {
		panic(err)
	}

	if err != nil {
		fmt.Println("Err", err.Error())
	}

	server := server.NewAPIServer(db)
	server.ConfigureServer()

	port := config.GetConfig().ServerPort

	log.Println("Starting web server on port " + port)

	err = http.ListenAndServe(":"+port, server.Router)

	if err != nil {
		panic(err)
	}
}
