package bigoquiz

import (
	"config"
	"encoding/gob"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"golang.org/x/oauth2"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
	"quiz"
)

const PATH_PARAM_QUIZ_ID = "quizId"
const PATH_PARAM_QUESTION_ID = "questionId"

const QUERY_PARAM_QUIZ_ID = "quiz-id"
const QUERY_PARAM_SECTION_ID = "section-id"
const QUERY_PARAM_QUESTION_ID = "question-id"
const QUERY_PARAM_LIST_ONLY = "list-only"
const QUERY_PARAM_NEXT_QUESTION_SECTION_ID = "next-question-section-id"

func init() {
	conf, err := config.GenerateConfig()
	if err != nil {
		log.Printf("Could not load conf file: %v\n", err)
		return
	}

	// Create the session cookie store,
	// using the secret key from the configuration file.
	store = sessions.NewCookieStore([]byte(conf.CookieKey))
	store.Options.HttpOnly = true
	store.Options.Secure = true // Only send via HTTPS connections, not HTTP.

	// Gob encoding for gorilla/sessions
	// Otherwise, we will see errors such as this when calling store.Save():
	// "
	// Could not save session:'securecookie: error - caused by: securecookie: error - caused by: gob: type not registered for interface: oauth2.Token'
	// "
	gob.Register(&oauth2.Token{})
	gob.Register(&datastore.Key{})

	loadQuizzesAndCaches()

	router := httprouter.New()
	router.GET("/api/quiz", restHandleQuizAll)
	router.GET("/api/quiz/:" + PATH_PARAM_QUIZ_ID, restHandleQuizById)
	router.GET("/api/quiz/:" + PATH_PARAM_QUIZ_ID + "/section", restHandleQuizSectionsByQuizId)
	router.GET("/api/quiz/:" + PATH_PARAM_QUIZ_ID + "/question/:" + PATH_PARAM_QUESTION_ID, restHandleQuizQuestionById)

	router.GET("/api/question/next", restHandleQuestionNext)

	router.GET("/api/user", restHandleUser)

	router.GET("/api/user-history", restHandleUserHistoryAll)
	router.GET("/api/user-history/:" + PATH_PARAM_QUIZ_ID, restHandleUserHistoryByQuizId)
	router.POST("/api/user-history/submit-answer", restHandleUserHistorySubmitAnswer)
	router.POST("/api/user-history/submit-dont-know-answer", restHandleUserHistorySubmitDontKnowAnswer)
	router.POST("/api/user-history/reset-sections", restHandleUserHistoryResetSections)

	router.GET("/login/login", handleGoogleLogin)
	router.GET("/login/callback", handleGoogleCallback)
	router.GET("/login/logout", handleGoogleLogout)

	// Allow Javascript requests from some domains other than the one serving this API.
	// The browser issue a CORS request before actually issuing the HTTP request.
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{config.BaseUrl},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true, // Note: The client needs to specify this too, or cookies won't be sent.
	})

	handler := c.Handler(router)
	http.Handle("/", handler)
}

var quizzes map[string]*quiz.Quiz
var quizzesListSimple []*quiz.Quiz
var quizzesListFull []*quiz.Quiz

func loadQuizzesAndCaches() {
	var err error
	quizzes, err = loadQuizzes()
	if err != nil {
		log.Printf("Could not load quiz files: %v\n", err)
		return
	}

	quizzesListSimple = buildQuizzesSimple(quizzes)
	quizzesListFull = buildQuizzesFull(quizzes)
}
