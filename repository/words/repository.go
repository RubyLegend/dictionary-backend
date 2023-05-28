package words

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	db "github.com/RubyLegend/dictionary-backend/middleware/database"
	dictionaryToWordsRepo "github.com/RubyLegend/dictionary-backend/repository/dictionaryToWords"
	translationRepo "github.com/RubyLegend/dictionary-backend/repository/translations"
)

type Word struct {
	WordId      int       `json:"id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
	IsLearned   bool      `json:"isLearned"`
	Translation string    `json:"translation"`
}

type WordWithDictId struct {
	Name         string `json:"name"`
	DictionaryId int    `json:"dictionaryId"`
	Translation  string `json:"translation"`
}

func (wordData WordWithDictId) ConvertToWord() Word {
	var word Word
	word.Name = wordData.Name
	word.Translation = wordData.Translation

	return word
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

	query, args, err := sqlx.In("select w.*, t.name from Words w left join Translation t on t.wordID = w.WordID where w.wordID in (?) order by createdAt desc", wordIds)

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
		err = rows.Scan(&word.WordId, &word.Name, &word.CreatedAt, &word.IsLearned, &word.Translation)
		if err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	return words, nil
}

func AddWord(wordData Word) (int, Word, error) {

	err := postValidation(wordData)

	if err != nil {
		return -1, Word{}, err
	}

	err = checkWordExistance(wordData)

	if err != nil {
		return -1, Word{}, err
	}

	dbCon := db.GetConnection()

	res, err := dbCon.Exec("insert into Words values (default, ?, CURRENT_TIMESTAMP(), default)", wordData.Name)

	if err != nil {
		return -1, Word{}, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		return -1, Word{}, err
	}

	err = translationRepo.AddTranslation(int(lastId), wordData.Translation)
	if err != nil {
		return -1, Word{}, err
	}

	wordData.WordId = int(lastId)
	wordData.CreatedAt = time.Now()

	return int(lastId), wordData, err
}

func UpdateWord(wordData Word) error {
	dbCon := db.GetConnection()

	_, err := dbCon.Exec("update Words set name = ? where wordId = ?", wordData.Name, wordData.WordId)

	if err != nil {
		return err
	}

	var translationData translationRepo.Translation
	translationData.WordId = wordData.WordId
	translationData.Name = wordData.Translation
	err = translationRepo.UpdateTranslation(translationData)

	if err != nil {
		err = translationRepo.AddTranslation(wordData.WordId, wordData.Translation)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteWord(wordData Word) error {
	dbCon := db.GetConnection()

	res, err := dbCon.Exec("delete from Words where wordId = ?", wordData.WordId)

	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.New("this word doesn't exist")
	}

	return nil
}
