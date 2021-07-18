package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	client = http.Client{
		Timeout: 15 * time.Second,
	}
)

func main() {
	port := *flag.Int("port", 8080, "the port to use for the server")
	flag.Parse()

	go func() {
		r := mux.NewRouter()

		r.HandleFunc("/{id}", FileGetHandler).Methods("GET")
		http.Handle("/", r)

		log.Printf("Started listening on :%d", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}()

	setupUI()
}
