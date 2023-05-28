package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"

	"github.com/RubyLegend/dictionary-backend/routes"
	dictionaryRoutes "github.com/RubyLegend/dictionary-backend/routes/dictionary"
	historyRoutes "github.com/RubyLegend/dictionary-backend/routes/history"
	quizRoutes "github.com/RubyLegend/dictionary-backend/routes/quiz"
	translationRoutes "github.com/RubyLegend/dictionary-backend/routes/translation"
	userRoutes "github.com/RubyLegend/dictionary-backend/routes/user"
	wordRoutes "github.com/RubyLegend/dictionary-backend/routes/word"

	"github.com/RubyLegend/dictionary-backend/middleware/cors"
	db "github.com/RubyLegend/dictionary-backend/middleware/database"
	userHelper "github.com/RubyLegend/dictionary-backend/middleware/users"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connecting to database
	db.OpenConnection()

	// Creating http router
	router := httprouter.New()

	router.GET("/", routes.Index)
	router.OPTIONS("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		cors.Setup(w, r)
	}) // done

	router.GET("/api/v1/word", wordRoutes.WordGet)
	router.POST("/api/v1/word", wordRoutes.WordPost)
	router.OPTIONS("/api/v1/word", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		cors.Setup(w, r)
	}) // done
	router.DELETE("/api/v1/word/:id", wordRoutes.WordDelete)
	router.PATCH("/api/v1/word/:id", wordRoutes.WordPatch)
	router.OPTIONS("/api/v1/word/:id", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		cors.Setup(w, r)
	}) // done

	router.POST("/api/v1/translation", translationRoutes.TranslationPost)
	router.DELETE("/api/v1/translation/:id", translationRoutes.TranslationDelete)
	router.GET("/api/v1/translation", translationRoutes.TranslationGet)
	router.PATCH("/api/v1/translation/:id", translationRoutes.TranslationPatch)

	router.POST("/api/v1/user/login", userRoutes.UserLogin) // done
	router.OPTIONS("/api/v1/user/login", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		cors.Setup(w, r)
	}) // done
	router.POST("/api/v1/user/signup", userRoutes.UserSignup) // done
	router.OPTIONS("/api/v1/user/signup", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		cors.Setup(w, r)
	}) // done
	router.POST("/api/v1/user/logout", userRoutes.UserLogout) // done
	router.OPTIONS("/api/v1/user/logout", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		cors.Setup(w, r)
	}) // done
	router.GET("/api/v1/user/status", userRoutes.UserStatus) // done
	router.OPTIONS("/api/v1/user/status", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		cors.Setup(w, r)
	}) // done
	router.POST("/api/v1/user/restore-username", userRoutes.UserRestoreUsername)
	router.POST("/api/v1/user/restore-password", userRoutes.UserRestorePassword)
	router.DELETE("/api/v1/user", userRoutes.UserDelete) // done
	router.PATCH("/api/v1/user", userRoutes.UserPatch)   // done
	router.OPTIONS("/api/v1/user", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		cors.Setup(w, r)
	}) // done

	router.GET("/api/v1/dictionary", dictionaryRoutes.DictionaryGet)   // done
	router.POST("/api/v1/dictionary", dictionaryRoutes.DictionaryPost) // done
	router.OPTIONS("/api/v1/dictionary", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		cors.Setup(w, r)
	}) // done
	router.GET("/api/v1/dictionary/:id", dictionaryRoutes.DictionaryGetWords)  // done
	router.PATCH("/api/v1/dictionary/:id", dictionaryRoutes.DictionaryPatch)   // done
	router.DELETE("/api/v1/dictionary/:id", dictionaryRoutes.DictionaryDelete) // done
	router.OPTIONS("/api/v1/dictionary/:id", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		cors.Setup(w, r)
	}) // done

	router.GET("/api/v1/history", historyRoutes.HistoryGet)
	router.DELETE("/api/v1/history", historyRoutes.HistoryDelete)
	router.OPTIONS("/api/v1/history", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		cors.Setup(w, r)
	}) // done

	router.GET("/api/v1/quiz/new", quizRoutes.QuizGetNew)
	router.POST("/api/v1/quiz/:quizId", quizRoutes.QuizPost)
	router.GET("/api/v1/quiz/status", quizRoutes.QuizGetStatus)

	// Just for now logout monitor will detach to it's own thread here
	go userHelper.LogoutMonitor()

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
