package routes

import (
  "fmt"
  "net/http"

  "github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Hello, World\n")
}


