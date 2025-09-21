package processing

import (
	"analysis-service/internal/minio"
	"context"
	"fmt"
	miniosdk "github.com/minio/minio-go/v7"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func downloadCSVToTmp(ctx context.Context, minioClient *minio.Client, bucket, objectKey, auth0ID string) (string, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
	}

	base := filepath.Join(os.TempDir(), auth0ID)
	if err := os.MkdirAll(base, 0o700); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", base, err)
	}

	tmp, err := os.MkdirTemp(base, "analysis_csv_")
	if err != nil {
		return "", fmt.Errorf("mkdtemp: %w", err)
	}

	rootDir := filepath.Dir(filepath.Dir(objectKey))
	rootDir = strings.ReplaceAll(rootDir, "\\", "/")

	posePrefix := path.Clean(path.Join(rootDir, "pose")) + "/"
	barPrefix := path.Clean(path.Join(rootDir, "barpath")) + "/"

	poseCSVKey, poseMetaKey, err := findOneCSVAndMeta(ctx, minioClient, bucket, posePrefix)
	if err != nil {
		return "", fmt.Errorf("list pose under %q: %w", posePrefix, err)
	}
	barCSVKey, barMetaKey, err := findOneCSVAndMeta(ctx, minioClient, bucket, barPrefix)
	if err != nil {
		return "", fmt.Errorf("list barpath under %q: %w", barPrefix, err)
	}

	poseDst := filepath.Join(tmp, "pose.csv")
	barDst := filepath.Join(tmp, "barpath.csv")
	poseMetaDst := filepath.Join(tmp, "pose_meta.json")
	barMetaDst := filepath.Join(tmp, "barpath_meta.json")

	if err := minioClient.Minio.FGetObject(ctx, bucket, poseCSVKey, poseDst, miniosdk.GetObjectOptions{}); err != nil {
		return "", fmt.Errorf("download pose csv (%s/%s): %w", bucket, poseCSVKey, err)
	}
	if err := minioClient.Minio.FGetObject(ctx, bucket, barCSVKey, barDst, miniosdk.GetObjectOptions{}); err != nil {
		return "", fmt.Errorf("download barpath csv (%s/%s): %w", bucket, barCSVKey, err)
	}

	if poseMetaKey != "" {
		_ = minioClient.Minio.FGetObject(ctx, bucket, poseMetaKey, poseMetaDst, miniosdk.GetObjectOptions{})
	}
	if barMetaKey != "" {
		_ = minioClient.Minio.FGetObject(ctx, bucket, barMetaKey, barMetaDst, miniosdk.GetObjectOptions{})
	}

	return tmp, nil
}

func findOneCSVAndMeta(ctx context.Context, minioClient *minio.Client, bucket, prefix string) (csvKey, metaKey string, err error) {
	var firstCSV, firstMeta string

	for obj := range minioClient.Minio.ListObjects(ctx, bucket, miniosdk.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}) {
		if obj.Err != nil {
			return "", "", obj.Err
		}
		key := obj.Key
		l := strings.ToLower(key)

		if strings.HasSuffix(l, ".csv") && firstCSV == "" {
			firstCSV = key
		}

		if (strings.HasSuffix(l, "_meta.json") ||
			strings.HasSuffix(l, "-meta.json") ||
			strings.HasSuffix(l, "meta.json")) && firstMeta == "" {
			firstMeta = key
		}

		if firstCSV != "" && firstMeta != "" {
			break
		}
	}

	if firstCSV == "" {
		return "", "", fmt.Errorf("no csv found under %s", prefix)
	}
	return firstCSV, firstMeta, nil
}
