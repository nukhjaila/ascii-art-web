package handlers

import (
	"fmt"
	"io"
	"net/http"
)

func inputHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("failed to read http request:", err)
	}

	input := string(httpRequestBody)
}
