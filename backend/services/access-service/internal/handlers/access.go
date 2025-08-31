package handler

import (
	"access-service/internal/auth"
	"access-service/internal/service"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type AccessHandler struct {
	svc *service.AccessService
}

func NewAccessHandler(svc *service.AccessService) *AccessHandler {
	return &AccessHandler{svc: svc}
}

func (h *AccessHandler) GetPresignedURL(w http.ResponseWriter, r *http.Request) {
	log.Println("GetPresignedURL called")
	user, err := auth.GetUserInfo(r.Context())
	if err != nil {
		log.Println("Unauthorized access attempt:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	videoID := chi.URLParam(r, "videoId")
	videoIDint, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		log.Println("Invalid video ID:", videoID, err)
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}
	url, err := h.svc.GeneratePresignedURL(r.Context(), user.Auth0ID, videoIDint)
	if err != nil {
		log.Println("Access denied for user", user.Auth0ID, "for video ID", videoIDint, ":", err)
		http.Error(w, "Access denied: "+err.Error(), http.StatusForbidden)
		return
	}
	log.Println(url)
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(map[string]string{"url": url})

}

func (h *AccessHandler) GetVideoAnalysis(w http.ResponseWriter, r *http.Request) {
	log.Println("GetVideoAnalysis called")

	user, err := auth.GetUserInfo(r.Context())
	if err != nil {
		log.Println("Unauthorized access attempt:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	videoID := chi.URLParam(r, "videoId")
	videoIDint, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		log.Println("Invalid video ID:", videoID, err)
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	bucket := r.URL.Query().Get("bucket")
	objectKey := r.URL.Query().Get("objectKey")
	videoType := r.URL.Query().Get("type")

	if bucket == "" || objectKey == "" || videoType == "" {
		log.Println("Missing bucket or objectKey or type query parameters")
		http.Error(w, "Missing bucket or objectKey or type query parameters", http.StatusBadRequest)
		return
	}

	switch videoType {
	case "barpath", "pose":
	default:
		http.Error(w, "Invalid type", http.StatusBadRequest)
		return
	}

	url, err := h.svc.GeneratePresignedAnalysisURL(r.Context(), user.Auth0ID, videoIDint, bucket, objectKey)
	if err != nil {
		log.Printf("Access denied for user %s for video ID %d: %v", user.Auth0ID, videoIDint, err)
		http.Error(w, "Access denied: "+err.Error(), http.StatusForbidden)
		return
	}

	type resp struct {
		URL    string  `json:"url"`
		CSVURL *string `json:"csvUrl,omitempty"`
	}

	out := resp{URL: url}

	if videoType == "barpath" {
		csvObjectKey := replaceExt(objectKey, ".csv")
		csvURL, err := h.svc.GeneratePresignedAnalysisURL(r.Context(), user.Auth0ID, videoIDint, bucket, csvObjectKey)
		if err != nil {
			log.Printf("Access denied (CSV) for user %s for video ID %d: %v", user.Auth0ID, videoIDint, err)
			http.Error(w, "Access denied: "+err.Error(), http.StatusForbidden)
			return
		}
		out.CSVURL = &csvURL
	}
	log.Println(videoType, url, out.CSVURL)
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(out); err != nil {
		log.Println("Failed to encode JSON:", err)
	}
}

func replaceExt(key, newExt string) string {
	ext := path.Ext(key)
	base := strings.TrimSuffix(key, ext)
	if !strings.HasPrefix(newExt, ".") {
		newExt = "." + newExt
	}
	return base + newExt
}
