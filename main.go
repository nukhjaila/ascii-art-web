package main

import (
	"ascii-art-web/handlers"
	"net/http"
)

func main() {

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/input", handlers.InputHandler)
	http.HandleFunc("/standard", handlers.StandardHandler)
	http.HandleFunc("/thinkertoy", handlers.ThinkertoyHandler)
	http.HandleFunc("/shadow", handlers.ShadowHandler)

	http.ListenAndServe(":9091", nil)

}
