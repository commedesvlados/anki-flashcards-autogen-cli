package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/commedesvlados/anki-flashcards-autogen-cli/internal/media"
	free_dictionary "github.com/commedesvlados/anki-flashcards-autogen-cli/pkg/clients/free-dictionary"
	"github.com/commedesvlados/anki-flashcards-autogen-cli/pkg/clients/unsplash"

	"go.uber.org/zap"
)

// EnrichmentService orchestrates the enrichment process
type EnrichmentService struct {
	dictionaryAPI *free_dictionary.API
	imageAPI      *unsplash.API
	downloader    *media.Downloader
	logger        *zap.Logger
}

// NewEnrichmentService creates a new enrichment service
func NewEnrichmentService(dictionaryAPI *free_dictionary.API, imageAPI *unsplash.API, downloader *media.Downloader, logger *zap.Logger) *EnrichmentService {
	return &EnrichmentService{
		dictionaryAPI: dictionaryAPI,
		imageAPI:      imageAPI,
		downloader:    downloader,
		logger:        logger,
	}
}

// EnrichFlashcard enriches a single flashcard with dictionary data and media
func (e *EnrichmentService) EnrichFlashcard(ctx context.Context, raw *RawFlashcard, id int) (*Flashcard, error) {
	e.logger.Info("Enriching flashcard", zap.String("english", raw.English), zap.String("russian", raw.Russian))

	// Check if phrase or not found in dictionary
	isPhrase := isMultiWordPhrase(raw.English)
	dictionaryData, err := e.dictionaryAPI.GetWordInfo(ctx, raw.English)
	if err != nil {
		e.logger.Warn("Failed to get dictionary data", zap.String("word", raw.English), zap.Error(err))
		if isPhrase {
			dictionaryData = nil
		} else {
			mainWord := extractMainWord(raw.English)
			if mainWord != raw.English {
				e.logger.Info("Trying to get dictionary data for main word", zap.String("main_word", mainWord))
				dictionaryData, err = e.dictionaryAPI.GetWordInfo(ctx, mainWord)
				if err != nil {
					e.logger.Warn("Failed to get dictionary data for main word", zap.String("main_word", mainWord), zap.Error(err))
					dictionaryData = nil
				} else {
					e.logger.Info("Successfully got dictionary data for main word", zap.String("main_word", mainWord))
				}
			}
		}
	} else {
		e.logger.Debug("Successfully got dictionary data", zap.String("word", raw.English))
	}

	flashcard := &Flashcard{
		ID:        id,
		Russian:   raw.Russian,
		English:   raw.English,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if dictionaryData != nil && !isPhrase {
		// Use the first entry for meaning/definition (most common usage)
		if len(dictionaryData) > 0 {
			firstEntry := dictionaryData[0]
			if len(firstEntry.Meanings) > 0 {
				meaning := firstEntry.Meanings[0]
				flashcard.PartOfSpeech = meaning.PartOfSpeech
				if len(meaning.Definitions) > 0 {
					flashcard.Definition = meaning.Definitions[0].Definition
				}
			}
		}

		// Find UK and US audio from all entries
		var ukAudio, usAudio, ukIPA, usIPA string

		for _, entry := range dictionaryData {
			for _, phonetic := range entry.Phonetics {
				if phonetic.Audio != "" {
					// Parse audio URL to determine region
					if strings.Contains(phonetic.Audio, "-uk.mp3") {
						if ukAudio == "" {
							ukAudio = phonetic.Audio
							ukIPA = phonetic.Text
						}
					} else if strings.Contains(phonetic.Audio, "-us.mp3") {
						if usAudio == "" {
							usAudio = phonetic.Audio
							usIPA = phonetic.Text
						}
					}
				}
			}
			// If we found both UK and US audio, we can stop searching
			if ukAudio != "" && usAudio != "" {
				break
			}
		}

		// Set IPA (use first available)
		if ukIPA != "" {
			flashcard.IPAUK = ukIPA
		} else if usIPA != "" {
			flashcard.IPAUK = usIPA
		}
		if usIPA != "" {
			flashcard.IPAUS = usIPA
		} else if ukIPA != "" {
			flashcard.IPAUS = ukIPA
		}

		// Download UK audio
		if ukAudio != "" {
			audioUKPath, err := e.downloader.DownloadAudio(ctx, ukAudio, fmt.Sprintf("%d_%s_uk.mp3", id, raw.English))
			if err == nil {
				flashcard.AudioUK = audioUKPath
			}
		}

		// Download US audio
		if usAudio != "" {
			audioUSPath, err := e.downloader.DownloadAudio(ctx, usAudio, fmt.Sprintf("%d_%s_us.mp3", id, raw.English))
			if err == nil {
				flashcard.AudioUS = audioUSPath
			}
		}
		// Image
		imageWord := raw.English
		if isMultiWordPhrase(raw.English) {
			mainWord := extractMainWord(raw.English)
			if mainWord != raw.English {
				imageWord = mainWord
			}
		}
		imageURL, err := e.imageAPI.GetImageURL(ctx, imageWord)
		if err == nil {
			imagePath, err := e.downloader.DownloadImage(ctx, imageURL, fmt.Sprintf("%d_%s.jpg", id, raw.English))
			if err == nil {
				flashcard.ImagePath = imagePath
			}
		}
	} else {
		// Not found or phrase: use only sheet data
		flashcard.PartOfSpeech = raw.PartOfSpeech
		flashcard.Definition = ""
		flashcard.Example = ""
		flashcard.IPAUK = ""
		flashcard.IPAUS = ""
		flashcard.AudioUK = ""
		flashcard.AudioUS = ""
		flashcard.ImagePath = ""
	}

	flashcard.UpdatedAt = time.Now()
	return flashcard, nil
}

// isMultiWordPhrase checks if the given string contains multiple words
func isMultiWordPhrase(phrase string) bool {
	words := strings.Fields(phrase)
	return len(words) > 1
}

// extractMainWord extracts the most significant word from a phrase
// For phrases like "keep something in mind", it returns "keep"
func extractMainWord(phrase string) string {
	words := strings.Fields(phrase)
	if len(words) == 0 {
		return phrase
	}

	// Common words to skip when looking for the main word
	skipWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "from": true, "up": true, "down": true,
		"something": true, "anything": true, "everything": true, "nothing": true,
		"someone": true, "anyone": true, "everyone": true, "no one": true,
	}

	// Find the first significant word
	for _, word := range words {
		cleanWord := strings.ToLower(strings.Trim(word, ".,!?"))
		if !skipWords[cleanWord] {
			return cleanWord
		}
	}

	// If all words are skip words, return the first word
	return strings.ToLower(words[0])
}

// EnrichFlashcards enriches multiple flashcards with progress tracking
func (e *EnrichmentService) EnrichFlashcards(ctx context.Context, rawFlashcards []*RawFlashcard) ([]*Flashcard, error) {
	enriched := make([]*Flashcard, 0, len(rawFlashcards))

	for i, raw := range rawFlashcards {
		select {
		case <-ctx.Done():
			return enriched, ctx.Err()
		default:
		}

		flashcard, err := e.EnrichFlashcard(ctx, raw, i+1)
		if err != nil {
			e.logger.Error("Failed to enrich flashcard", zap.String("english", raw.English), zap.Error(err))
			// Create basic flashcard without enrichment
			flashcard = &Flashcard{
				ID:        i + 1,
				Russian:   raw.Russian,
				English:   raw.English,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}

		enriched = append(enriched, flashcard)
	}

	return enriched, nil
}
