package history

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
func TestGetHistory(t *testing.T) {

	histories, err := GetHistory(1)

	if err != nil {
		t.Errorf("Received an error: received %v, expected %v", err, nil)
	}

	expectedCount := 0
	if len(histories) != expectedCount {
		t.Errorf("Received incorrect number of histories: received %v, expected %v", len(histories), expectedCount)
	}
}

func TestDeleteHistory(t *testing.T) {

	err := DeleteHistory(1, 1)

	if err != nil {
		t.Errorf("Received an error: received %v, expected %v", err, nil)
	}

}
