// main.go

package main

import (
	"Ozon/db"
	"Ozon/handler"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	var storageType string
	flag.StringVar(&storageType, "storage", "memory", "Specify the storage type: memory or postgres")
	flag.Parse()

	var storage db.Storage
	switch storageType {
	case "memory":
		storage = db.NewMemoryStorage()
	case "postgres":
		storage = db.NewPostgresStorage()
	default:
		log.Fatal("Invalid storage type specified")
	}

	handler := handler.NewHandler(storage)

	r := mux.NewRouter()
	r.HandleFunc("/shorten", handler.ShortenURL).Methods("POST")
	r.HandleFunc("/{shortLink}", handler.ExpandURL).Methods("GET")

	log.Println("Server is running...")
	log.Fatal(http.ListenAndServe(":8085", r))
}
