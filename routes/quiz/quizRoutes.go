package quizRoutes

import (
	"fmt"
	"net/http"
  
	"github.com/julienschmidt/httprouter"
  )

  func QuizGetNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Not Implemented\n")
  }
  
  func QuizPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Not Implemented\n")
  }
  
  func QuizGetStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Not Implemented\n")
  }