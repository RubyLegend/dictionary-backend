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
	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/julienschmidt/httprouter"
)

type dictionaryToWordsRepoWrapper interface {
	GetWords(int, int, int) ([]dictionaryToWordsRepo.DictionaryToWords, int, error)
	AddConnection(int, int) error
}

type wordRepoWrapper interface {
	AddWord(wordRepo.Word) (int, wordRepo.Word, error)
	UpdateWord(wordRepo.Word) error
	DeleteWord(wordRepo.Word) error
}

type dtwrWrap struct{}

func (d dtwrWrap) GetWords(DictionaryId int, page int, limit int) ([]dictionaryToWordsRepo.DictionaryToWords, int, error) {
	return dictionaryToWordsRepo.GetWords(DictionaryId, page, limit)
}

func (d dtwrWrap) AddConnection(DictionaryId int, WordId int) error {
	return dictionaryToWordsRepo.AddConnection(DictionaryId, WordId)
}

type wrWrap struct{}

func (w wrWrap) AddWord(wordData wordRepo.Word) (int, wordRepo.Word, error) {
	return wordRepo.AddWord(wordData)
}

func (w wrWrap) UpdateWord(wordData wordRepo.Word) error {
	return wordRepo.UpdateWord(wordData)
}

func (w wrWrap) DeleteWord(wordData wordRepo.Word) error {
	return wordRepo.DeleteWord(wordData)
}

var dtwr dictionaryToWordsRepoWrapper
var wr wordRepoWrapper

type userHelperWrapper interface {
	VerifyJWT(http.ResponseWriter, *http.Request, map[string]any) jwt.MapClaims
}

type userHelpWrap struct{}

func (u userHelpWrap) VerifyJWT(w http.ResponseWriter, r *http.Request, resp map[string]any) jwt.MapClaims {
	return userHelper.VerifyJWT(w, r, resp)
}

var (
	userHelperW userHelperWrapper
)

func init() {
	dtwr = dtwrWrap{}
	wr = wrWrap{}
	userHelperW = userHelpWrap{}
}

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

	_ = userHelperW.VerifyJWT(w, r, resp)
	if resp["error"] == nil {
		word := wordData.ConvertToWord()
		lastId, word, err := wr.AddWord(word)

		if err != nil {
			resp["error"] = []string{err.Error()}
			w.WriteHeader(http.StatusNotAcceptable)
		} else {
			err = dtwr.AddConnection(wordData.DictionaryId, lastId)

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

	resp := make(map[string]any)

	_ = userHelperW.VerifyJWT(w, r, resp)
	if resp["error"] == nil {
		wordId, err := strconv.Atoi(ps.ByName("id"))

		if err != nil {
			resp["error"] = []string{err.Error()}
			w.WriteHeader(http.StatusNotAcceptable)
		} else {
			wordData.WordId = wordId
			err := wr.DeleteWord(wordData)

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

	_ = userHelperW.VerifyJWT(w, r, resp)
	if resp["error"] == nil {
		wordId, err := strconv.Atoi(ps.ByName("id"))

		if err != nil {
			resp["error"] = []string{err.Error()}
			w.WriteHeader(http.StatusNotAcceptable)
		} else {
			wordData.WordId = wordId
			err := wr.UpdateWord(wordData)

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
