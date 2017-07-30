package bigoquiz

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func init() {
	router := httprouter.New()
	router.GET("/api/quiz", restHandleQuiz)

	http.Handle("/", router)
}
