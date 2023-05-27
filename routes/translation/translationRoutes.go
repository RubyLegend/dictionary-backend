package translationRoutes

import (
	"encoding/json"
	"net/http"

	translationRepo "github.com/RubyLegend/dictionary-backend/repository/translations"
	"github.com/julienschmidt/httprouter"
)

func TranslationGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	var WordId int = 1

	resp := make(map[string]any)

	errors, translation := translationRepo.GetTranslation(WordId)

	if len(errors) > 0 {
		var errorMessages []string
		for _, v := range errors {
			errorMessages = append(errorMessages, v.Error())
		}
		resp["error"] = errorMessages
		w.WriteHeader(http.StatusBadRequest)
	} else {
		resp["translation"] = translation
		w.WriteHeader(http.StatusOK)
	}
	_ = json.NewEncoder(w).Encode(resp)

}
func TranslationPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var WordId = 1

	var translationData translationRepo.Translation
	_ = json.NewDecoder(r.Body).Decode(&translationData)

	resp := make(map[string]any)
	translationData.WordId = WordId

	err := translationRepo.AddTranslation(translationData)
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

func TranslationDelete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]any)

	var WordId int = 2
	TranslationId := p.ByName("id")

	errors := translationRepo.DeleteTranslation(WordId, TranslationId)

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

func TranslationPatch(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]any)

	var WordId int = 2
	TranslationId := p.ByName("id")

	var translationData translationRepo.Translation
	_ = json.NewDecoder(r.Body).Decode(&translationData)

	errors, UpdatedTranslation := translationRepo.UpdateTranslation(WordId, TranslationId, translationData)

	if len(errors) > 0 {
		var errorMessages []string
		for _, v := range errors {
			errorMessages = append(errorMessages, v.Error())
		}
		resp["error"] = errorMessages
		w.WriteHeader(http.StatusBadRequest)
	} else {
		resp["translation"] = UpdatedTranslation
		w.WriteHeader(http.StatusOK)
	}
	_ = json.NewEncoder(w).Encode(resp)

}
