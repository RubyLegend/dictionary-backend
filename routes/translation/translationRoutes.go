package translationRoutes

import (
	"fmt"
	"net/http"
  
	"github.com/julienschmidt/httprouter"
  )

  
	func TranslationGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "Not Implemented\n")
  }
  
  func TranslationPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "Not Implemented\n")
  }
  
  func TranslationDelete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "Not Implemented\n")
  }
  
  func TranslationPatch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "Not Implemented\n")
  }
  