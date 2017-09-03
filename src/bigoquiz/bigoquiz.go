package bigoquiz

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func init() {
	router := httprouter.New()
	router.GET("/api/quiz", restHandleQuizAll)
	router.GET("/api/quiz/:quizId", restHandleQuizById)
	router.GET("/api/quiz/:quizId/section", restHandleQuizSectionsByQuizId)
	router.GET("/api/quiz/:quizId/question/:questionId", restHandleQuizQuestionById)

	http.Handle("/", router)
}
