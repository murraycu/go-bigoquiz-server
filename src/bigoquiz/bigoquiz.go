package bigoquiz

import (
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"golang.org/x/oauth2"
	"net/http"
	"encoding/gob"
	"log"
	"config"
)

func init() {
	config, err := config.GenerateConfig()
	if err != nil {
		log.Println("Could not load config file: %v", err)
		return
	}

	// Create the session cookie store,
	// using the secret key from the configuration file.
	store = sessions.NewCookieStore([]byte(config.CookieKey))
	store.Options.HttpOnly = true

	// Gob encoding for gorilla/sessions
	// Otherwise, we will see errors such as this when calling store.Save():
	// "
	// Could not save session:'securecookie: error - caused by: securecookie: error - caused by: gob: type not registered for interface: oauth2.Token'
	// "
	gob.Register(&oauth2.Token{})

	router := httprouter.New()
	router.GET("/api/quiz", restHandleQuizAll)
	router.GET("/api/quiz/:quizId", restHandleQuizById)
	router.GET("/api/quiz/:quizId/section", restHandleQuizSectionsByQuizId)
	router.GET("/api/quiz/:quizId/question/:questionId", restHandleQuizQuestionById)
	router.GET("/api/question/next", restHandleQuestionNext)
	router.GET("/api/user", restHandleUser)
	router.GET("/api/user-history", restHandleUserHistoryAll)
	router.GET("/api/user-history/:quizId", restHandleUserHistoryByQuizId)

	router.GET("/login/login", handleGoogleLogin)
	router.GET("/login/callback", handleGoogleCallback)
	router.GET("/login/logout", handleGoogleLogout)

	// Allow Javascript requests from some domains other than the one serving this API.
	// The browser issue a CORS request before actually issuing the HTTP request.
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://beta.bigoquiz.com"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true, // Note: The client needs to specify this too, or cookies won't be sent.
	})

	handler := c.Handler(router)
	http.Handle("/", handler)
}
