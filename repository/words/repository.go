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
	WordId    int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

func checkWordExistance(word Word) error {
	// for _, v := range WordAndDictionaryTable {
	// 	for _, w := range Words {
	// 		if w.Name == words.Name {
	// 			err = append(err, errors.New("Word "+words.Name+" already exists"))
	// 		}
	// 	}
	// }

	return nil
}

func postValidation(wordData Word) error {

	if len(wordData.Name) == 0 {
		return errors.New("name is required field")
	}

	return nil
}

func WordIDtoWords(dictToWords []dictionaryToWordsRepo.DictionaryToWords) ([]Word, error) {
	var words []Word
	var wordIds []int

	if len(dictToWords) == 0 {
		return []Word{}, nil
	}

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

func AddWord(wordData Word) (int, error) {

	err := postValidation(wordData)

	if err != nil {
		return -1, err
	}

	err = checkWordExistance(wordData)

	if err != nil {
		return -1, err
	}

	dbCon := db.GetConnection()

	res, err := dbCon.Exec("insert into Words values (default, ?, CURRENT_TIMESTAMP())", wordData.Name)

	if err != nil {
		return -1, err
	}

	lastId, err := res.LastInsertId()

	return int(lastId), err
}
