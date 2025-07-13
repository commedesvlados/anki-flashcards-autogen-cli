package core

import (
	"testing"
	"time"
)

func TestFlashcard_ToExportFlash(t *testing.T) {
	now := time.Now()
	flashcard := &Flashcard{
		ID:           1,
		Russian:      "яблоко",
		English:      "apple",
		PartOfSpeech: "noun",
		Definition:   "A round fruit with red or green skin",
		Example:      "I eat an apple every day",
		IPAUK:        "/ˈæp.əl/",
		IPAUS:        "/ˈæp.əl/",
		AudioUK:      "1_uk.mp3",
		AudioUS:      "1_us.mp3",
		ImagePath:    "1.jpg",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	export := flashcard.ToExportFlash()

	// Verify all fields are correctly copied
	if export.ID != flashcard.ID {
		t.Errorf("Expected ID %d, got %d", flashcard.ID, export.ID)
	}
	if export.Russian != flashcard.Russian {
		t.Errorf("Expected Russian %s, got %s", flashcard.Russian, export.Russian)
	}
	if export.English != flashcard.English {
		t.Errorf("Expected English %s, got %s", flashcard.English, export.English)
	}
	if export.PartOfSpeech != flashcard.PartOfSpeech {
		t.Errorf("Expected PartOfSpeech %s, got %s", flashcard.PartOfSpeech, export.PartOfSpeech)
	}
	if export.Definition != flashcard.Definition {
		t.Errorf("Expected Definition %s, got %s", flashcard.Definition, export.Definition)
	}
	if export.Example != flashcard.Example {
		t.Errorf("Expected Example %s, got %s", flashcard.Example, export.Example)
	}
	if export.IPAUK != flashcard.IPAUK {
		t.Errorf("Expected IPAUK %s, got %s", flashcard.IPAUK, export.IPAUK)
	}
	if export.IPAUS != flashcard.IPAUS {
		t.Errorf("Expected IPAUS %s, got %s", flashcard.IPAUS, export.IPAUS)
	}
	if export.AudioUK != flashcard.AudioUK {
		t.Errorf("Expected AudioUK %s, got %s", flashcard.AudioUK, export.AudioUK)
	}
	if export.AudioUS != flashcard.AudioUS {
		t.Errorf("Expected AudioUS %s, got %s", flashcard.AudioUS, export.AudioUS)
	}
	if export.ImagePath != flashcard.ImagePath {
		t.Errorf("Expected ImagePath %s, got %s", flashcard.ImagePath, export.ImagePath)
	}
}

func TestRawFlashcard(t *testing.T) {
	raw := &RawFlashcard{
		Russian: "дом",
		English: "house",
	}

	if raw.Russian != "дом" {
		t.Errorf("Expected Russian 'дом', got %s", raw.Russian)
	}
	if raw.English != "house" {
		t.Errorf("Expected English 'house', got %s", raw.English)
	}
}
