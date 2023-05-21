package dictionaryRoutes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	dictionaryRepo "github.com/RubyLegend/dictionary-backend/repository/dictionary"
	"github.com/julienschmidt/httprouter"
)

func DictionaryGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err []error
	w.Header().Set("Content-Type", "application/json")
	var UserId int = 1
	DictionaryIdStr := ps.ByName("id")

	DictionaryId, errAtoi := strconv.Atoi(DictionaryIdStr)
	if errAtoi != nil {
		err = append(err, errors.New("Occurred some error"))
	}

	resp := make(map[string]any)

	errors, dictionary := dictionaryRepo.GetDictionary(UserId, DictionaryId)

	if errors != nil {
		err = append(err, errors...)
	}

	if len(errors) > 0 {
		var errorMessages []string
		for _, v := range errors {
			errorMessages = append(errorMessages, v.Error())
		}
		resp["error"] = errorMessages
		w.WriteHeader(http.StatusBadRequest)
	} else {
		resp["dictionary"] = dictionary
		w.WriteHeader(http.StatusOK)
	}
	_ = json.NewEncoder(w).Encode(resp)

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
