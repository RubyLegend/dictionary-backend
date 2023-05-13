package routes

import (
  "fmt"
  "net/http"
  "log"
  "encoding/json"

  userHelper "github.com/RubyLegend/dictionary-backend/helpers/users"
  userRepo "github.com/RubyLegend/dictionary-backend/repository/users"

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
  var userData userRepo.User
  _ = json.NewDecoder(r.Body).Decode(&userData)

  if userData.Username == "" {
    fmt.Fprintf(w, "Username not provided. Cannot authorize.\n")
  } else {
    err := userHelper.VerifyCredentials(userData)
    if err != nil {
      log.Println(err)
      fmt.Fprintf(w, "%v\n", err)
    } else {
      token, err := userHelper.GenerateJWT(userData.Username)
      fmt.Fprintf(w, "Token: " + token + "\n")
      fmt.Fprintf(w, "Errors: %v\n", err)
    }
  }
  fmt.Fprintf(w, "Under Construction\n")
}

func UserSignup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  var userData userRepo.User
  _ = json.NewDecoder(r.Body).Decode(&userData)
  log.Println("Decoded Data:", userData)
  
  err := userRepo.AddUser(userData)

  if err != nil {
    log.Println(err)
    fmt.Fprintf(w, "%s", err)
  } else {
    fmt.Fprintf(w, "User added successfully\n")
  }

  fmt.Fprintf(w, "Under Construction\n")
}

func UserLogout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprintf(w, "Not Implemented\n")
}

func UserStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  if r.Header["Authorization"] == nil {
    fmt.Fprintf(w, "Authorization header not found.\n")
  } else {
    tokenString := r.Header["Authorization"][0]
    if ok := userHelper.VerifyAuthorizationToken(tokenString); !ok {
      log.Println("Token doesn't start with 'Bearer '. Token incorrect.")
    } else { 
      tokenClear := tokenString[7:]
      claims, err := userHelper.VerifyJWT(tokenClear)
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

