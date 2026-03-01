package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"conceptual-lan/internals/communication"
	"conceptual-lan/internals/discovery"
	"conceptual-lan/utils"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()

	// Chat endpoints
	r.HandleFunc("/api/message", communication.ReceiveMessage).Methods("POST")
	r.HandleFunc("/api/send", communication.SendMessageFromBrowser).Methods("POST")
	r.HandleFunc("/api/messages", communication.GetMessages).Methods("GET")

	// Discovery endpoint
	r.HandleFunc("/get-my-network-info", discovery.GetMyNetworkInfo).Methods("GET")

	// File server (must come LAST)
	fs := http.FileServer(http.Dir("/"))
	r.PathPrefix("/").Handler(fs)

	handler := utils.CORSMiddleware(r)

	// Start discovery routines
	go discovery.BroadcastListener()
	go discovery.BroadcastSender()
	go discovery.CleanupPeers()

	fmt.Printf("🚀 Server running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
