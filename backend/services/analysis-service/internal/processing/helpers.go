package processing

import (
	"analysis-service/internal/minio"
	"context"
	"fmt"
	miniosdk "github.com/minio/minio-go/v7"
	"os"
	"path/filepath"
	"strings"
)

func downloadCSVToTmp(ctx context.Context, minioClient *minio.Client, bucket, objectKey, auth0ID string) (string, error) {
	base := filepath.Join(os.TempDir(), auth0ID)
	if err := os.MkdirAll(base, 0o700); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", base, err)
	}

	rootDir := filepath.Dir(filepath.Dir(objectKey))

	baseName := strings.TrimSuffix(filepath.Base(objectKey), filepath.Ext(objectKey)) + ".csv"
	poseKey := filepath.ToSlash(filepath.Join(rootDir, "pose", baseName))
	barKey := filepath.ToSlash(filepath.Join(rootDir, "barpath", baseName))

	tmpDir, err := os.MkdirTemp(base, "analysis_csv_")
	if err != nil {
		return "", fmt.Errorf("mkdtemp: %w", err)
	}

	poseDst := filepath.Join(tmpDir, "pose.csv")
	barDst := filepath.Join(tmpDir, "barpath.csv")

	if err := minioClient.Minio.FGetObject(ctx, bucket, poseKey, poseDst, miniosdk.GetObjectOptions{}); err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", fmt.Errorf("download pose CSV (%s/%s): %w", bucket, poseKey, err)
	}
	if err := minioClient.Minio.FGetObject(ctx, bucket, barKey, barDst, miniosdk.GetObjectOptions{}); err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", fmt.Errorf("download barpath CSV (%s/%s): %w", bucket, barKey, err)
	}

	return tmpDir, nil
}
