package handlers

import (
	orPb "api-gateway/proto/analysis"
	dbPb "api-gateway/proto/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func (h *VideoHandler) parseAndPrepareFormData(w http.ResponseWriter, r *http.Request) (exercise, cleanFilename string, tmpFile *os.File, err error) {
	err = r.ParseMultipartForm(100 << 20) // 100 MB
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

	cleanFilename = strings.NewReplacer("(", "_", ")", "_", " ", "_").Replace(handler.Filename)

	exercise = r.FormValue("exercise")
	if exercise == "" {
		http.Error(w, "Missing exercise field", http.StatusBadRequest)
		return
	}
	exercise = strings.ReplaceAll(exercise, " ", "_")

	tmpFile, err = os.CreateTemp("", "upload-*"+filepath.Ext(cleanFilename))
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}

	file.Seek(0, io.SeekStart)
	_, err = io.Copy(tmpFile, file)
	if err != nil {
		http.Error(w, "Failed to write uploaded file", http.StatusInternalServerError)
		return
	}

	return exercise, cleanFilename, tmpFile, nil
}

func openConvertedFile(path string) (*os.File, os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, nil, err
	}
	return f, info, nil
}

func (h *VideoHandler) uploadToMinIO(auth0ID, exercise, baseFilename string, file *os.File, info os.FileInfo) (*MinioObjectRef, error) {
	baseFilename = strings.TrimSuffix(baseFilename, filepath.Ext(baseFilename))
	objectKey := fmt.Sprintf("%s/%s/original/%s.mp4", auth0ID, exercise, baseFilename)
	bucket := "videos"

	_, err := h.minio.PutObject(context.Background(), bucket, objectKey, file, info.Size(), minio.PutObjectOptions{
		ContentType: "video/mp4",
	})
	if err != nil {
		return nil, err
	}

	return &MinioObjectRef{
		Bucket:    bucket,
		ObjectKey: objectKey,
	}, nil
}

func (h *VideoHandler) saveVideoToDB(auth0ID string, ref *MinioObjectRef, exercise string) (*dbPb.UploadVideoResponse, error) {
	resp, err := h.grpcDBClient.DBService.SaveUploadedVideo(context.Background(), &dbPb.UploadVideoRequest{
		Auth0Id:      auth0ID,
		Bucket:       ref.Bucket,
		ObjectKey:    ref.ObjectKey,
		ExerciseName: exercise,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *VideoHandler) sendVideoToAnalyze(ref *MinioObjectRef, exercise, auth0ID string, videoID int64) (*orPb.VideoToAnalyzeResponse, error) {
	resp, err := h.grpcOrchestratorClient.OrchestratorService.AnalyzeVideo(context.Background(), &orPb.VideoToAnalyzeRequest{
		Bucket:       ref.Bucket,
		ObjectKey:    ref.ObjectKey,
		ExerciseName: exercise,
		Auth0Id:      auth0ID,
		VideoId:      videoID,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func fetchPresignedURL(ctx context.Context, authorization string, videoID int64) string {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		fmt.Sprintf("%s/access/video/%d", "http://access-service:8082", videoID),
		nil,
	)
	if err != nil {
		return ""
	}
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var payload struct {
		URL string `json:"url"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return ""
	}
	return payload.URL
}

func fetchPresignedAnalysisURL(ctx context.Context, authorization string, videoID int64, bucket, objectKey string) string {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		fmt.Sprintf("http://access-service:8082/access/video-analysis/%d?bucket=%s&objectKey=%s",
			videoID,
			url.QueryEscape(bucket),
			url.QueryEscape(objectKey),
		),
		nil,
	)
	if err != nil {
		return ""
	}
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var payload struct {
		URL string `json:"url"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return ""
	}
	return payload.URL
}
