package bigoquiz

import (
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"net/http"
)

func init() {
	store = sessions.NewCookieStore([]byte("secret")) // TODO: Make secret configurable.
	store.Options.HttpOnly = true

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

	// Allow Javascript requests from some domains other than the one serving this API.
	// The browser issue a CORS request before actually issuing the HTTP request.
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://beta.bigoquiz.com"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	})
	handler := c.Handler(router)
	http.Handle("/", handler)
}
