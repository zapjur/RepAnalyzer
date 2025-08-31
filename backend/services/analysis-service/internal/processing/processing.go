package processing

import (
	"analysis-service/internal/minio"
	"analysis-service/types"
	"context"
	"log"
	"os"
	"path/filepath"
)

func GenerateAnalysis(ctx context.Context, minioClient *minio.Client, req types.AnalysisRequest) error {
	tmpDir, err := downloadCSVToTmp(ctx, minioClient, req.Bucket, req.ObjectKey, req.Auth0Id)
	if err != nil {
		return err
	}

	posePath := filepath.Join(tmpDir, "pose.csv")
	barPath := filepath.Join(tmpDir, "barpath.csv")

	defer func() {
		_ = os.Remove(posePath)
		_ = os.Remove(barPath)
		_ = os.RemoveAll(tmpDir)
	}()
	log.Println("Downloaded CSVs to", tmpDir)

	return nil
}
