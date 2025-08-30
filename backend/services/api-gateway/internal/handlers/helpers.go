package handlers

import (
	orPb "api-gateway/proto/analysis"
	dbPb "api-gateway/proto/db"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/oklog/ulid/v2"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var entropy = ulid.Monotonic(rand.Reader, 0)

func GenerateULID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}

func (h *VideoHandler) parseAndPrepareFormData(w http.ResponseWriter, r *http.Request) (string, string, *os.File, error) {
	err := r.ParseMultipartForm(100 << 20) // 100 MB
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return "", "", nil, err
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Could not get file from request", http.StatusBadRequest)
		return "", "", nil, err
	}
	defer file.Close()

	filename := GenerateULID()

	exercise := r.FormValue("exercise")
	if exercise == "" {
		http.Error(w, "Missing exercise field", http.StatusBadRequest)
		return "", "", nil, fmt.Errorf("missing exercise field")
	}
	exercise = strings.ReplaceAll(exercise, " ", "_")

	tmpFile, err := os.CreateTemp("", "upload-*"+filepath.Ext(filename))
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return "", "", nil, err
	}

	file.Seek(0, io.SeekStart)
	_, err = io.Copy(tmpFile, file)
	if err != nil {
		http.Error(w, "Failed to write uploaded file", http.StatusInternalServerError)
		return "", "", nil, err
	}

	return exercise, filename, tmpFile, nil
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

func (h *VideoHandler) uploadToMinIO(ctx context.Context, auth0ID, exercise, baseFilename string, file *os.File, info os.FileInfo) (*MinioObjectRef, error) {
	objectKey := fmt.Sprintf("%s/%s/original/%s.mp4", auth0ID, exercise, baseFilename)
	bucket := "videos"

	_, err := h.minio.PutObject(ctx, bucket, objectKey, file, info.Size(), minio.PutObjectOptions{
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

func fetchPresignedAnalysisURL(ctx context.Context, authorization string, videoID int64, bucket, objectKey, videoType string) (string, *string) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		fmt.Sprintf("http://access-service:8082/access/video-analysis/%d?bucket=%s&objectKey=%s&type=%s",
			videoID,
			url.QueryEscape(bucket),
			url.QueryEscape(objectKey),
			url.QueryEscape(videoType),
		),
		nil,
	)
	if err != nil {
		return "", nil
	}
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", nil
	}

	var payload struct {
		URL    string  `json:"url"`
		CSVURL *string `json:"csvUrl,omitempty"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", nil
	}
	return payload.URL, payload.CSVURL
}
