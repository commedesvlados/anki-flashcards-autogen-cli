package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/commedesvlados/anki-flashcards-autogen-cli/internal/util"

	"go.uber.org/zap"
)

// Downloader handles downloading media files
type Downloader struct {
	client   *http.Client
	mediaDir string
	logger   *zap.Logger
}

// NewDownloader creates a new media downloader
func NewDownloader(mediaDir string, logger *zap.Logger) *Downloader {
	return &Downloader{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		mediaDir: mediaDir,
		logger:   logger,
	}
}

// sanitizeFilename removes or replaces special characters to make filename safe
func sanitizeFilename(filename string) string {
	// Replace spaces and special characters with underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	sanitized := re.ReplaceAllString(filename, "_")

	// Remove multiple consecutive underscores
	re = regexp.MustCompile(`_+`)
	sanitized = re.ReplaceAllString(sanitized, "_")

	// Remove leading/trailing underscores
	sanitized = strings.Trim(sanitized, "_")

	return sanitized
}

// downloadFile is a helper for downloading a file from a URL and saving it to the media directory.
func (d *Downloader) downloadFile(ctx context.Context, url, filename, logType string) (string, error) {
	safeFilename := sanitizeFilename(filename)
	d.logger.Debug("Downloading "+logType, zap.String("url", url), zap.String("filename", safeFilename))

	if err := os.MkdirAll(d.mediaDir, 0755); err != nil { //nolint:mnd
		return "", fmt.Errorf("failed to create media directory: %w", err)
	}

	filePath := filepath.Join(d.mediaDir, safeFilename)
	if _, err := os.Stat(filePath); err == nil {
		d.logger.Debug(logType+" already exists", zap.String("path", filePath))
		return safeFilename, nil
	}

	var resp *http.Response
	err := util.RetryWithBackoff(ctx, 3, func() error {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil) //nolint:gocritic
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		resp, err = d.client.Do(req) //nolint:bodyclose
		if err != nil {
			return fmt.Errorf("failed to download "+logType+": %w", err)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download %s: status %d", logType, resp.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	d.logger.Debug("Successfully downloaded "+logType, zap.String("path", filePath))
	return safeFilename, nil
}

func (d *Downloader) DownloadImage(ctx context.Context, imageURL, filename string) (string, error) {
	return d.downloadFile(ctx, imageURL, filename, "image")
}

func (d *Downloader) DownloadAudio(ctx context.Context, audioURL, filename string) (string, error) {
	return d.downloadFile(ctx, audioURL, filename, "audio")
}
