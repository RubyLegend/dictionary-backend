package dictionarytowords

import (
	"log"
	"testing"

	db "github.com/RubyLegend/dictionary-backend/middleware/database"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.OpenConnection()
}

func TestGetWords(t *testing.T) {
	dictionaryID := 1
	page := 1
	limit := 1

	// Assuming GetWords returns an error in this example
	words, total, err := GetWords(dictionaryID, page, limit)

	if err == nil {
		t.Errorf("Expected an error, but received none")
	}

	expectedNumWords := 0
	if len(words) != expectedNumWords {
		t.Errorf("Received incorrect number of words: received %v, expected %v", len(words), expectedNumWords)
	}

	if total == expectedNumWords {
		t.Errorf("Received incorrect total count: received %v, expected %v", total, expectedNumWords)
	}
}
func TestAddConnection(t *testing.T) {

	err := AddConnection(1, 2)

	if err == nil {
		t.Errorf("Received an error: received %v, expected %v", err, nil)
	}

}
