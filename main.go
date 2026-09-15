package main

import (
	"ascii-art-web/handlers"
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/input", handlers.InputHandler)
	http.HandleFunc("/standard", handlers.StandardHandler)
	http.HandleFunc("/thinkertoy", handlers.ThinkertoyHandler)
	http.HandleFunc("/shadow", handlers.ShadowHandler)

	fmt.Println("The server started at: localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to run http server", err)
	}

}
