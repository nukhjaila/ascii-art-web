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

<<<<<<< HEAD
	fmt.Println("Server started at localhost:9091")
=======
	fmt.Println("The server started at: localhost:8080")
>>>>>>> dc5fd21946c56b466a2f03fda0d1065f55bf512f
	if err := http.ListenAndServe(":9091", nil); err != nil {
		fmt.Println("Failed to run http server", err)
	}

}
