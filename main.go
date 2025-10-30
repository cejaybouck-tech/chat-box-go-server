package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)



func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
		return
	}

	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		ServeWebSockets(hub, w, r)
	});

	// Start the HTTP server
	port := os.Getenv("PORT")
	log.Printf("Starting server at %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}