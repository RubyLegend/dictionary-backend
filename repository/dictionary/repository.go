package users

import (
	"time"
	//   "fmt"
	userRepo "github.com/RubyLegend/dictionary-backend/repository/users"
	//   "log"
	"errors"
)

type Dictionary struct {
	DictionaryId int       `json:"dictionaryId"`
	UserId       int       `json:"userId"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"createdAt"`
}

var Dictionaries []Dictionary

func checkDictionaryExistance(dictionaryData Dictionary) []error {
	var err []error

	for _, v := range Dictionaries {
		if v.Name == dictionaryData.Name {
			err = append(err, errors.New("Dictionary "+dictionaryData.Name+" already exists"))
		}
	}

	return err
}

func validation(dictionaryData Dictionary) []error {
	var err []error

	if len(dictionaryData.Name) == 0 {
		err = append(err, errors.New("Name is required field"))
	}


	found := false
	for _, user := range userRepo.Users {
		if user.UserId == dictionaryData.UserId {
			found = true
			break
		}
	}

	if !found {
		err = append(err, errors.New("Occurred some error"))
	}

	return err
}

func AddDictionary(dictionaryData Dictionary) []error {
	var err []error

	err = append(err, validation(dictionaryData)...)
	err = append(err, checkDictionaryExistance(dictionaryData)...)

	if err == nil {
		lastElementIndex := len(Dictionaries) - 1
		if lastElementIndex < 0 {
			dictionaryData.DictionaryId = 0
		} else {
			dictionaryData.DictionaryId = Dictionaries[lastElementIndex].DictionaryId + 1
		}

		dictionaryData.CreatedAt = time.Now()
		Dictionaries = append(Dictionaries, dictionaryData)

		return nil
	} else {
		return err
	}

}
