package words

import (
	"log"
	"testing"

	db "github.com/RubyLegend/dictionary-backend/middleware/database"
	dictionaryToWordsRepo "github.com/RubyLegend/dictionary-backend/repository/dictionaryToWords"
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

	validWord := Word{
		Name:        "valid",
		Translation: "valid translation",
	}

	err := postValidation(validWord)
	if err != nil {
		t.Errorf("Received an error: %v, expected nil", err)
	}

	invalidWord := Word{
		Name:        "",
		Translation: "invalid translation",
	}

	err = postValidation(invalidWord)
	if err == nil {
		t.Error("Expected an error for empty name, but got nil")
	}
}
func TestWordIDtoWords(t *testing.T) {

	dictToWords := []dictionaryToWordsRepo.DictionaryToWords{}

	words, err := WordIDtoWords(dictToWords)
	if err != nil {
		t.Errorf("Received an error: %v, expected nil", err)
	}

	if len(words) != 0 {
		t.Errorf("Expected 0 words, got: %d", len(words))
	}

	dictToWords = []dictionaryToWordsRepo.DictionaryToWords{
		{
			WordId: 2,
		},
	}

	words, err = WordIDtoWords(dictToWords)
	if err != nil {
		t.Errorf("Received an error: %v, expected nil", err)
	}

	if len(words) != 1 {
		t.Errorf("Expected 1 word, got: %d", len(words))
	}

	if words[0].WordId != 2 {
		t.Errorf("Expected word ID: 1, got: %d", words[0].WordId)
	}
}

func TestAddWord(t *testing.T) {

	wordData := Word{
		Name:        "new word",
		Translation: "new translation",
	}


	_, _, err := AddWord(wordData)

	if err != nil {
		t.Errorf("Received an error: %v, expected nil", err)
	}


}

func TestUpdateWord(t *testing.T) {

	wordData := Word{
		WordId:      1,
		Name:        "updated word",
		Translation: "updated translation",
	}

	err := UpdateWord(wordData)

	if err != nil {
		t.Errorf("Received an error: %v, expected nil", err)
	}

}

func TestDeleteWord(t *testing.T) {

	wordData := Word{
		WordId: 1,
	}

	err := DeleteWord(wordData)

	if err == nil {
		t.Errorf("Received an error: %v, expected nil", err)
	}

}
