package main

import (
	"conceptual-lan/internals"
	"conceptual-lan/utils"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	hub := internals.NewHub()
	go hub.Run() // Run the hub in the background

	r := mux.NewRouter()
	r.HandleFunc("/get-my-network-info", internals.GetMyNetworkInfo).Methods("GET")
	r.HandleFunc("/ws/chat", hub.ChatWS)
	// r.HandleFunc("/fs/list", internals.ListFiles).Methods("GET")
	// r.HandleFunc("/peer/fs", internals.GetPeerFiles).Methods("GET")

	//Middlewares
	go utils.BroadcastListener()
	go utils.CleanupPeers()
	go utils.BroadcastSender()
	r.Use(utils.CORSMiddleware)

	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
