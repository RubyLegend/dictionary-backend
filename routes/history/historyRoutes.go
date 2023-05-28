package historyRoutes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RubyLegend/dictionary-backend/middleware/cors"
	userHelper "github.com/RubyLegend/dictionary-backend/middleware/users"
	historyRepo "github.com/RubyLegend/dictionary-backend/repository/history"
	userRepo "github.com/RubyLegend/dictionary-backend/repository/users"
	"github.com/julienschmidt/httprouter"
)

func HistoryGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	resp := make(map[string]any)

	claims := userHelper.VerifyJWT(w, r, resp)
	if resp["error"] == nil {
		var userData userRepo.User
		userData.Username = claims["username"].(string)
		userData, err := userRepo.GetUser(userData)

		if err != nil {
			resp["error"] = []string{err.Error()}
			_ = json.NewEncoder(w).Encode(resp)
			w.WriteHeader(http.StatusNotFound)
		} else {
			history, err := historyRepo.GetHistory(userData.UserId)

			if err != nil {
				resp["error"] = []string{err.Error()}
				_ = json.NewEncoder(w).Encode(resp)
				w.WriteHeader(http.StatusBadRequest)
			} else {
				_ = json.NewEncoder(w).Encode(history)
			}
		}
	} else {
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func HistoryDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	resp := make(map[string]any)

	claims := userHelper.VerifyJWT(w, r, resp)
	if resp["error"] == nil {
		var userData userRepo.User
		userData.Username = claims["username"].(string)
		userData, err := userRepo.GetUser(userData)

		if err != nil {
			resp["error"] = []string{err.Error()}
			w.WriteHeader(http.StatusNotFound)
		} else {
			HistoryId, err := strconv.Atoi(ps.ByName("id"))

			if err != nil {
				resp["error"] = []string{err.Error()}
			} else {

				errors := historyRepo.DeleteHistory(userData.UserId, HistoryId)

				if errors != nil {
					resp["error"] = []string{errors.Error()}
					w.WriteHeader(http.StatusBadRequest)
				} else {
					w.WriteHeader(http.StatusOK)
				}
			}
		}
	}
	_ = json.NewEncoder(w).Encode(resp)
}
