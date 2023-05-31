package translations

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
func TestPostValidation(t *testing.T) {
	invalidData := Translation{
		WordId:   -1,
		Name:     "",
		Language: "",
	}

	err := postValidation(invalidData)

	if err == nil {
		t.Error("Expected an error, but received nil")
	} else {
		expectedError := "name is required field"
		if err.Error() != expectedError {
			t.Errorf("Received incorrect error message: received %v, expected %v", err.Error(), expectedError)
		}
	}

	validData := Translation{
		WordId:   1,
		Name:     "translationName",
		Language: "English",
	}

	err = postValidation(validData)

	if err != nil {
		t.Errorf("Received an error: received %v, expected %v", err, nil)
	}
}
func TestCheckTranslationExistance(t *testing.T) {
	existingTranslation := Translation{
		TranslationId: 1,
		WordId:        1,
		Name:          "",
		Language:      "English",
	}

	err := checkTranslationExistance(existingTranslation)

	if err != nil {
		t.Errorf("Received an error: %v, expected nil", err)
	}
}

func TestGetTranslation(t *testing.T) {

	translations, err := GetTranslation(1)

	if err != nil {
		t.Errorf("Received an error: received %v, expected %v", err, nil)
	}

	expectedCount := 0
	if len(translations) != expectedCount {
		t.Errorf("Received incorrect number of translations: received %v, expected %v", len(translations), expectedCount)
	}
}
func TestAddTranslation(t *testing.T) {
	wordID := 1
	translation := "new translation"

	err := AddTranslation(wordID, translation)

	if err == nil {
		t.Errorf("Received an error: %v, expected %v", err, nil)
	}

}
func TestUpdateTranslation(t *testing.T) {

	translation := Translation{
		TranslationId: 1,
		WordId:        1,
		Name:          "updated translation",
		Language:      "English",
	}

	err := UpdateTranslation(translation)

	if err != nil {
		t.Errorf("Received an error: received %v, expected %v", err, nil)
	}

}

func TestDeleteTranslation(t *testing.T) {
	translation := Translation{
		TranslationId: 1,
		WordId:        1,
		Name:          "updated translation",
		Language:      "English",
	}
	err := DeleteTranslation(translation)

	if err == nil {
		t.Errorf("Received an error: %v, expected nil", err)
	}
}
