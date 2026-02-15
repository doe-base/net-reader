package utils

import (
	"encoding/json"
	"fmt"
	"net/http" // Ensure this is imported
	"os"
	"sync" // Required for peerMu
)

var peerMu2 sync.Mutex // Add this to your utils package

func StartFileServer(port string, sharedDir string) {
	mux := http.NewServeMux()

	// Corrected signature: (w http.ResponseWriter, r *http.Request)
	mux.HandleFunc("/api/list", func(w http.ResponseWriter, r *http.Request) {
		files, err := os.ReadDir(sharedDir)
		if err != nil {
			http.Error(w, "Unable to read directory", http.StatusInternalServerError)
			return
		}

		var fileNames []string
		for _, f := range files {
			fileNames = append(fileNames, f.Name())
		}

		w.Header().Set("Content-Type", "应用/json")
		json.NewEncoder(w).Encode(fileNames)
	})

	// Serve actual files
	fs := http.FileServer(http.Dir(sharedDir))
	mux.Handle("/files/", http.StripPrefix("/files/", fs))

	fmt.Printf("📂 File server starting on :%s\n", port)

	// Using a local server instance is cleaner for Go apps
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
