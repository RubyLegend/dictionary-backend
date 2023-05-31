package quizRoutes

import (
	"encoding/json"
	"net/http"

	"github.com/RubyLegend/dictionary-backend/middleware/cors"
	"github.com/julienschmidt/httprouter"
)

type WordResp struct {
	QuizId      int    `json:"quizId"`
	Name        string `json:"name"`
	Translation string `json:"translation"`
}

var words = []WordResp{
	{1, "111", "111"},
	{1, "qweqwe33333", "wewew"},
	{1, "wewqe", "dsdsa"},
	{1, "wewqeqw", "eqweqw"},
	{1, "123123123", "12312"},
}

var quizResult = make([]struct {
	Word     string `json:"word"`
	Expected string `json:"expected"`
	Answer   string `json:"answer"`
	Status   bool   `json:"status"`
}, len(words))

var i = 0

func QuizGetNewWord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	resp := make(map[string]any)

	if i >= len(words) {
		resp["results"] = quizResult
		i = 0
	} else {
		resp["quizId"] = words[i].QuizId
		if len(words[i].Name) != 0 {
			resp["name"] = words[i].Name
		} else {
			resp["translation"] = words[i].Translation
		}
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func QuizPostNewQuizWords(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	resp := make(map[string]any)
	resp["quizId"] = 1
	_ = json.NewEncoder(w).Encode(resp)
}

func QuizPostVerifyWord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	resp := make(map[string]any)

	var word WordResp
	_ = json.NewDecoder(r.Body).Decode(&word)

	if word.Translation == words[i].Translation {
		resp["status"] = "correct"
		quizResult[i].Word = words[i].Name
		quizResult[i].Expected = words[i].Translation
		quizResult[i].Answer = word.Translation
		quizResult[i].Status = true
		i = i + 1
	} else {
		resp["error"] = "answer incorrect"
		quizResult[i].Word = words[i].Name
		quizResult[i].Expected = words[i].Translation
		quizResult[i].Answer = word.Translation
		quizResult[i].Status = false
		i = i + 1
		w.WriteHeader(http.StatusNotModified)
	}

	_ = json.NewEncoder(w).Encode(resp)
}
