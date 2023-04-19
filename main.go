package main

import (
	"log"
	"net/http"

	"github.com/alexgaudon/budgie/server"
	"github.com/alexgaudon/budgie/storage"
)

func main() {
	r := server.SetupServer()

	err := storage.SetupDatabase()

	if err != nil {
		log.Fatal("Error setting up database: ", err)
	}

	if err := storage.DB.Init(); err != nil {
		log.Fatal("Error initializing database: ", err)
	}

	log.Println("Starting server on :3000")

	err = http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
