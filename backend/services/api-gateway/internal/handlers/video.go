package handlers

import (
	"api-gateway/internal/auth"
	"api-gateway/internal/grpc"
	miniohelpers "api-gateway/internal/minio"
	"api-gateway/internal/utils"
	dbPb "api-gateway/proto/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/minio/minio-go/v7"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type VideoHandler struct {
	minio                  *minio.Client
	grpcDBClient           *grpc.Client
	grpcOrchestratorClient *grpc.Client
}

type MinioObjectRef struct {
	Bucket    string
	ObjectKey string
}

func NewVideoHandler(minioClient *minio.Client, grpcDBClient, grpcOrchestratorClient *grpc.Client) *VideoHandler {
	return &VideoHandler{minio: minioClient, grpcDBClient: grpcDBClient, grpcOrchestratorClient: grpcOrchestratorClient}
}

func (h *VideoHandler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	exercise, cleanFilename, tmpInput, err := h.parseAndPrepareFormData(w, r)
	if err != nil {
		return
	}
	defer os.Remove(tmpInput.Name())
	defer tmpInput.Close()

	user, err := auth.GetUserInfo(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	auth0ID := user.Auth0ID
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

	minioObjectRef, err := h.uploadToMinIO(auth0IDEdited, exercise, cleanFilename, convertedFile, fileInfo)
	if err != nil {
		http.Error(w, "Failed to upload to MinIO: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err = miniohelpers.WaitForObject(h.minio, minioObjectRef.Bucket, minioObjectRef.ObjectKey, 5, 1*time.Second); err != nil {
		http.Error(w, "MinIO object not available after upload: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Uploaded to MinIO and verified presence: %s/%s", minioObjectRef.Bucket, minioObjectRef.ObjectKey)
		return
	}

	dbResp, err := h.saveVideoToDB(auth0ID, minioObjectRef, exercise)
	if err != nil {
		http.Error(w, "Failed to save video info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !dbResp.Success {
		http.Error(w, "Failed to save video info: "+dbResp.Message, http.StatusInternalServerError)
		return
	}

	videoID := dbResp.VideoId

	resp, err := h.sendVideoToAnalyze(minioObjectRef, exercise, auth0ID, videoID)
	if err != nil {
		http.Error(w, "Failed to send video for analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !resp.Success {
		http.Error(w, "Failed to send video for analysis: "+resp.Message, http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "File uploaded, video id: %d", videoID)

}

func (h *VideoHandler) GetVideosByExercise(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUserInfo(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	exercise := chi.URLParam(r, "exercise")
	if exercise == "" {
		http.Error(w, "Missing exercise path parameter", http.StatusBadRequest)
		return
	}

	response, err := h.grpcDBClient.DBService.GetUserVideosByExercise(context.Background(), &dbPb.GetUserVideosByExerciseRequest{
		Auth0Id:      user.Auth0ID,
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
