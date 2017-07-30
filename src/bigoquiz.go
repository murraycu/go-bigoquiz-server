package bigoquiz

import (
	"net/http"
)

func init() {
	http.HandleFunc("/api/quiz", restHandleQuiz)
}
