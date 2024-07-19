package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var LOGGER = log.New(os.Stdout, "INFO:", log.Ldate)

func ping(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Pong\n"))
}

func call(w http.ResponseWriter, req *http.Request) {
	port := os.Getenv("APP_PORT")
	LOGGER.Println("got a call request in port: ", port)
	w.WriteHeader(200)
	io.Copy(w, req.Body)
}

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "9000"
	}
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/call", call)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
