package wordRoutes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RubyLegend/dictionary-backend/middleware/cors"
	wordRepo "github.com/RubyLegend/dictionary-backend/repository/words"
	"github.com/julienschmidt/httprouter"
)

func WordGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cors.Setup(w, r)
	fmt.Fprintf(w, "Not Implemented\n")
}

func WordPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)

	var wordData wordRepo.RequestType
	_ = json.NewDecoder(r.Body).Decode(&wordData)

	// resp := make(map[string]any)

	wordRepo.AddWord(wordData)

	fmt.Fprintf(w, "Not Implemented\n")
}

func WordDelete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cors.Setup(w, r)
	fmt.Fprintf(w, "Not Implemented\n")
}

func WordPatch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cors.Setup(w, r)
	fmt.Fprintf(w, "Not Implemented\n")
}
