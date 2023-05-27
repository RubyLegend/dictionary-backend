package translations

import (
	"errors"
	"strconv"

	wordsRepo "github.com/RubyLegend/dictionary-backend/repository/words"
)

type Translation struct {
	TranslationId int    `json:"translaonId"`
	WordId        int    `json:"wordId"`
	Name          string `json:"name"`
	Language      string `json:"language"`
}

var Translations []Translation

func checkTranslationExistance(translationData Translation) []error {
	var err []error

	for _, v := range Translations {
		if v.Name == translationData.Name && v.WordId == translationData.WordId {
			err = append(err, errors.New("translation "+translationData.Name+" already exists"))
		}
	}

	return err
}

func validation(translationData Translation) []error {
	var err []error

	if len(translationData.Name) == 0 {
		err = append(err, errors.New("name is required field"))
	}
	found := false
	for _, translation := range wordsRepo.Words {
		if translation.WordId == translationData.WordId {
			found = true
			break
		}
	}
	if !found {
		err = append(err, errors.New("occured some error"))
	}
	return err
}

// func findTranslation(params ...interface{}) (int, error) {
// 	translationData, ok := params[0].(Translation)
// 	if ok {
// 		for i, v := range Translations {
// 			if v.Name == translationData.Name {
// 				return i, nil
// 			}
// 		}
// 	} else {
// 		name, ok := params[0].(string)
// 		if ok {
// 			for i, v := range Translations {
// 				if v.Name == name {
// 					return i, nil
// 				}
// 			}
// 		} else {
// 			return -1, errors.New("Unknown parameter passed")
// 		}
// 	}

//		return -1, errors.New("Translation not found")
//	}
func GetTranslation(TranslationId int) ([]error, Translation) {
	var err []error
	Translations = append(Translations, Translation{TranslationId: 1, WordId: 1, Name: "ehgfwe", Language: "ehgfwe"})
	var FinedTranslation Translation

	for _, v := range Translations {
		if v.TranslationId == TranslationId {
			FinedTranslation = v
		}
	}

	if (FinedTranslation == Translation{}) {
		err = append(err, errors.New("translation not found"))
	}

	return err, FinedTranslation
}

func AddTranslation(translationData Translation) []error {
	var err []error

	err = append(err, validation(translationData)...)
	err = append(err, checkTranslationExistance(translationData)...)

	if err == nil {
		lastElementIndex := len(Translations) - 1
		if lastElementIndex < 0 {
			translationData.WordId = 0
		} else {
			translationData.WordId = Translations[lastElementIndex].WordId + 1
		}

		Translations = append(Translations, translationData)

		return nil
	} else {
		return nil
	}
}

func DeleteTranslation(WordId int, TranslationId string) []error {
	var err []error
	id, error := strconv.Atoi(TranslationId)

	if error != nil {
		err = append(err, errors.New("invalid id params"))
		return err
	}

	var isDeleted = false
	for i, v := range Translations {
		if v.WordId == WordId && v.TranslationId == id {
			Translations = append(Translations[:i], Translations[i+1:]...)
			isDeleted = true
			break
		}
	}

	if !isDeleted {
		err = append(err, errors.New("translation not found"))
	}

	return err
}
func updateValidation(translationData Translation) []error {
	var err []error
	var isFieldInRequest = false
	if len(translationData.Name) != 0 {
		isFieldInRequest = true
	}
	if len(translationData.Language) != 0 {
		isFieldInRequest = true
	}
	if !isFieldInRequest {
		err = append(err, errors.New("incorrect body"))
	}
	return err
}

func UpdateTranslation(WordId int, TranslationId string, translationData Translation) ([]error, Translation) {
	var err []error
	id, error := strconv.Atoi(TranslationId)

	if error != nil {
		err = append(err, errors.New("invalid id params"))
	}

	err = append(err, updateValidation(translationData)...)

	var UpdatedTranslation Translation
	if err == nil {
		for i, v := range Translations {
			if v.WordId == WordId && v.TranslationId == id {
				Translations[i] = translationData
				UpdatedTranslation = Translations[i]
				break
			}
		}
		if UpdatedTranslation == (Translation{}) {
			err = append(err, errors.New("translation not found"))
		}
	}

	return err, UpdatedTranslation

}
