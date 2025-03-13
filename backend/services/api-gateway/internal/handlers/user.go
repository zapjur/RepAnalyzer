package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"api-gateway/internal/grpc"
	pb "api-gateway/proto"
)

type UserHandler struct {
	client *grpc.Client
}

func NewUserHandler(client *grpc.Client) *UserHandler {
	return &UserHandler{client: client}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	auth0ID := strings.TrimPrefix(r.URL.Path, "/users/")

	resp, err := h.client.UserService.GetUser(r.Context(), &pb.GetUserRequest{Auth0Id: auth0ID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
