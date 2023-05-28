package dictionary

import (
	"errors"
	"time"

	db "github.com/RubyLegend/dictionary-backend/middleware/database"
)

type Dictionary struct {
	DictionaryId int       `json:"id"`
	UserId       int       `json:"userId"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"createdAt"`
	Total        int       `json:"total"`
}

var Dictionaries []Dictionary

func checkDictionaryExistance(dictionaryData Dictionary) error {
	dbCon := db.GetConnection()
	var res int
	err := dbCon.QueryRow("select count(*) from Dictionaries where Name = ? and userID = ?", dictionaryData.Name, dictionaryData.UserId).Scan(&res)

	if err != nil {
		return err
	}

	if res != 0 {
		return errors.New("dictionary already exist")
	}

	return nil
}

func postValidation(dictionaryData Dictionary) []error {
	var err []error

	if len(dictionaryData.Name) == 0 {
		err = append(err, errors.New("name is required field"))
	}

	dbCon := db.GetConnection()

	var res int
	err2 := dbCon.QueryRow("select count(*) from Users where userID = ?", dictionaryData.UserId).Scan(&res)

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
		rows.Scan(&dict.DictionaryId, &dict.UserId, &dict.Name, &dict.CreatedAt, &dict.Total)
		Dictionaries = append(Dictionaries, dict)
	}

	if len(Dictionaries) == 0 {
		return []Dictionary{}, nil
	}

	return Dictionaries, nil
}

func AddDictionary(dictionaryData Dictionary) []error {
	var err []error

	err2 := postValidation(dictionaryData)
	if err2 != nil {
		err = append(err, err2...)
	}

	err3 := checkDictionaryExistance(dictionaryData)
	if err3 != nil {
		err = append(err, err3)
	}

	if err != nil {
		return err
	}

	dbCon := db.GetConnection()

	_, err3 = dbCon.Exec("insert into Dictionaries values (default, ?, ?, CURRENT_TIMESTAMP(), default)", dictionaryData.UserId, dictionaryData.Name)

	if err3 != nil {
		err = append(err, err3)
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

	var dictWithUserId Dictionary
	dictWithUserId.UserId = UserId
	dictWithUserId.Name = dictionaryData.Name

	err = checkDictionaryExistance(dictWithUserId)
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

	row.Scan(&dict.DictionaryId, &dict.UserId, &dict.Name, &dict.CreatedAt, &dict.Total)

	return dict, nil
}
