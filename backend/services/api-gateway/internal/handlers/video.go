package handlers

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"net/http"
)

type VideoHandler struct {
	minio *minio.Client
}

func NewVideoHandler(minioClient *minio.Client) *VideoHandler {
	return &VideoHandler{minio: minioClient}
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

	objectName := fmt.Sprintf("%s/%s", exercise, handler.Filename)
	contentType := handler.Header.Get("Content-Type")

	uploadInfo, err := h.minio.PutObject(context.Background(), "videos", objectName, file, handler.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		http.Error(w, "Failed to upload to MinIO: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded to MinIO: %s (size: %d bytes)", uploadInfo.Key, uploadInfo.Size)
}
