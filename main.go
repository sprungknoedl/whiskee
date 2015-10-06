package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set.")
	}

	http.HandleFunc("/", IndexR)
	http.ListenAndServe(":"+port, nil)
}

func IndexR(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}
