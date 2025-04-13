package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"

	"api-gateway/internal/grpc"
	dbPb "api-gateway/proto/db"

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

	auth0ID := chi.URLParam(r, "auth0ID")
	if auth0ID == "" {
		http.Error(w, "Missing user ID in path", http.StatusBadRequest)
		return
	}

	resp, err := h.client.DBService.GetUser(r.Context(), &dbPb.GetUserRequest{
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
