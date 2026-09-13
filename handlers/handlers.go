package handlers

import (
	"ascii-art-web/pkg/banner"
	"ascii-art-web/pkg/parser"
	"ascii-art-web/pkg/renderer"
	"fmt"
	"io"
	"net/http"
	"os"
)

func InputHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("failed to read http request:", err)
	}

	input := string(httpRequestBody)

	bannerMap, err := banner.Load("banners/standard.txt")
	if err != nil {
		fmt.Printf("Error loading template: %v\n", err)
	}
	if input == "\\n" {
		fmt.Println()
		return
	}
	matrix, err := parser.Parse(input)
	if err != nil {
		fmt.Printf("Error parsing text line: %v\n", err)
		os.Exit(1)
	}

	renderer.Render(os.Stdout, matrix, bannerMap)
}
