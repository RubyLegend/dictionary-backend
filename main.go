package main

import (
  "log"
  "net/http"

  "github.com/julienschmidt/httprouter"

  "github.com/RubyLegend/dictionary-backend/routes"
)

func main() {
  // Creating http router
  router := httprouter.New()

  router.GET("/", routes.Index)

  router.POST("/api/v1/word", routes.WordPost)
  router.DELETE("/api/v1/word/:id", routes.WordDelete)
  router.GET("/api/v1/word", routes.WordGet)
  router.PATCH("/api/v1/word/:id", routes.WordPatch)

  router.POST("/api/v1/translation", routes.TranslationPost)
  router.DELETE("/api/v1/translation/:id", routes.TranslationDelete)
  router.GET("/api/v1/translation", routes.TranslationGet)
  router.PATCH("/api/v1/translation/:id", routes.TranslationPatch)

  router.POST("/api/v1/user/login", routes.UserLogin)
  router.POST("/api/v1/user/signup", routes.UserSignup)
  router.POST("/api/v1/user/logout", routes.UserLogout)
  router.GET("/api/v1/user/status", routes.UserStatus)
  router.POST("/api/v1/user/restore-username", routes.UserRestoreUsername)
  router.POST("/api/v1/user/restore-password", routes.UserRestorePassword)
  router.DELETE("/api/v1/user", routes.UserDelete)
  router.PATCH("/api/v1/user", routes.UserPatch)
  
  router.GET("/api/v1/dictionary", routes.DictionaryGet)
  router.POST("/api/v1/dictionary", routes.DictionaryPost)
  router.PATCH("/api/v1/dictionary/:id", routes.DictionaryPatch)
  router.DELETE("/api/v1/dictionary/:id", routes.DictionaryDelete)

  router.GET("/api/v1/history", routes.HistoryGet)
  router.DELETE("/api/v1/history", routes.HistoryDelete)

  router.GET("/api/v1/quiz/new", routes.QuizGetNew)
  router.POST("/api/v1/quiz/:quizId", routes.QuizPost)
  router.GET("/api/v1/quiz/status", routes.QuizGetStatus)

  log.Println("Server started at http://localhost:8080")
  log.Fatal(http.ListenAndServe(":8080", router))
}
