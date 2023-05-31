package historyRoutes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RubyLegend/dictionary-backend/middleware/cors"
	userHelper "github.com/RubyLegend/dictionary-backend/middleware/users"
	historyRepo "github.com/RubyLegend/dictionary-backend/repository/history"
	userRepo "github.com/RubyLegend/dictionary-backend/repository/users"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

type userHelperWrapper interface {
	VerifyJWT(http.ResponseWriter, *http.Request, map[string]any) jwt.MapClaims
}

type userHelpWrap struct{}

func (u userHelpWrap) VerifyJWT(w http.ResponseWriter, r *http.Request, resp map[string]any) jwt.MapClaims {
	return userHelper.VerifyJWT(w, r, resp)
}

type userRepoWrapper interface {
	GetUser(userRepo.User) (userRepo.User, error)
}

type userRepoWrap struct{}

func (u userRepoWrap) GetUser(userData userRepo.User) (userRepo.User, error) {
	return userRepo.GetUser(userData)
}

type historyRepoWrapper interface {
	GetHistory(int) ([]historyRepo.History, error)
	DeleteHistory(int, int) error
}

type historyRepoWrap struct{}

func (h historyRepoWrap) GetHistory(UserId int) ([]historyRepo.History, error) {
	return historyRepo.GetHistory(UserId)
}

func (h historyRepoWrap) DeleteHistory(UserId int, HistoryId int) error {
	return historyRepo.DeleteHistory(UserId, HistoryId)
}

var (
	userHelp     userHelperWrapper
	userRepoW    userRepoWrapper
	historyRepoW historyRepoWrapper
)

func init() {
	userHelp = userHelpWrap{}
	userRepoW = userRepoWrap{}
	historyRepoW = historyRepoWrap{}
}

func HistoryGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	resp := make(map[string]any)

	claims := userHelp.VerifyJWT(w, r, resp)
	if resp["error"] == nil {
		var userData userRepo.User
		userData.Username = claims["username"].(string)
		userData, err := userRepoW.GetUser(userData)

		if err != nil {
			resp["error"] = []string{err.Error()}
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(resp)
		} else {
			history, err := historyRepoW.GetHistory(userData.UserId)

			if err != nil {
				resp["error"] = []string{err.Error()}
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(resp)
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

	claims := userHelp.VerifyJWT(w, r, resp)
	if resp["error"] == nil {
		var userData userRepo.User
		userData.Username = claims["username"].(string)
		userData, err := userRepoW.GetUser(userData)

		if err != nil {
			resp["error"] = []string{err.Error()}
			w.WriteHeader(http.StatusNotFound)
		} else {
			HistoryId, err := strconv.Atoi(ps.ByName("id"))

			if err != nil {
				resp["error"] = []string{err.Error()}
				w.WriteHeader(http.StatusNotAcceptable)
			} else {

				errors := historyRepoW.DeleteHistory(userData.UserId, HistoryId)

				if errors != nil {
					resp["error"] = []string{errors.Error()}
					w.WriteHeader(http.StatusBadRequest)
				} else {
					resp["status"] = "success"
					w.WriteHeader(http.StatusOK)
				}
			}
		}
	}
	_ = json.NewEncoder(w).Encode(resp)
}
