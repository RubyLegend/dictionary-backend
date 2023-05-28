package dictionarytowords

import (
	db "github.com/RubyLegend/dictionary-backend/middleware/database"
)

type DictionaryToWords struct {
	DictionaryId int `json:"dictionaryId"`
	WordId       int `json:"wordId"`
}

func GetWords(DictionaryId int, page int, limit int) ([]DictionaryToWords, int, error) {
	var count int

	dbCon := db.GetConnection()

	rows, err := dbCon.Query("select dw.* from DictionariesWords dw join Words w on w.wordID = dw.wordID where dictionaryID = ? order by w.createdAt desc limit ?,?",
		DictionaryId, page*limit, limit)

	if err != nil {
		return nil, 0, err
	}

	var res []DictionaryToWords

	for rows.Next() {
		var word DictionaryToWords
		err = rows.Scan(&word.DictionaryId, &word.WordId)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, word)
	}

	if len(res) == 0 {
		return []DictionaryToWords{}, 0, nil
	}

	err = dbCon.QueryRow("select count(*) from DictionariesWords where dictionaryID = ?", DictionaryId).Scan(&count)

	if err != nil {
		return nil, 0, err
	}

	return res, count, nil
}

func AddConnection(DictionaryId int, WordId int) error {
	dbCon := db.GetConnection()

	_, err := dbCon.Exec("insert into DictionariesWords values (?, ?)", DictionaryId, WordId)

	if err != nil {
		return err
	}

	return nil
}
