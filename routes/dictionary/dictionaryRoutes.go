package dictionaryRoutes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RubyLegend/dictionary-backend/middleware/cors"
	userHelper "github.com/RubyLegend/dictionary-backend/middleware/users"
	dictionaryRepo "github.com/RubyLegend/dictionary-backend/repository/dictionary"
	dictionaryToWordsRepo "github.com/RubyLegend/dictionary-backend/repository/dictionaryToWords"
	userRepo "github.com/RubyLegend/dictionary-backend/repository/users"
	wordsRepo "github.com/RubyLegend/dictionary-backend/repository/words"
	"github.com/julienschmidt/httprouter"
)

func DictionaryGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
			dictionary, err := dictionaryRepo.GetDictionary(userData.UserId)

			if err != nil {
				resp["error"] = []string{err.Error()}
				_ = json.NewEncoder(w).Encode(resp)
				w.WriteHeader(http.StatusBadRequest)
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

	claims := userHelper.VerifyJWT(w, r, resp)
	if resp["error"] == nil {
		var userData userRepo.User
		userData.Username = claims["username"].(string)
		userData, err := userRepo.GetUser(userData)

		if err != nil {
			resp["error"] = []string{err.Error()}
			w.WriteHeader(http.StatusNotFound)
		} else {
			DictionaryId, err := strconv.Atoi(ps.ByName("id"))

			if err != nil {
				resp["error"] = []string{err.Error()}
			} else {
				wordIds, err := dictionaryToWordsRepo.GetWords(DictionaryId)

				if err != nil {
					resp["error"] = []string{err.Error()}
					w.WriteHeader(http.StatusBadRequest)
				} else {
					words, err := wordsRepo.WordIDtoWords(wordIds)
					if err != nil {
						resp["error"] = []string{err.Error()}
						w.WriteHeader(http.StatusInternalServerError)
					} else {
						resp["words"] = words
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
	var dictionaryData dictionaryRepo.Dictionary
	_ = json.NewDecoder(r.Body).Decode(&dictionaryData)
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
			dictionaryData.UserId = userData.UserId

			err := dictionaryRepo.AddDictionary(dictionaryData)
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

	claims := userHelper.VerifyJWT(w, r, resp)
	if resp["error"] == nil {
		var userData userRepo.User
		userData.Username = claims["username"].(string)
		userData, err := userRepo.GetUser(userData)

		if err != nil {
			resp["error"] = []string{err.Error()}
			w.WriteHeader(http.StatusNotFound)
		} else {
			DictionaryId, err := strconv.Atoi(ps.ByName("id"))

			if err != nil {
				resp["error"] = []string{err.Error()}
			} else {

				var dictionaryData dictionaryRepo.Dictionary
				_ = json.NewDecoder(r.Body).Decode(&dictionaryData)

				dict, errors := dictionaryRepo.UpdateDictionary(userData.UserId, DictionaryId, dictionaryData)

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

	claims := userHelper.VerifyJWT(w, r, resp)
	if resp["error"] == nil {
		var userData userRepo.User
		userData.Username = claims["username"].(string)
		userData, err := userRepo.GetUser(userData)

		if err != nil {
			resp["error"] = []string{err.Error()}
			w.WriteHeader(http.StatusNotFound)
		} else {
			DictionaryId, err := strconv.Atoi(ps.ByName("id"))

			if err != nil {
				resp["error"] = []string{err.Error()}
			} else {

				errors := dictionaryRepo.DeleteDictionary(userData.UserId, DictionaryId)

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
