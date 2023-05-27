package users

import (
	"errors"
	"strconv"
	"time"

	userRepo "github.com/RubyLegend/dictionary-backend/repository/users"
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

	// userRepo.Users = append(userRepo.Users, User{UserId: 1, Email: "ehgfwe", Password: "dssf", CreatedAt: time.Now()})

	for _, v := range Dictionaries {
		if v.Name == dictionaryData.Name && v.UserId == dictionaryData.UserId {
			err = append(err, errors.New("Dictionary "+dictionaryData.Name+" already exists"))
		}
	}

	return err
}

func postValidation(dictionaryData Dictionary) []error {
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

func GetDictionary(UserId int) ([]error, Dictionary) {
	var err []error
	// userRepo.Users = append(userRepo.Users, userRepo.User{UserId: 1, Username: "puk", Email: "ehgfwe", Password: "dssf", CreatedAt: time.Now()})
	// Dictionaries = append(Dictionaries, Dictionary{DictionaryId: 2, UserId: 1, Name: "ehgfwe", CreatedAt: time.Now()})

	var FinedDictionary Dictionary

	for _, v := range Dictionaries {
		if v.UserId == UserId {
			FinedDictionary = v
		}
	}

	if (FinedDictionary == Dictionary{}) {
		err = append(err, errors.New("Dictionary not found"))
	}

	return err, FinedDictionary
}

func AddDictionary(dictionaryData Dictionary) []error {
	var err []error

	// userRepo.Users = append(userRepo.Users, userRepo.User{UserId: 1, Username: "puk", Email: "ehgfwe", Password: "dssf", CreatedAt: time.Now()})

	err = append(err, postValidation(dictionaryData)...)
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

func DeleteDictionary(UserId int, DictionaryId string) []error {
	var err []error
	id, error := strconv.Atoi(DictionaryId)

	if error != nil {
		err = append(err, errors.New("Invalid id params"))
	}

	// userRepo.Users = append(userRepo.Users, userRepo.User{UserId: 2, Username: "puk", Email: "ehgfwe", Password: "dssf", CreatedAt: time.Now()})
	// Dictionaries = append(Dictionaries, Dictionary{DictionaryId: 2, UserId: 2, Name: "ehgfwe", CreatedAt: time.Now()})

	var isDeleted = false
	if err == nil {
		for i, v := range Dictionaries {
			if v.UserId == UserId && v.DictionaryId == id {
				Dictionaries = append(Dictionaries[:i], Dictionaries[i+1:]...)
				isDeleted = true
				break
			}
		}
		if !isDeleted {
			err = append(err, errors.New("Dictionary not found"))
		}
	}

	return err
}

func updateValidation(dictionaryData Dictionary) []error {
	var err []error
	var isFieldInRequest = false
	if len(dictionaryData.Name) != 0 {
		isFieldInRequest = true
	}
	if !dictionaryData.CreatedAt.IsZero() {
		isFieldInRequest = true
	}
	if !isFieldInRequest {
		err = append(err, errors.New("Incorrect body"))
	}
	return err
}

func UpdateDictionary(UserId int, DictionaryId string, dictionaryData Dictionary) ([]error, Dictionary) {
	var err []error
	id, error := strconv.Atoi(DictionaryId)

	if error != nil {
		err = append(err, errors.New("Invalid id params"))
	}

	err = append(err, updateValidation(dictionaryData)...)

	userRepo.Users = append(userRepo.Users, userRepo.User{UserId: 2, Username: "puk", Email: "ehgfwe", Password: "dssf", CreatedAt: time.Now()})
	Dictionaries = append(Dictionaries, Dictionary{DictionaryId: 2, UserId: 2, Name: "dictionary", CreatedAt: time.Now()})

	var UpdatedDictionary Dictionary
	if err == nil {
		for i, v := range Dictionaries {
			if v.UserId == UserId && v.DictionaryId == id {
				Dictionaries[i] = dictionaryData
				UpdatedDictionary = Dictionaries[i]
				break
			}
		}
		if UpdatedDictionary == (Dictionary{}) {
			err = append(err, errors.New("Dictionary not found"))
		}
	}

	return err, UpdatedDictionary
}
