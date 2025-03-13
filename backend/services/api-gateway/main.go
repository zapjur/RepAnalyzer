package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	pb "api-gateway/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("user-service:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Nie można połączyć się z user-service: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		auth0ID := strings.TrimPrefix(r.URL.Path, "/users/")

		resp, err := client.GetUser(context.Background(), &pb.GetUserRequest{Auth0Id: auth0ID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("API Gateway działa na porcie :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Błąd serwera HTTP: %v", err)
	}
}
