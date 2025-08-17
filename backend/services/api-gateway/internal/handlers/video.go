package handlers

import (
	"api-gateway/internal/auth"
	"api-gateway/internal/grpc"
	miniohelpers "api-gateway/internal/minio"
	"api-gateway/internal/types"
	"api-gateway/internal/utils"
	dbPb "api-gateway/proto/db"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/minio/minio-go/v7"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
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

var httpClient = &http.Client{Timeout: 3 * time.Second}

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

	resp, err := h.grpcDBClient.DBService.GetUserVideosByExercise(
		r.Context(),
		&dbPb.GetUserVideosByExerciseRequest{
			Auth0Id:      user.Auth0ID,
			ExerciseName: exercise,
		})
	if err != nil {
		http.Error(w, "Failed to get videos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !resp.Success {
		http.Error(w, "Failed to get videos: "+resp.Message, http.StatusInternalServerError)
		return
	}

	vids := resp.Videos
	out := make([]types.VideoWithURL, len(vids))

	sem := make(chan struct{}, 8)
	var wg sync.WaitGroup
	authHeader := r.Header.Get("Authorization")

	for i := range vids {
		i := i
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			urlStr := fetchPresignedURL(r.Context(), authHeader, vids[i].Id)
			out[i] = types.VideoWithURL{
				Id:           vids[i].Id,
				Bucket:       vids[i].Bucket,
				ObjectKey:    vids[i].ObjectKey,
				ExerciseName: vids[i].ExerciseName,
				CreatedAt:    vids[i].CreatedAt,
				Url:          urlStr,
			}
			log.Println(out)
		}()
	}
	wg.Wait()

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(out); err != nil {
		http.Error(w, "encode error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *VideoHandler) GetVideoAnalysis(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "videoId")
	if videoID == "" {
		http.Error(w, "Missing video ID path parameter", http.StatusBadRequest)
		return
	}
	videoIDint, err := strconv.ParseInt(videoID, 10, 64)
	if err != nil {
		log.Println("Invalid video ID:", videoID, err)
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}
	log.Println("GetVideoAnalysis called for video ID:", videoIDint)

	resp, err := h.grpcDBClient.DBService.GetVideoAnalysis(
		r.Context(),
		&dbPb.GetVideoAnalysisRequest{
			VideoId: videoIDint,
		},
	)
	if err != nil {
		http.Error(w, "Failed to get videos analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !resp.Success {
		http.Error(w, "Failed to get videos analysis: "+resp.Message, http.StatusInternalServerError)
		return
	}

	vids := resp.Analyses
	out := make([]types.VideoAnalysisWithURL, len(vids))

	sem := make(chan struct{}, 8)
	var wg sync.WaitGroup
	authHeader := r.Header.Get("Authorization")

	for i := range vids {
		i := i
		wg.Add(1)
		sem <- struct{}{}
		log.Println("Fetching presigned URL for video analysis:", vids[i].VideoId, vids[i].Bucket, vids[i].ObjectKey)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			urlStr := fetchPresignedAnalysisURL(r.Context(), authHeader, vids[i].VideoId, vids[i].Bucket, vids[i].ObjectKey)
			out[i] = types.VideoAnalysisWithURL{
				Id:        vids[i].Id,
				Bucket:    vids[i].Bucket,
				ObjectKey: vids[i].ObjectKey,
				Type:      vids[i].Type,
				VideoId:   vids[i].VideoId,
				Url:       urlStr,
			}
			log.Println(out)
		}()
	}
	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	log.Println("Returning video analysis response with url", out[0].Url)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(out); err != nil {
		http.Error(w, "encode error: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
