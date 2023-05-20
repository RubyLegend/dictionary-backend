package dictionaryRoutes

import (
	"encoding/json"
	"fmt"
	"net/http"

	dictionaryRepo "github.com/RubyLegend/dictionary-backend/repository/dictionary"
	"github.com/julienschmidt/httprouter"
)

func DictionaryGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Not Implemented\n")
}

func DictionaryPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var UserId = 1

	var dictionaryData dictionaryRepo.Dictionary
	_ = json.NewDecoder(r.Body).Decode(&dictionaryData)

	resp := make(map[string]any)
	dictionaryData.UserId = UserId

	err := dictionaryRepo.AddDictionary(dictionaryData)
	if err != nil {
		var errors []string
		for _, v := range err {
			errors = append(errors, v.Error())
		}
		resp["error"] = errors
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
  _ = json.NewEncoder(w).Encode(resp)
}

func DictionaryPatch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Not Implemented\n")
}

func DictionaryDelete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Not Implemented\n")
}
