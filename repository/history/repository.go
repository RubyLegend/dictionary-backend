package histories

import (
	"errors"
	"time"
)

type History struct {
	HistoryId int       `json:"HistoryId"`
	UserId    int       `json:"UserId"`
	WordId    int       `json:"WordId"`
	isCorrect bool      //`json:"isCorrect"`
	createdAt time.Time //`json:"createdAt"`
}

var Histories []History

//func checkHistoryExistance(historyData History) []error {
//	var err []error
//
//	for _, v := range Histories {
//		if v.isCorrect == historyData.isCorrect && v.UserId == historyData.UserId {
//			err = append(err, errors.New("History  already exists"))
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

func GetHistory(historyId int) ([]error, History) {
	var err []error
	Histories = append(Histories, History{HistoryId: 1, UserId: 1, WordId: 1, isCorrect: true, createdAt: time.Now()})
	var foundHistory History

	for _, v := range Histories {
		if v.HistoryId == historyId {
			foundHistory = v
			break
		}
	}

	if (foundHistory == History{}) {
		err = append(err, errors.New("History not found"))
	}

	return err, foundHistory
}
func DeleteHistory(userId, historyId int) []error {
	var err []error
	var isDeleted = false

	for i, v := range Histories {
		if v.UserId == userId && v.HistoryId == historyId {
			Histories = append(Histories[:i], Histories[i+1:]...)
			isDeleted = true
			break
		}
	}

	if !isDeleted {
		err = append(err, errors.New("History not found"))
	}

	return err
}
