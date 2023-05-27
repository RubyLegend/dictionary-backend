package historyRoutes

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func HistoryDelete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Not Implemented\n")
}
