//nolint:stylecheck
package free_dictionary

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

// API client for Free Dictionary API
type API struct {
	client  *http.Client
	baseURL string
	logger  *zap.Logger
}

// NewAPI creates a new dictionary API client
func NewAPI(logger *zap.Logger) *API {
	return &API{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.dictionaryapi.dev/api/v2/entries/en",
		logger:  logger,
	}
}

// GetWordInfo retrieves word information from Free Dictionary API
func (api *API) GetWordInfo(ctx context.Context, word string) ([]WordInfoResp, error) {
	// Clean and encode the word
	cleanWord := url.QueryEscape(word)
	apiURL := fmt.Sprintf("%s/%s", api.baseURL, cleanWord)

	api.logger.Debug("Fetching word info", zap.String("word", word), zap.String("url", apiURL))

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

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
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("word '%s' not found", word)
	}

	if resp.StatusCode != http.StatusOK {
		// Try to parse error response
		var errorResp WordInfoErrorResp
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			// If we can parse the error response, use it
			if errorResp.Title != "" {
				return nil, fmt.Errorf("API error: %s - %s", errorResp.Title, errorResp.Message)
			}
		}
		// Fallback to generic error
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Parse response
	var apiResponse []WordInfoResp
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if len(apiResponse) == 0 {
		return nil, fmt.Errorf("no data returned for word '%s'", word)
	}

	api.logger.Debug("Successfully fetched word info", zap.String("word", word), zap.Int("entries", len(apiResponse)))
	return apiResponse, nil
}
