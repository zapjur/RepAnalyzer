package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func ConvertToMP4(inputPath string) (string, error) {
	if strings.HasSuffix(strings.ToLower(inputPath), ".mp4") {
		return inputPath, nil
	}

	ext := filepath.Ext(inputPath)
	outputPath := strings.TrimSuffix(inputPath, ext) + ".mp4"

	cmd := exec.Command("ffmpeg", "-y", "-i", inputPath, "-vcodec", "libx264", "-acodec", "aac", outputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg conversion failed: %v\n%s", err, string(output))
	}

	return outputPath, nil
}
