package handlers

import (
	"api-gateway/internal/auth"
	"api-gateway/internal/grpc"
	dbPb "api-gateway/proto/db"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	client *grpc.Client
}

func NewUserHandler(client *grpc.Client) *UserHandler {
	return &UserHandler{client: client}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUserInfo(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	resp, err := h.client.DBService.GetUser(r.Context(), &dbPb.GetUserRequest{
		Auth0Id: user.Auth0ID,
		Email:   user.Email,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
