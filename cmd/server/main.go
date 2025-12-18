package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./web/static/"))

	mux.Handle("/", fs)

	log.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
