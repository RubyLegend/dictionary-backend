package dictionarytowords

import (
	db "github.com/RubyLegend/dictionary-backend/middleware/database"
)

type DictionaryToWords struct {
	DictionaryId int `json:"dictionaryId"`
	WordId       int `json:"wordId"`
}

func GetWords(DictionaryId int) ([]DictionaryToWords, error) {
	dbCon := db.GetConnection()

	rows, err := dbCon.Query("select * from DictionariesWords where dictionaryID = ?", DictionaryId)

	if err != nil {
		return nil, err
	}

	var res []DictionaryToWords

	for rows.Next() {
		var word DictionaryToWords
		rows.Scan(&word.DictionaryId, &word.WordId)
		res = append(res, word)
	}

	if len(res) == 0 {
		return []DictionaryToWords{}, nil
	}

	return res, nil
}
