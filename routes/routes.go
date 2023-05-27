package routes

import (
	"fmt"
	"net/http"

	"github.com/RubyLegend/dictionary-backend/middleware/cors"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cors.Setup(w, r)
	fmt.Fprintf(w, "Hello, World\n")
}
