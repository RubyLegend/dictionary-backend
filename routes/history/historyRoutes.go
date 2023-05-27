package historyRoutes

import (
	"encoding/json"
	"net/http"
	"strconv"

	historyRepo "github.com/RubyLegend/dictionary-backend/repository/history"
	"github.com/julienschmidt/httprouter"
)

func HistoryGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var historyId int = 1

	resp := make(map[string]interface{})

	errors, history := historyRepo.GetHistory(historyId)

	if len(errors) > 0 {
		var errorMessages []string
		for _, v := range errors {
			errorMessages = append(errorMessages, v.Error())
		}
		resp["error"] = errorMessages
		w.WriteHeader(http.StatusBadRequest)
	} else {
		resp["history"] = history
		w.WriteHeader(http.StatusOK)
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func HistoryDelete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]interface{})

	var UserId int = 1
	HistoryId := p.ByName("id")

	historyId, err := strconv.Atoi(HistoryId)
	if err != nil {
		resp["error"] = []string{"Invalid history ID"}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	errors := historyRepo.DeleteHistory(UserId, historyId)

	if len(errors) > 0 {
		var errorMessages []string
		for _, v := range errors {
			errorMessages = append(errorMessages, v.Error())
		}
		resp["error"] = errorMessages
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_ = json.NewEncoder(w).Encode(resp)
}
