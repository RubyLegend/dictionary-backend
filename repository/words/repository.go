package words

import (
	"errors"

	// "log"
	"time"
	//   "fmt"
	//   "bytes"
	"github.com/jmoiron/sqlx"

	db "github.com/RubyLegend/dictionary-backend/middleware/database"
	dictionaryToWordsRepo "github.com/RubyLegend/dictionary-backend/repository/dictionaryToWords"
)

type Word struct {
	WordId    int       `json:"wordId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type RequestType struct {
	DictionaryId int    `json:"dictionaryId"`
	Name         string `json:"name"`
}

type WordAndDictionary struct {
	WordId       int `json:"wordId"`
	DictionaryId int `json:"dictionaryId"`
}

var Words []Word
var WordAndDictionaryTable []WordAndDictionary

func checkWordExistance(words RequestType) []error {
	var err []error
	// for _, v := range WordAndDictionaryTable {
	// 	for _, w := range Words {
	// 		if w.Name == words.Name {
	// 			err = append(err, errors.New("Word "+words.Name+" already exists"))
	// 		}
	// 	}
	// }

	return err
}
func postValidation(wordData RequestType) []error {
	var err []error

	if len(wordData.Name) == 0 {
		err = append(err, errors.New("name is required field"))
	}
	if wordData.DictionaryId == 0 {
		err = append(err, errors.New("dictionaryId is required field"))
	}

	return err
}

func WordIDtoWords(dictToWords []dictionaryToWordsRepo.DictionaryToWords) ([]Word, error) {
	var words []Word
	var wordIds []int

	for _, v := range dictToWords {
		wordIds = append(wordIds, v.WordId)
	}

	dbCon := db.GetConnection()

	query, args, err := sqlx.In("select * from Words where wordID in (?)", wordIds)

	if err != nil {
		return nil, err
	}

	query = dbCon.Rebind(query)
	rows, err := dbCon.Query(query, args...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var word Word
		err = rows.Scan(&word.WordId, &word.Name, &word.CreatedAt)
		if err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	return words, nil
}

func AddWord(wordData RequestType) []error {

	var err []error

	err = append(err, checkWordExistance(wordData)...)
	err = append(err, postValidation(wordData)...)

	if err == nil {

		var newWord Word
		newWord.Name = wordData.Name
		var newConnection WordAndDictionary

		lastElementIndex := len(Words) - 1
		if lastElementIndex < 0 {
			newConnection.DictionaryId = wordData.DictionaryId
			newConnection.WordId = 0
			newWord.WordId = newConnection.WordId
		} else {
			newConnection.WordId = Words[lastElementIndex].WordId + 1
			newWord.WordId = newConnection.WordId
		}

		newWord.CreatedAt = time.Now()

		Words = append(Words, newWord)
		WordAndDictionaryTable = append(WordAndDictionaryTable, newConnection)

		return nil
	} else {
		return err
	}
}
