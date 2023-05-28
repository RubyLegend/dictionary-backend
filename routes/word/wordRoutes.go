package wordRoutes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/RubyLegend/dictionary-backend/middleware/cors"
	httpHelper "github.com/RubyLegend/dictionary-backend/middleware/httphelper"
	userHelper "github.com/RubyLegend/dictionary-backend/middleware/users"
	dictionaryToWordsRepo "github.com/RubyLegend/dictionary-backend/repository/dictionaryToWords"
	wordRepo "github.com/RubyLegend/dictionary-backend/repository/words"

	"github.com/julienschmidt/httprouter"
)

func WordGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	fmt.Fprintf(w, "Not Implemented\n")
}

func WordPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)

	var wordData wordRepo.WordWithDictId
	_ = json.NewDecoder(r.Body).Decode(&wordData)

	resp := make(map[string]any)

	_ = userHelper.VerifyJWT(w, r, resp)
	if resp["error"] == nil {
		word := wordData.ConvertToWord()
		lastId, word, err := wordRepo.AddWord(word)

		if err != nil {
			resp["error"] = []string{err.Error()}
			w.WriteHeader(http.StatusNotAcceptable)
		} else {
			err = dictionaryToWordsRepo.AddConnection(wordData.DictionaryId, lastId)

			if err != nil {
				resp["error"] = []string{err.Error()}
				w.WriteHeader(http.StatusNotAcceptable)
			} else {
				resp["status"] = "success"
				httpHelper.UnpackToResp(word, resp)
			}
		}
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func WordDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)

	var wordData wordRepo.Word
	_ = json.NewDecoder(r.Body).Decode(&wordData)

	resp := make(map[string]any)

	_ = userHelper.VerifyJWT(w, r, resp)
	if resp["error"] == nil {
		wordId, err := strconv.Atoi(ps.ByName("id"))

		if err != nil {
			resp["error"] = []string{err.Error()}
			w.WriteHeader(http.StatusNotAcceptable)
		} else {
			wordData.WordId = wordId
			err := wordRepo.DeleteWord(wordData)

			if err != nil {
				resp["error"] = []string{err.Error()}
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				resp["status"] = "success"
			}
		}
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func WordPatch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)

	var wordData wordRepo.Word
	_ = json.NewDecoder(r.Body).Decode(&wordData)

	resp := make(map[string]any)

	_ = userHelper.VerifyJWT(w, r, resp)
	if resp["error"] == nil {
		wordId, err := strconv.Atoi(ps.ByName("id"))

		if err != nil {
			resp["error"] = []string{err.Error()}
			w.WriteHeader(http.StatusNotAcceptable)
		} else {
			wordData.WordId = wordId
			err := wordRepo.UpdateWord(wordData)

			if err != nil {
				resp["error"] = []string{err.Error()}
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				resp["status"] = "success"
			}
		}
	}

	_ = json.NewEncoder(w).Encode(resp)
}
