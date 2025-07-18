//nolint:gofmt
package free_dictionary //nolint:stylecheck

type WordInfoErrorResp struct {
	Title      string `json:"title"`
	Message    string `json:"message"`
	Resolution string `json:"resolution"`
}

type WordInfoResp struct {
	Word      string `json:"word"`
	Phonetic  string `json:"phonetic"`
	Phonetics []struct {
		Text      string `json:"text"`
		Audio     string `json:"audio"`
		SourceURL string `json:"sourceUrl"`
		License   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"license"`
	} `json:"phonetics"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string        `json:"definition"`
			Synonyms   []interface{} `json:"synonyms"`
			Antonyms   []interface{} `json:"antonyms"`
			Example    string        `json:"example,omitempty"`
		} `json:"definitions"`
		Synonyms []string `json:"synonyms"`
		Antonyms []string `json:"antonyms"`
	} `json:"meanings"`
	License struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"license"`
	SourceURLs []string `json:"sourceUrls"`
}
