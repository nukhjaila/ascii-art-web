package handlers

import (
	"ascii-art-web/pkg/banner"
	"ascii-art-web/pkg/parser"
	"ascii-art-web/pkg/renderer"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

var (
	filePath = "banners/standard.txt"
	mu       sync.Mutex
)

func InputHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("failed to read http request:", err)
	}

	input := string(httpRequestBody)

	filePath := getFilePath()

	bannerMap, err := banner.Load(filePath)
	if err != nil {
		fmt.Printf("Error loading template: %v\n", err)
	}

	matrix, err := parser.Parse(input)
	if err != nil {
		fmt.Printf("Error parsing text line: %v\n", err)
		os.Exit(1)
	}

	res := renderer.Render(matrix, bannerMap)

	_, err = w.Write(res)
	if err != nil {
		fmt.Println("failed to write response:", err)
	}

}

func getFilePath() string {
	mu.Lock()
	defer mu.Unlock()

	return filePath
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/index.html")
}

func StandardHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	filePath = "banners/standard.txt"
	mu.Unlock()

	fmt.Fprintf(w, "%s", filePath)
}

func ShadowHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	filePath = "banners/shadow.txt"
	mu.Unlock()

	fmt.Fprintf(w, "%s", filePath)
}
func ThinkertoyHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	filePath = "banners/thinkertoy.txt"
	mu.Unlock()

	fmt.Fprintf(w, "%s", filePath)
}
