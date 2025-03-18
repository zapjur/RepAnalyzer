package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"api-gateway/internal/grpc"
	pb "api-gateway/proto"

	"github.com/golang-jwt/jwt/v4"
)

type UserHandler struct {
	client *grpc.Client
}

func NewUserHandler(client *grpc.Client) *UserHandler {
	return &UserHandler{client: client}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	email, err := extractEmailFromJWT(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	auth0ID := strings.TrimPrefix(r.URL.Path, "/users/")

	resp, err := h.client.UserService.GetUser(r.Context(), &pb.GetUserRequest{
		Auth0Id: auth0ID,
		Email:   email,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func extractEmailFromJWT(tokenString string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", err
	}

	return email, nil
}
