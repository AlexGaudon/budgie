package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alexgaudon/budgie/config"
	"github.com/alexgaudon/budgie/server"
	"github.com/alexgaudon/budgie/storage"
)

func main() {
	config.LoadConfig()
	r := server.SetupServer()

	err := storage.SetupDatabase()

	if err != nil {
		log.Fatal("Error setting up database: ", err)
	}

	if err := storage.DB.Init(); err != nil {
		log.Fatal("Error initializing database: ", err)
	}

	port := os.Getenv("SERVER_PORT")
	log.Println("Starting server on :" + port)

	if port == "" {
		port = "3000"
	}

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
