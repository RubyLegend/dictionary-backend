package routes

import (
  "fmt"
  "net/http"
  "log"

  "github.com/RubyLegend/dictionary-backend/helpers/users"

  "github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Hello, World\n")
}

func WordGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Not Implemented\n")
}

func WordPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Not Implemented\n")
}

func WordDelete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Not Implemented\n")
}

func WordPatch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Not Implemented\n")
}

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

func UserLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  r.ParseForm()
  if r.Form["username"] == nil {
    fmt.Fprintf(w, "Username not provided. Cannot authorize.\n")
  } else {
    token, err := users.GenerateJWT(r.Form["username"][0])
    fmt.Fprintf(w, "Token: " + token + "\n")
    fmt.Fprintf(w, "Errors: %v\n", err)
  }
  fmt.Fprintf(w, "Under Construction\n")
}

func UserSignup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Not Implemented\n")
}

func UserLogout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Not Implemented\n")
}

func UserStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  if r.Header["Authorization"] == nil {
    fmt.Fprintf(w, "Authorization header not found.\n")
  } else {
    tokenString := r.Header["Authorization"][0]
    log.Println("Supplied token: " + tokenString)
    if ok := users.VerifyAuthorizationToken(tokenString); !ok {
      log.Println("Token doesn't start with 'Bearer '. Token incorrect.")
    } else { 
      tokenClear := tokenString[7:]
      claims, err := users.VerifyJWT(tokenClear)
      if err != nil {
        log.Println("Errors while parsing token.")
      } else {
        fmt.Fprintf(w, claims + "\n")
      }
    }
  }
  fmt.Fprintf(w, "Under Construction\n")
}

func UserRestoreUsername(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Not Implemented\n")
}

func UserRestorePassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Not Implemented\n")
}

func UserDelete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Not Implemented\n")
}

func UserPatch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Not Implemented\n")
}

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

func HistoryGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Not Implemented\n")
}

func HistoryDelete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Not Implemented\n")
}

func QuizGetNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Not Implemented\n")
}

func QuizPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Not Implemented\n")
}

func QuizGetStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Not Implemented\n")
}

