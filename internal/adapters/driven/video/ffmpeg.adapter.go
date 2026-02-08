package video

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"video_processor_service/internal/core/domain/ports"
)

type FFmpegVideoProcessor struct {
	storageService ports.IStorageService
}

func NewFFmpegVideoProcessor(storageService ports.IStorageService) *FFmpegVideoProcessor {
	return &FFmpegVideoProcessor{
		storageService: storageService,
	}
}

func (f *FFmpegVideoProcessor) GetVideoDuration(ctx context.Context, bucket, key string) (float64, error) {
	tmpFile, err := f.downloadToTemp(ctx, bucket, key)
	if err != nil {
		return 0, fmt.Errorf("failed to download video: %w", err)
	}
	defer os.Remove(tmpFile)

	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		tmpFile,
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return 0, fmt.Errorf("failed to get video duration: %w, stderr: %s", err, stderr.String())
	}

	durationStr := strings.TrimSpace(out.String())
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	return duration, nil
}

func (f *FFmpegVideoProcessor) ExtractFrame(ctx context.Context, bucket, key string, timestamp float64) ([]byte, error) {
	tmpFile, err := f.downloadToTemp(ctx, bucket, key)
	if err != nil {
		return nil, fmt.Errorf("failed to download video: %w", err)
	}
	defer os.Remove(tmpFile)

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-ss", fmt.Sprintf("%.2f", timestamp),
		"-i", tmpFile,
		"-vframes", "1",
		"-f", "image2pipe",
		"-vcodec", "mjpeg",
		"-",
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to extract frame: %w, stderr: %s", err, stderr.String())
	}

	return out.Bytes(), nil
}

func (f *FFmpegVideoProcessor) downloadToTemp(ctx context.Context, bucket, key string) (string, error) {
	size, err := f.storageService.GetObjectMetadata(ctx, bucket, key)
	if err != nil {
		return "", err
	}

	reader, err := f.storageService.GetObjectRange(ctx, bucket, key, 0, size-1)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("video_%s", filepath.Base(key)))

	file, err := os.Create(tmpFile)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		os.Remove(tmpFile)
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}

	return tmpFile, nil
}
