package main

import (
	"conceptual-lan/internals"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// utils.DownloadOUI()
	//Middlewares
	// go utils.BroadcastListener()
	// go utils.BroadcastSender()

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	r := mux.NewRouter()
	r.HandleFunc("/get-my-network-info", internals.GetMyNetworkInfo).Methods("GET")

	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
