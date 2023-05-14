package dictionaryRoutes

import (
	"fmt"
	"net/http"
  
	"github.com/julienschmidt/httprouter"
  )

  func DictionaryGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Not Implemented\n")
  }
  
  func DictionaryPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Not Implemented\n")
  }
  
  func DictionaryPatch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Not Implemented\n")
  }
  
  func DictionaryDelete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Not Implemented\n")
  }