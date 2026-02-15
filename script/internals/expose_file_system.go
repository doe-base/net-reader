package internals

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func ListFiles(w http.ResponseWriter, r *http.Request) {
	baseDir := "/home" // 👈 LIMIT ACCESS (important)

	path := r.URL.Query().Get("path")
	if path == "" {
		path = baseDir
	}

	// Security: prevent directory escape
	if !strings.HasPrefix(path, baseDir) {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type FileEntry struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	var files []FileEntry
	for _, e := range entries {
		t := "file"
		if e.IsDir() {
			t = "dir"
		}
		files = append(files, FileEntry{
			Name: e.Name(),
			Type: t,
		})
	}

	json.NewEncoder(w).Encode(map[string]any{
		"path":    path,
		"entries": files,
	})
}

func GetPeerFiles(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	path := r.URL.Query().Get("path")

	if ip == "" {
		http.Error(w, "missing peer ip", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("http://%s:8080/fs/list?path=%s", ip, path)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "peer unreachable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}
