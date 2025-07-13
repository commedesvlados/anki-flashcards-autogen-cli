package core

import "time"

// Flashcard represents a single flashcard with Russian-English word pair
type Flashcard struct {
	ID           int       `json:"id"`
	Russian      string    `json:"russian"`
	English      string    `json:"english"`
	PartOfSpeech string    `json:"part_of_speech"`
	Definition   string    `json:"definition"`
	Example      string    `json:"example"`
	IPAUK        string    `json:"ipa_uk"`
	IPAUS        string    `json:"ipa_us"`
	AudioUK      string    `json:"audio_uk"` // e.g., "uk.mp3"
	AudioUS      string    `json:"audio_us"` // e.g., "us.mp3"
	ImagePath    string    `json:"image"`    // e.g., "apple.jpg"
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ExportFlash is the struct used for genanki export
type ExportFlash struct {
	ID           int    `json:"id"`
	Russian      string `json:"russian"`
	English      string `json:"english"`
	PartOfSpeech string `json:"part_of_speech"`
	Definition   string `json:"definition"`
	Example      string `json:"example"`
	IPAUK        string `json:"ipa_uk"`
	IPAUS        string `json:"ipa_us"`
	AudioUK      string `json:"audio_uk"` // e.g., "uk.mp3"
	AudioUS      string `json:"audio_us"` // e.g., "us.mp3"
	ImagePath    string `json:"image"`    // e.g., "apple.jpg"
}

// ToExportFlash converts Flashcard to ExportFlash
func (f *Flashcard) ToExportFlash() *ExportFlash {
	return &ExportFlash{
		ID:           f.ID,
		Russian:      f.Russian,
		English:      f.English,
		PartOfSpeech: f.PartOfSpeech,
		Definition:   f.Definition,
		Example:      f.Example,
		IPAUK:        f.IPAUK,
		IPAUS:        f.IPAUS,
		AudioUK:      f.AudioUK,
		AudioUS:      f.AudioUS,
		ImagePath:    f.ImagePath,
	}
}

// RawFlashcard represents the basic word pair from Excel
type RawFlashcard struct {
	Russian      string `json:"russian"`
	English      string `json:"english"`
	PartOfSpeech string `json:"part_of_speech"`
}
