package handlers

import (
	"api-gateway/internal/grpc"
	pb "api-gateway/proto"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/minio/minio-go/v7"
	"net/http"
	"strings"
)

type VideoHandler struct {
	minio      *minio.Client
	grpcClient *grpc.Client
}

func NewVideoHandler(minioClient *minio.Client, grpcClient *grpc.Client) *VideoHandler {
	return &VideoHandler{minio: minioClient, grpcClient: grpcClient}
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

	cleanFilename := strings.NewReplacer("(", "_", ")", "_", " ", "_").Replace(handler.Filename)

	exercise := r.FormValue("exercise")
	if exercise == "" {
		http.Error(w, "Missing exercise field", http.StatusBadRequest)
		return
	}
	exercise = strings.ReplaceAll(exercise, " ", "_")

	auth0ID := chi.URLParam(r, "auth0ID")
	if auth0ID == "" {
		http.Error(w, "Missing user ID in path", http.StatusBadRequest)
		return
	}

	auth0IDEdited := strings.ReplaceAll(auth0ID, "|", "_")

	objectName := fmt.Sprintf("%s/%s/%s", auth0IDEdited, exercise, cleanFilename)
	contentType := handler.Header.Get("Content-Type")

	uploadInfo, err := h.minio.PutObject(context.Background(), "videos", objectName, file, handler.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		http.Error(w, "Failed to upload to MinIO: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded to MinIO: %s (size: %d bytes)", uploadInfo.Key, uploadInfo.Size)

	baseURL := "http://localhost:9000"
	bucketName := "videos"
	objectURL := fmt.Sprintf("%s/%s/%s", baseURL, bucketName, uploadInfo.Key)

	h.grpcClient.UserService.SaveUploadedVideo(context.Background(), &pb.UploadVideoRequest{
		Auth0Id:      auth0ID,
		Url:          objectURL,
		ExerciseName: exercise,
	})
}

func (h *VideoHandler) GetVideosByExercise(w http.ResponseWriter, r *http.Request) {
	auth0ID := chi.URLParam(r, "auth0ID")
	if auth0ID == "" {
		http.Error(w, "Missing user ID in path", http.StatusBadRequest)
		return
	}

	exercise := chi.URLParam(r, "exercise")
	if exercise == "" {
		http.Error(w, "Missing exercise path parameter", http.StatusBadRequest)
		return
	}

	response, err := h.grpcClient.UserService.GetUserVideosByExercise(context.Background(), &pb.GetUserVideosByExerciseRequest{
		Auth0Id:      auth0ID,
		ExerciseName: exercise,
	})
	if err != nil {
		http.Error(w, "Failed to get videos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !response.Success {
		http.Error(w, "Failed to get videos: "+response.Message, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.Videos)
}
