package dictionary

import (
	"database/sql"
	"log"
	"testing"
	"time"

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
func TestCheckDictionaryExistance(t *testing.T) {
	userID := 1
	name := "existing dictionary"

	err := CheckDictionaryExistance(userID, name)

	if err != nil {
		t.Errorf("didn't expect an error but dictionary already exists")
	}

}
func TestPostValidation(t *testing.T) {

	dictionaryData := DictionaryPost{
		Name: "newDictionary",
	}
	userId := 1

	_ = postValidation(userId, dictionaryData)

	dictionaryData = DictionaryPost{
		Name: "",
	}

	errs := postValidation(userId, dictionaryData)

	{
		expectedError := "name is required field"
		if errs[0].Error() != expectedError {
			t.Errorf("Received incorrect error message: received %v, expected %v", errs[0].Error(), expectedError)
		}
	}

	dictionaryData = DictionaryPost{
		Name: "newDictionary",
	}
	userId = 2

	_ = postValidation(userId, dictionaryData)
}
func TestGetDictionary(t *testing.T) {

	dictionaries, err := GetDictionary(1)

	if err != nil {
		t.Errorf("Received an error: received %v, expected %v", err, nil)
	}

	expectedCount := 2
	if len(dictionaries) != expectedCount {
		t.Errorf("Received incorrect number of dictionaries: received %v, expected %v", len(dictionaries), expectedCount)
	}
}

func TestAddDictionary(t *testing.T) {
	dictionaryData := DictionaryPost{
		Name: "newDictionary",
	}

	err := AddDictionary(3, dictionaryData)

	if err == nil {
		t.Errorf("Received an error: %v, expected %v", err, nil)
	}

	err = AddDictionary(3, dictionaryData)

	if err == nil {
		t.Errorf("Expected an error, but received nil")
	} else {
		expectedError := "dictionary already exists"
		if err[0].Error() == expectedError {
			t.Errorf("Received incorrect error message: received %v, expected %v", err[0].Error(), expectedError)
		}
	}

	emptyNameData := DictionaryPost{
		Name: "",
	}

	err = AddDictionary(1, emptyNameData)

	if err == nil {
		t.Errorf("Expected an error, but received nil")
	} else {
		expectedError := "name is required field"
		if err[0].Error() != expectedError {
			t.Errorf("Received incorrect error message: received %v, expected %v", err[0].Error(), expectedError)
		}
	}

	nonExistentUser := DictionaryPost{
		Name: "newDictionary",
	}

	err = AddDictionary(2, nonExistentUser)

	if err == nil {
		t.Errorf("Expected an error, but received nil")
	} else {
		expectedError := "dictionary owner not found"
		if err[0].Error() == expectedError {
			t.Errorf("Received incorrect error message: received %v, expected %v", err[0].Error(), expectedError)
		}
	}
}
func TestDeleteDictionary(t *testing.T) {

	err := DeleteDictionary(1, 1)

	if err != nil {
		t.Errorf("Received an error: received %v, expected %v", err, nil)
	}

}

func TestUpdateValidation(t *testing.T) {
	emptyNameData := Dictionary{
		Name: "",
	}

	errs := updateValidation(emptyNameData)

	if errs == nil {
		t.Errorf("Expected 0 errors, but received error : %v", errs)
	}

	validData := Dictionary{
		Name: "updatedDictionary",
	}

	errs = updateValidation(validData)

	if errs != nil {
		t.Errorf("Expected no errors, but received error : %v", errs)
	}
}
func TestUpdateDictionary(t *testing.T) {
	dictionary, err := UpdateDictionary(1, 5, Dictionary{
		DictionaryId: 5,
		UserId:       1,
		Name:         "testtest",
		CreatedAt:    time.Now(),
	})

	if err == nil {
		if err != sql.ErrNoRows {
			t.Errorf("Dictionary or user not found: %v", err)
		} else {
			t.Errorf("Received an unexpected error: %v", err)
		}
		return
	}

	expectedName := "updatedDictionary"
	if dictionary.Name == expectedName {
		t.Errorf("Failed to update dictionary name: received %v, expected %v", dictionary.Name, expectedName)
	}
}
