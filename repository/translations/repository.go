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

func checkTranslationExistance(translation Translation) error {
	dbCon := db.GetConnection()
	var res int
	err := dbCon.QueryRow("select count(*) from Translation where binary Name = binary ? and wordID = ?", translation.Name, translation.WordId).Scan(&res)

	if err != nil {
		return err
	}

	return nil
}

func GetTranslation(WordId int) ([]Translation, error) {
	var Translations []Translation
	dbCon := db.GetConnection()

	rows, err := dbCon.Query("select * from Translation where wordID = ?", WordId)

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

func AddTranslation(WordId int, translation string) error {
	var translationData Translation
	translationData.WordId = WordId
	translationData.Name = translation

	err := postValidation(translationData)

	if err != nil {
		return err
	}

	err = checkTranslationExistance(translationData)

	if err != nil {
		return err
	}

	dbCon := db.GetConnection()

	_, err = dbCon.Exec("insert into Translation values (default, ?, ?, default)", translationData.WordId, translationData.Name)

	if err != nil {
		return err
	}

	return nil
}

func updateValidation(translationData Translation) error {
	if len(translationData.Name) == 0 {
		return errors.New("name not supplied")
	}
	// if len(translationData.Language) == 0 {
	// 	return errors.New("new dictionary language not supplied")
	// }

	return nil
}

func UpdateTranslation(translationData Translation) error {
	err := updateValidation(translationData)

	if err != nil {
		return err
	}

	err = checkTranslationExistance(translationData)
	if err != nil {
		return err
	}

	dbCon := db.GetConnection()

	_, err = dbCon.Exec("update Translation set name = ? where wordID = ?", translationData.Name, translationData.WordId)

	if err != nil {
		return err
	}

	return nil
}

func DeleteTranslation(translationData Translation) error {
	dbCon := db.GetConnection()

	res, err := dbCon.Exec("delete from Translation where wordID = ?", translationData.WordId)

	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.New("this translation doesn't exist")
	}

	return nil
}
