package main

import (
	"cloud.google.com/go/datastore"
	"encoding/gob"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/murraycu/go-bigoquiz-server/config"
	"github.com/murraycu/go-bigoquiz-server/repositories/quizzes"
	"github.com/murraycu/go-bigoquiz-server/server/loginserver"
	"github.com/murraycu/go-bigoquiz-server/server/restserver"
	"github.com/murraycu/go-bigoquiz-server/server/usersessionstore"
	"github.com/rs/cors"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
)

func main() {
	conf, err := config.GenerateConfig()
	if err != nil {
		log.Fatalf("Could not load conf file: %v\n", err)
		return
	}

	// Gob encoding for gorilla/sessions
	// Otherwise, we will see errors such as this when calling store.Save():
	// "
	// Could not save session:'securecookie: error - caused by: securecookie: error - caused by: gob: type not registered for interface: oauth2.Token'
	// "
	gob.Register(&oauth2.Token{})
	gob.Register(&datastore.Key{})

	userSessionStore, err := usersessionstore.NewUserSessionStore(conf.CookieKey)
	if err != nil {
		log.Fatalf("NewUserSessionStore failed: %v\n", err)
		return
	}

	quizzesStore, err := quizzes.NewQuizzesRepository()
	if err != nil {
		log.Fatalf("NewQuizzesRepository failed: %v\n", err)
		return
	}

	restServer, err := restserver.NewRestServer(quizzesStore, userSessionStore)
	if err != nil {
		log.Fatalf("NewRestServer failed: %v\n", err)
		return
	}

	loginServer, err := loginserver.NewLoginServer(userSessionStore)
	if err != nil {
		log.Fatalf("NewLoginServer failed: %v\n", err)
		return
	}

	router := httprouter.New()
	router.GET("/api/quiz", restServer.HandleQuizAll)
	router.GET("/api/quiz/:"+restserver.PATH_PARAM_QUIZ_ID, restServer.HandleQuizById)
	router.GET("/api/quiz/:"+restserver.PATH_PARAM_QUIZ_ID+"/section", restServer.HandleQuizSectionsByQuizId)
	router.GET("/api/quiz/:"+restserver.PATH_PARAM_QUIZ_ID+"/question/:"+restserver.PATH_PARAM_QUESTION_ID, restServer.HandleQuizQuestionById)

	router.GET("/api/question/next", restServer.HandleQuestionNext)

	router.GET("/api/user", restServer.HandleUser)

	router.GET("/api/user-history", restServer.HandleUserHistoryAll)
	router.GET("/api/user-history/:"+restserver.PATH_PARAM_QUIZ_ID, restServer.HandleUserHistoryByQuizId)
	router.POST("/api/user-history/submit-answer", restServer.HandleUserHistorySubmitAnswer)
	router.POST("/api/user-history/submit-dont-know-answer", restServer.HandleUserHistorySubmitDontKnowAnswer)
	router.POST("/api/user-history/reset-sections", restServer.HandleUserHistoryResetSections)

	router.GET("/login/login-google", loginServer.HandleGoogleLogin)
	router.GET("/login/"+config.PART_URL_LOGIN_CALLBACK_GOOGLE, loginServer.HandleGoogleCallback)
	router.GET("/login/login-github", loginServer.HandleGitHubLogin)
	router.GET("/login/"+config.PART_URL_LOGIN_CALLBACK_GITHUB, loginServer.HandleGitHubCallback)
	router.GET("/login/login-facebook", loginServer.HandleFacebookLogin)
	router.GET("/login/"+config.PART_URL_LOGIN_CALLBACK_FACEBOOK, loginServer.HandleFacebookCallback)
	router.GET("/login/logout", loginServer.HandleLogout)

	// Allow Javascript requests from some domains other than the one serving this API.
	// The browser issue a CORS request before actually issuing the HTTP request.
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{config.BaseUrl},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true, // Note: The client needs to specify this too, or cookies won't be sent.
	})

	handler := c.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handler))
}
