package history

import (
	"time"

	db "github.com/RubyLegend/dictionary-backend/middleware/database"
)

type History struct {
	HistoryId int       `json:"HistoryId"`
	UserId    int       `json:"UserId"`
	WordId    int       `json:"WordId"`
	IsCorrect bool      `json:"isCorrect"`
	CreatedAt time.Time `json:"createdAt"`
}

var Histories []History

//func checkHistoryExistance(historyData History) []error {
//	var err []error
//
//	for _, v := range Histories {
//		if v.isCorrect == historyData.isCorrect && v.UserId == historyData.UserId {
//			err = append(err, errors.New("history  already exists"))
//		}
//	}
//
//	return err
//}

//func findHistory(histories []History, historyID int) (History, error) {
//	for _, history := range histories {
//		if history.HistoryId == historyID {
//			return history, nil
//		}
//	}
//	return History{}, fmt.Errorf("history with ID %d not found", historyID)
//}

func GetHistory(UserId int) ([]History, error) {
	var Histories []History
	dbCon := db.GetConnection()

	rows, err := dbCon.Query("select * from Histories where userID = ?", UserId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var hist History
		err = rows.Scan(&hist.HistoryId, &hist.UserId, &hist.WordId, &hist.IsCorrect, &hist.CreatedAt)
		if err != nil {
			return nil, err
		}
		Histories = append(Histories, hist)
	}

	if len(Histories) == 0 {
		return []History{}, nil
	}

	return Histories, nil
}
func DeleteHistory(UserId int, HistoryId int) error {
	dbCon := db.GetConnection()

	_, error := dbCon.Exec("delete from Histories where historyID = ? and userID = ?", HistoryId, UserId)

	if error != nil {
		return error
	}

	return nil
}
