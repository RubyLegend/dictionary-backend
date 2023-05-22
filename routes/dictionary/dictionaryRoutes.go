package dictionaryRoutes

import (
	"encoding/json"
	// "log"
	"net/http"

	dictionaryRepo "github.com/RubyLegend/dictionary-backend/repository/dictionary"
	"github.com/julienschmidt/httprouter"
)

func DictionaryGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	var UserId int = 1

	resp := make(map[string]any)

	errors, dictionary := dictionaryRepo.GetDictionary(UserId)

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

func DictionaryPatch(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]any)

	var UserId int = 2
	DictionaryId := p.ByName("id")

	var dictionaryData dictionaryRepo.Dictionary
	_ = json.NewDecoder(r.Body).Decode(&dictionaryData)

	errors, UpdatedDictionary := dictionaryRepo.UpdateDictionary(UserId, DictionaryId, dictionaryData)

	if len(errors) > 0 {
		var errorMessages []string
		for _, v := range errors {
			errorMessages = append(errorMessages, v.Error())
		}
		resp["error"] = errorMessages
		w.WriteHeader(http.StatusBadRequest)
	} else {
		resp["dictionary"] = UpdatedDictionary
		w.WriteHeader(http.StatusOK)
	}
	_ = json.NewEncoder(w).Encode(resp)

}

func DictionaryDelete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]any)

	var UserId int = 2
	DictionaryId := p.ByName("id")

	errors := dictionaryRepo.DeleteDictionary(UserId, DictionaryId)

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
