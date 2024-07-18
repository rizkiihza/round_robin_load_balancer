package main

import (
	"fmt"
	"net/http"
	"os"
)

func ping(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Pong\n"))
	req.Body.Close()
}

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "9000"
	}
	http.HandleFunc("/ping", ping)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
