package unsplash

// ImageErrorResp represents the Unsplash API error response.
type ImageErrorResp struct {
	Errors []string `json:"errors"`
}

// Response represents the Unsplash API response.
type ImageResp struct {
	Results []struct {
		URLs struct {
			Regular string `json:"regular"`
		} `json:"urls"`
	} `json:"results"`
}
