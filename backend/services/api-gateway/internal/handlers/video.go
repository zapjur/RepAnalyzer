package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type VideoHandler struct{}

func NewVideoHandler() *VideoHandler {
	return &VideoHandler{}
}

func (h *VideoHandler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100 << 20) // 100 MB
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Could not get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	exercise := r.FormValue("exercise")
	if exercise == "" {
		http.Error(w, "Missing exercise field", http.StatusBadRequest)
		return
	}

	uploadDir := "./uploads/" + exercise
	os.MkdirAll(uploadDir, os.ModePerm)
	filePath := filepath.Join(uploadDir, handler.Filename)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Uploaded successfully to %s", filePath)
}
