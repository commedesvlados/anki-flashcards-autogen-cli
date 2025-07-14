package unsplash

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/commedesvlados/anki-flashcards-autogen-cli/internal/util"

	"go.uber.org/zap"
)

// API client for Unsplash API
type API struct {
	client    *http.Client
	baseURL   string
	accessKey string
	logger    *zap.Logger
}

// NewAPI creates a new image API client
func NewAPI(accessKey string, logger *zap.Logger) *API {
	return &API{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:   "https://api.unsplash.com/search/photos",
		accessKey: accessKey,
		logger:    logger,
	}
}

// GetImageURL retrieves an image URL from Unsplash API
func (api *API) GetImageURL(ctx context.Context, query string) (string, error) {
	// Build URL with query parameters
	params := url.Values{}
	params.Add("query", query)
	params.Add("per_page", "1")
	params.Add("orientation", "landscape")

	apiURL := fmt.Sprintf("%s?%s", api.baseURL, params.Encode())

	api.logger.Debug("Fetching image URL", zap.String("query", query), zap.String("url", apiURL))

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", api.accessKey))

	// Add retry logic
	var resp *http.Response
	err = util.RetryWithBackoff(ctx, 3, func() error {
		resp, err = api.client.Do(req) //nolint:bodyclose
		if err != nil {
			return fmt.Errorf("failed to make request: %w", err)
		}
		return nil
	})

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Try to parse error response
		var errorResp ImageErrorResp
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			// If we can parse the error response, use it
			if len(errorResp.Errors) > 0 {
				return "", fmt.Errorf("API error: %s", errorResp.Errors[0])
			}
		}

		// Handle specific status codes
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return "", fmt.Errorf("unauthorized: invalid access token")
		case http.StatusForbidden:
			return "", fmt.Errorf("forbidden: missing permissions")
		case http.StatusTooManyRequests:
			return "", fmt.Errorf("rate limit exceeded: too many requests")
		default:
			return "", fmt.Errorf("API returned status %d", resp.StatusCode)
		}
	}

	// Parse response
	var apiResponse ImageResp
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResponse.Results) == 0 {
		return "", fmt.Errorf("no images found for query '%s'", query)
	}

	imageURL := apiResponse.Results[0].URLs.Regular
	api.logger.Debug("Successfully fetched image URL", zap.String("query", query), zap.String("url", imageURL))
	return imageURL, nil
}
