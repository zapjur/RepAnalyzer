package handler

import (
	"access-service/internal/auth"
	"access-service/internal/service"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type AccessHandler struct {
	svc *service.AccessService
}

func NewAccessHandler(svc *service.AccessService) *AccessHandler {
	return &AccessHandler{svc: svc}
}

func (h *AccessHandler) GetPresignedURL(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUserInfo(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	videoID := chi.URLParam(r, "videoId")
	videoIDint, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}
	url, err := h.svc.GeneratePresignedURL(r.Context(), user.Auth0ID, videoIDint)
	if err != nil {
		http.Error(w, "Access denied: "+err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"url": url})
}
