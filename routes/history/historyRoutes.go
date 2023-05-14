package historyRoutes

import (
	"fmt"
	"net/http"
  
	"github.com/julienschmidt/httprouter"
  )

  
  func HistoryGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Not Implemented\n")
  }
  
  func HistoryDelete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Not Implemented\n")
  }