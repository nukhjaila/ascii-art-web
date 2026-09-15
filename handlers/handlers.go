package handlers

import (
	"ascii-art-web/pkg/banner"
	"ascii-art-web/pkg/parser"
	"ascii-art-web/pkg/renderer"
	"fmt"
	"io"
	"net/http"
	"sync"
)

var (
	filePath = "banners/standard.txt"
	mu       sync.Mutex
)

func InputHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not allowed", http.StatusMethodNotAllowed)
		return
	}

	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	input := string(httpRequestBody)

	filePath := getFilePath()

	bannerMap, err := banner.Load(filePath)
	if err != nil {
		fmt.Printf("Error loading template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	matrix, err := parser.Parse(input)
	if err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	res := renderer.Render(matrix, bannerMap)

	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(res); err != nil {
		fmt.Println("failed to write response:", err)
	}

}

func getFilePath() string {
	mu.Lock()
	defer mu.Unlock()

	return filePath
}

func setFilePath(p string) {
	mu.Lock()
	filePath = p
	mu.Unlock()
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "frontend/index.html")
}

func StandardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	setFilePath("banners/standard.txt")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", getFilePath())

}

func ShadowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	setFilePath("banners/shadow.txt")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", getFilePath())
}
func ThinkertoyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	setFilePath("banners/thinkertoy.txt")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", getFilePath())
}
