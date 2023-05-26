package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/RubyLegend/dictionary-backend/routes"
	dictionaryRoutes "github.com/RubyLegend/dictionary-backend/routes/dictionary"
	historyRoutes "github.com/RubyLegend/dictionary-backend/routes/history"
	quizRoutes "github.com/RubyLegend/dictionary-backend/routes/quiz"
	translationRoutes "github.com/RubyLegend/dictionary-backend/routes/translation"
	userRoutes "github.com/RubyLegend/dictionary-backend/routes/user"
	wordRoutes "github.com/RubyLegend/dictionary-backend/routes/word"

	db "github.com/RubyLegend/dictionary-backend/middleware/database"
	userHelper "github.com/RubyLegend/dictionary-backend/middleware/users"
)

func main() {
	// Connecting to database
	db.OpenConnection()

	// Creating http router
	router := httprouter.New()

	router.GET("/", routes.Index)

	router.POST("/api/v1/word", wordRoutes.WordPost)
	router.DELETE("/api/v1/word/:id", wordRoutes.WordDelete)
	router.GET("/api/v1/word", wordRoutes.WordGet)
	router.PATCH("/api/v1/word/:id", wordRoutes.WordPatch)

	router.POST("/api/v1/translation", translationRoutes.TranslationPost)
	router.DELETE("/api/v1/translation/:id", translationRoutes.TranslationDelete)
	router.GET("/api/v1/translation", translationRoutes.TranslationGet)
	router.PATCH("/api/v1/translation/:id", translationRoutes.TranslationPatch)

	router.POST("/api/v1/user/login", userRoutes.UserLogin)   // done
	router.POST("/api/v1/user/signup", userRoutes.UserSignup) // done
	router.POST("/api/v1/user/logout", userRoutes.UserLogout) // done
	router.GET("/api/v1/user/status", userRoutes.UserStatus)  // done
	router.POST("/api/v1/user/restore-username", userRoutes.UserRestoreUsername)
	router.POST("/api/v1/user/restore-password", userRoutes.UserRestorePassword)
	router.DELETE("/api/v1/user", userRoutes.UserDelete) // done
	router.PATCH("/api/v1/user", userRoutes.UserPatch)   // done

	router.GET("/api/v1/dictionary", dictionaryRoutes.DictionaryGet)
	router.POST("/api/v1/dictionary", dictionaryRoutes.DictionaryPost)
	router.PATCH("/api/v1/dictionary/:id", dictionaryRoutes.DictionaryPatch)
	router.DELETE("/api/v1/dictionary/:id", dictionaryRoutes.DictionaryDelete)

	router.GET("/api/v1/history", historyRoutes.HistoryGet)
	router.DELETE("/api/v1/history", historyRoutes.HistoryDelete)

	router.GET("/api/v1/quiz/new", quizRoutes.QuizGetNew)
	router.POST("/api/v1/quiz/:quizId", quizRoutes.QuizPost)
	router.GET("/api/v1/quiz/status", quizRoutes.QuizGetStatus)

	// Just for now logout monitor will detach to it's own thread here
	go userHelper.LogoutMonitor()

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
