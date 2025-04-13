package handlers

import (
	"api-gateway/internal/grpc"
	"api-gateway/internal/utils"
	dbPb "api-gateway/proto/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/minio/minio-go/v7"
	"net/http"
	"os"
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
	auth0ID, exercise, cleanFilename, tmpInput, err := h.parseAndPrepareFormData(w, r)
	if err != nil {
		return
	}
	defer os.Remove(tmpInput.Name())
	defer tmpInput.Close()

	auth0IDEdited := strings.ReplaceAll(auth0ID, "|", "_")

	convertedPath, err := utils.ConvertToMP4(tmpInput.Name())
	if err != nil {
		http.Error(w, "Conversion failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(convertedPath)

	convertedFile, fileInfo, err := openConvertedFile(convertedPath)
	if err != nil {
		http.Error(w, "Failed to open converted file", http.StatusInternalServerError)
		return
	}
	defer convertedFile.Close()

	objectURL, err := h.uploadToMinIO(auth0IDEdited, exercise, cleanFilename, convertedFile, fileInfo)
	if err != nil {
		http.Error(w, "Failed to upload to MinIO: "+err.Error(), http.StatusInternalServerError)
		return
	}

	videoID, err := h.saveVideoToDB(auth0ID, objectURL, exercise)
	if err != nil {
		http.Error(w, "Failed to save video info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded to minIO, video id: %s", videoID)
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

	response, err := h.grpcClient.DBService.GetUserVideosByExercise(context.Background(), &dbPb.GetUserVideosByExerciseRequest{
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
