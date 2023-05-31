package dictionaryRoutes

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/RubyLegend/dictionary-backend/middleware/cors"
	userHelper "github.com/RubyLegend/dictionary-backend/middleware/users"
	dictionaryRepo "github.com/RubyLegend/dictionary-backend/repository/dictionary"
	dictionaryToWordsRepo "github.com/RubyLegend/dictionary-backend/repository/dictionaryToWords"
	userRepo "github.com/RubyLegend/dictionary-backend/repository/users"
	wordRepo "github.com/RubyLegend/dictionary-backend/repository/words"
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

type dictionaryRepoWrapper interface {
	GetDictionary(int) ([]dictionaryRepo.Dictionary, error)
	AddDictionary(int, dictionaryRepo.DictionaryPost) []error
	DeleteDictionary(int, int) error
	UpdateDictionary(int, int, dictionaryRepo.Dictionary) (dictionaryRepo.Dictionary, error)
}

type dictionaryRepoWrap struct{}

func (d dictionaryRepoWrap) GetDictionary(UserId int) ([]dictionaryRepo.Dictionary, error) {
	return dictionaryRepo.GetDictionary(UserId)
}

func (d dictionaryRepoWrap) AddDictionary(UserId int, dictionaryData dictionaryRepo.DictionaryPost) []error {
	return dictionaryRepo.AddDictionary(UserId, dictionaryData)
}

func (d dictionaryRepoWrap) DeleteDictionary(UserId int, DictionaryId int) error {
	return dictionaryRepo.DeleteDictionary(UserId, DictionaryId)
}

func (d dictionaryRepoWrap) UpdateDictionary(UserId int, DictionaryId int, dictionaryData dictionaryRepo.Dictionary) (dictionaryRepo.Dictionary, error) {
	return dictionaryRepo.UpdateDictionary(UserId, DictionaryId, dictionaryData)
}

type dictionaryToWordsRepoWrapper interface {
	GetWords(int, int, int) ([]dictionaryToWordsRepo.DictionaryToWords, int, error)
}

type wordRepoWrapper interface {
	WordIDtoWords([]dictionaryToWordsRepo.DictionaryToWords) ([]wordRepo.Word, error)
}

type dtwrWrap struct{}

func (d dtwrWrap) GetWords(DictionaryId int, page int, limit int) ([]dictionaryToWordsRepo.DictionaryToWords, int, error) {
	return dictionaryToWordsRepo.GetWords(DictionaryId, page, limit)
}

type wrWrap struct{}

func (w wrWrap) WordIDtoWords(dictToWords []dictionaryToWordsRepo.DictionaryToWords) ([]wordRepo.Word, error) {
	return wordRepo.WordIDtoWords(dictToWords)
}

var (
	userHelp  userHelperWrapper
	userRepoW userRepoWrapper
	drW       dictionaryRepoWrapper
	dtwr      dictionaryToWordsRepoWrapper
	wr        wordRepoWrapper
)

func init() {
	userHelp = userHelpWrap{}
	userRepoW = userRepoWrap{}
	drW = dictionaryRepoWrap{}
	dtwr = dtwrWrap{}
	wr = wrWrap{}
}

func DictionaryGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
			dictionary, err := drW.GetDictionary(userData.UserId)

			if err != nil {
				resp["error"] = []string{err.Error()}
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(resp)
			} else {
				resp["dictionary"] = dictionary
				_ = json.NewEncoder(w).Encode(dictionary)
			}
		}
	} else {
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func DictionaryGetWords(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	resp := make(map[string]any)

	claims := userHelp.VerifyJWT(w, r, resp)
	if resp["error"] == nil {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))

		if err != nil {
			page = 0
		} else {
			page = page - 1
		}

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

		if err != nil {
			limit = int(^uint(0) >> 1) // Magic value for upper limit of integer in Go
		}

		var userData userRepo.User
		userData.Username = claims["username"].(string)
		userData, err = userRepoW.GetUser(userData)

		if err != nil {
			resp["error"] = []string{err.Error()}
			w.WriteHeader(http.StatusNotFound)
		} else {
			DictionaryId, err := strconv.Atoi(ps.ByName("id"))

			if err != nil {
				resp["error"] = []string{err.Error()}
				w.WriteHeader(http.StatusBadRequest)
			} else {
				wordIds, count, err := dtwr.GetWords(DictionaryId, page, limit)

				if err != nil {
					resp["error"] = []string{err.Error()}
					w.WriteHeader(http.StatusBadRequest)
				} else {
					words, err := wr.WordIDtoWords(wordIds)
					if err != nil {
						resp["error"] = []string{err.Error()}
						w.WriteHeader(http.StatusInternalServerError)
					} else {
						resp["words"] = words
						resp["count"] = count
						resp["limit"] = limit
						resp["page"] = page + 1
						resp["pages"] = math.Ceil(float64(count) / float64(limit))
					}
				}
			}
		}
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func DictionaryPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	var dictionaryData dictionaryRepo.DictionaryPost
	_ = json.NewDecoder(r.Body).Decode(&dictionaryData)
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
			err := drW.AddDictionary(userData.UserId, dictionaryData)
			if err != nil {
				var errors []string
				for _, v := range err {
					errors = append(errors, v.Error())
				}
				resp["error"] = errors
				w.WriteHeader(http.StatusBadRequest)
			} else {
				resp["status"] = "success"
				w.WriteHeader(http.StatusOK)
			}
		}
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func DictionaryPatch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
			DictionaryId, err := strconv.Atoi(ps.ByName("id"))

			if err != nil {
				resp["error"] = []string{err.Error()}
				w.WriteHeader(http.StatusInternalServerError)
			} else {

				var dictionaryData dictionaryRepo.Dictionary
				_ = json.NewDecoder(r.Body).Decode(&dictionaryData)

				dict, errors := drW.UpdateDictionary(userData.UserId, DictionaryId, dictionaryData)

				if errors != nil {
					resp["error"] = []string{errors.Error()}
					w.WriteHeader(http.StatusBadRequest)
				} else {
					resp["dictionary"] = dict
					w.WriteHeader(http.StatusOK)
				}
			}
		}
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func DictionaryDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
			DictionaryId, err := strconv.Atoi(ps.ByName("id"))

			if err != nil {
				resp["error"] = []string{err.Error()}
				w.WriteHeader(http.StatusInternalServerError)
			} else {

				errors := drW.DeleteDictionary(userData.UserId, DictionaryId)

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
