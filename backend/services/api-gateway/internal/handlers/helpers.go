package handlers

import (
	dbPb "api-gateway/proto/db"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/minio/minio-go/v7"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (h *VideoHandler) parseAndPrepareFormData(w http.ResponseWriter, r *http.Request) (auth0ID, exercise, cleanFilename string, tmpFile *os.File, err error) {
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

	auth0ID = chi.URLParam(r, "auth0ID")
	if auth0ID == "" {
		http.Error(w, "Missing user ID in path", http.StatusBadRequest)
		return
	}

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

	return auth0ID, exercise, cleanFilename, tmpFile, nil
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

func (h *VideoHandler) uploadToMinIO(auth0ID, exercise, baseFilename string, file *os.File, info os.FileInfo) (string, error) {
	baseFilename = strings.TrimSuffix(baseFilename, filepath.Ext(baseFilename))
	objectName := fmt.Sprintf("%s/%s/%s.mp4", auth0ID, exercise, baseFilename)

	uploadInfo, err := h.minio.PutObject(context.Background(), "videos", objectName, file, info.Size(), minio.PutObjectOptions{
		ContentType: "video/mp4",
	})
	if err != nil {
		return "", err
	}

	baseURL := "http://localhost:9000"
	bucketName := "videos"
	objectURL := fmt.Sprintf("%s/%s/%s", baseURL, bucketName, uploadInfo.Key)
	return objectURL, nil
}

func (h *VideoHandler) saveVideoToDB(auth0ID, url, exercise string) (int64, error) {
	resp, err := h.grpcClient.DBService.SaveUploadedVideo(context.Background(), &dbPb.UploadVideoRequest{
		Auth0Id:      auth0ID,
		Url:          url,
		ExerciseName: exercise,
	})
	if err != nil {
		return 0, err
	}
	return resp.GetVideoId(), nil
}
