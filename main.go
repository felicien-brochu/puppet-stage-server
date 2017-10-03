package main

import (
	"log"
	"net/http"
)

func main() {
	router := getRouter()
	log.Println("Server listens on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
