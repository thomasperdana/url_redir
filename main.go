package main

import (
	"fmt"
	"log"
	"net/http"

	"api/api"
)

func main() {
	http.HandleFunc("/", api.UrlHandler)
	
	fmt.Println("Starting URL redirector server on :8080")
	fmt.Println("Visit http://localhost:8080/example to test redirects")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
