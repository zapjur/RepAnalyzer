package main

import (
	"access-service/internal/router"
	"log"
	"net/http"
)

func main() {
	r := router.Setup()

	log.Println("Access Service running on :8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}
