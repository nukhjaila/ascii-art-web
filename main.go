package main

import (
	"ascii-art-web/handlers"
	"net/http"
)

func main() {

	http.HandleFunc("/", handlers.InputHandler)

	http.ListenAndServe(":9091", nil)

}
