package translations

import (
	"errors"

	
	db "github.com/RubyLegend/dictionary-backend/middleware/database"
)

type Translation struct {
	TranslationId int    `json:"translaonId"`
	WordId        int    `json:"wordId"`
	Name          string `json:"name"`
	Language      string `json:"language"`
}
type TranslationWithWordId struct {
	Name         string `json:"name"`
	DictionaryId int    `json:"dictionaryId"`
}

func checkTranslationExistance(WordId int, Name string) error {
	dbCon := db.GetConnection()
	var res int
	err := dbCon.QueryRow("select count(*) from Translations where binary Name = binary ? and wordID = ?", Name, WordId).Scan(&res)

	if err != nil {
		return err
	}

	return nil
}

func validation(translationData Translation) []error {
	var err []error

	if len(translationData.Name) == 0 {
		err = append(err, errors.New("name is required field"))
	}
	
	dbCon := db.GetConnection()

	var res int
	err2 := dbCon.QueryRow("select count(*) from Words where wordID = ?", WordId).Scan(&res)

	if err2 != nil {
		err = append(err, err2)
		return err
	}

	if res == 0 {
		err = append(err, errors.New("translation not found"))
	}

	return err
}

func GetTranslation(WordId int) ([]Translation, error) {
	var Translations []Translation
	dbCon := db.GetConnection()

	rows, err := dbCon.Query("select * from Translations where wordID = ?", WordId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var trans Translation
		err = rows.Scan(&trans.TranslationId, &trans.WordId, &trans.Name, &trans.Language)
		if err != nil {
			return nil, err
		}
		Translations = append(Translations, trans)
	}

	if len(Translations) == 0 {
		return []Translation{}, nil
	}

	return Translations, nil
}
func postValidation(translationData Translation) error {

	if len(translationData.Name) == 0 {
		return errors.New("name is required field")
	}

	return nil
}

func AddTranslation(translationData Translation) (int, Translation,error) {
	
	err := postValidation(translationData)

	if err != nil {
		return -1, Translation{}, err
	}

	err = checkTranslationExistance(translationData)

	if err != nil {
		return -1, Translation{}, err
	}

	dbCon := db.GetConnection()

	res, err := dbCon.Exec("insert into Words values (default, ?, CURRENT_TIMESTAMP(), default)", wordData.Name)

	if err != nil {
		return -1, Word{}, err
	}

	lastId, err := res.LastInsertId()

	wordData.WordId = int(lastId)
	wordData.CreatedAt = time.Now()

	return int(lastId), wordData, err
}

func updateValidation(translationData Translation) error {
	if len(translationData.Name) == 0 {
		return errors.New("new dictionary name not supplied")
	}
	if len(translationData.Language) == 0 {
		return errors.New("new dictionary language not supplied")
	}

	return nil
}

func UpdateTranslation(translationData Translation)error {
	err := updateValidation(translationData)

	if err != nil {
		return Translation{}, err
	}

	err = checkTranslationExistence(WordId, translationData.Name)
	if err != nil {
		return Translation{}, err
	}

	dbCon := db.GetConnection()

	_, err = dbCon.Exec("update Translation set name = ?, language = ? where translationID = ?", translationData.Name, TranslationId)

	if err != nil {
		return Translation{}, err
	}

	row := dbCon.QueryRow("select * from Translations where translationID = ?", TranslationId)

	var trans Translation

	err = row.Scan(&trans.TranslationId, &trans.WordId, &trans.Name, &trans.Language)

	if err != nil {
		return Translation{}, err
	}

	return trans, nil
}

func DeleteTranslation(translationData Translation) error {
	dbCon := db.GetConnection()

	res, err:= dbCon.Exec("delete from Translations where translationID = ? and wordID = ?", TranslationId, WordId)

	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.New("this translation doesn't exist")
	}

	return nil
}