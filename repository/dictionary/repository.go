package dictionary

import (
	"errors"
	"log"
	"time"

	db "github.com/RubyLegend/dictionary-backend/middleware/database"
	dictionaryToWordsRepo "github.com/RubyLegend/dictionary-backend/repository/dictionaryToWords"
	wordsRepo "github.com/RubyLegend/dictionary-backend/repository/words"
)

type Dictionary struct {
	DictionaryId int       `json:"id"`
	UserId       int       `json:"userId"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"createdAt"`
	Total        int       `json:"total"`
}

type DictionaryPost struct {
	Name  string           `json:"dictionaryName"`
	Words []wordsRepo.Word `json:"words"`
}

var Dictionaries []Dictionary

func CheckDictionaryExistance(UserId int, Name string) error {
	dbCon := db.GetConnection()
	var res int
	err := dbCon.QueryRow("select count(*) from Dictionaries where binary Name = binary ? and userID = ?", Name, UserId).Scan(&res)

	if err != nil {
		return err
	}

	if res != 0 {
		return errors.New("dictionary already exist")
	}

	return nil
}

func postValidation(UserId int, dictionaryData DictionaryPost) []error {
	var err []error

	if len(dictionaryData.Name) == 0 {
		err = append(err, errors.New("name is required field"))
	}

	dbCon := db.GetConnection()

	var res int
	err2 := dbCon.QueryRow("select count(*) from Users where userID = ?", UserId).Scan(&res)

	if err2 != nil {
		err = append(err, err2)
		return err
	}

	if res == 0 {
		err = append(err, errors.New("dictionary owner not found"))
	}

	return err
}

func GetDictionary(UserId int) ([]Dictionary, error) {
	var Dictionaries []Dictionary
	dbCon := db.GetConnection()

	rows, err := dbCon.Query("select * from Dictionaries where userID = ?", UserId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var dict Dictionary
		err = rows.Scan(&dict.DictionaryId, &dict.UserId, &dict.Name, &dict.CreatedAt, &dict.Total)
		if err != nil {
			return nil, err
		}
		Dictionaries = append(Dictionaries, dict)
	}

	if len(Dictionaries) == 0 {
		return []Dictionary{}, nil
	}

	return Dictionaries, nil
}

func AddDictionary(UserId int, dictionaryData DictionaryPost) []error {
	err := postValidation(UserId, dictionaryData)
	if err != nil {
		return err
	}

	err2 := CheckDictionaryExistance(UserId, dictionaryData.Name)
	if err2 != nil {
		return []error{err2}
	}

	dbCon := db.GetConnection()

	res, err2 := dbCon.Exec("insert into Dictionaries values (default, ?, ?, CURRENT_TIMESTAMP(), default)", UserId, dictionaryData.Name)

	if err2 != nil {
		return []error{err2}
	}

	lastDictionaryId, err2 := res.LastInsertId()

	if err2 != nil {
		return []error{err2}
	}

	for i, v := range dictionaryData.Words {
		lastId, _, err2 := wordsRepo.AddWord(v)
		if err2 != nil {
			return []error{err2}
		}

		err2 = dictionaryToWordsRepo.AddConnection(int(lastDictionaryId), lastId)

		if err2 != nil {
			return []error{err2}
		}

		log.Println(i, v)
	}

	return err
}

func DeleteDictionary(UserId int, DictionaryId int) error {
	dbCon := db.GetConnection()

	_, error := dbCon.Exec("delete from Dictionaries where dictionaryID = ? and userID = ?", DictionaryId, UserId)

	if error != nil {
		return error
	}

	return nil
}

func updateValidation(dictionaryData Dictionary) error {
	if len(dictionaryData.Name) == 0 {
		return errors.New("new dictionary name not supplied")
	}

	return nil
}

func UpdateDictionary(UserId int, DictionaryId int, dictionaryData Dictionary) (Dictionary, error) {
	err := updateValidation(dictionaryData)

	if err != nil {
		return Dictionary{}, err
	}

	err = CheckDictionaryExistance(UserId, dictionaryData.Name)
	if err != nil {
		return Dictionary{}, err
	}

	dbCon := db.GetConnection()

	_, err = dbCon.Exec("update Dictionaries set name = ? where dictionaryID = ?", dictionaryData.Name, DictionaryId)

	if err != nil {
		return Dictionary{}, err
	}

	row := dbCon.QueryRow("select * from Dictionaries where dictionaryID = ?", DictionaryId)

	var dict Dictionary

	err = row.Scan(&dict.DictionaryId, &dict.UserId, &dict.Name, &dict.CreatedAt, &dict.Total)

	if err != nil {
		return Dictionary{}, err
	}

	return dict, nil
}
