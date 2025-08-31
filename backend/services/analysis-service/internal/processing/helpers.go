package processing

import (
	"analysis-service/internal/minio"
	"context"
	"fmt"
	miniosdk "github.com/minio/minio-go/v7"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func downloadCSVToTmp(ctx context.Context, minioClient *minio.Client, bucket, objectKey, auth0ID string) (string, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	base := filepath.Join(os.TempDir(), auth0ID)
	if err := os.MkdirAll(base, 0o700); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", base, err)
	}

	rootDir := filepath.Dir(filepath.Dir(objectKey))

	stem := strings.TrimSuffix(filepath.Base(objectKey), filepath.Ext(objectKey))

	var poseCSV, barCSV string
	switch {
	case strings.Contains(objectKey, "/pose/"):
		poseCSV = stem + ".csv"
		barCSV = strings.TrimSuffix(stem, "-pose") + "-barpath.csv"
	case strings.Contains(objectKey, "/barpath/"):
		barCSV = stem + ".csv"
		poseCSV = strings.TrimSuffix(stem, "-barpath") + "-pose.csv"
	default:
		poseCSV = stem + "-pose.csv"
		barCSV = stem + "-barpath.csv"
	}

	poseKey := filepath.ToSlash(filepath.Join(rootDir, "pose", poseCSV))
	barKey := filepath.ToSlash(filepath.Join(rootDir, "barpath", barCSV))

	tmpDir, err := os.MkdirTemp(base, "analysis_csv_")
	if err != nil {
		return "", fmt.Errorf("mkdtemp: %w", err)
	}

	poseDst := filepath.Join(tmpDir, "pose.csv")
	barDst := filepath.Join(tmpDir, "barpath.csv")

	if err = minioClient.Minio.FGetObject(ctx, bucket, poseKey, poseDst, miniosdk.GetObjectOptions{}); err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", fmt.Errorf("download pose CSV (%s/%s): %w", bucket, poseKey, err)
	}
	if err = minioClient.Minio.FGetObject(ctx, bucket, barKey, barDst, miniosdk.GetObjectOptions{}); err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", fmt.Errorf("download barpath CSV (%s/%s): %w", bucket, barKey, err)
	}

	return tmpDir, nil
}
