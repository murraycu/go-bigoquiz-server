package bigoquiz

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func restHandleQuestionNext(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var quizId string
	var sectionId string
	// var sectionId string
	queryValues := r.URL.Query()
	if queryValues != nil {
		quizId = queryValues.Get("quiz-id")
		sectionId = queryValues.Get("section-id")
	}

	if len(quizId) == 0 {
		// TODO: One day we might let the user answer questions from a
		// random quiz, so they wouldn't have to specify a quiz-id.
		http.Error(w, "No quiz-id specified", http.StatusBadRequest)
		return
	}

	q := getQuiz(quizId)
	if q == nil {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	}

	question := q.GetRandomQuestion(sectionId)
	if question == nil {
		http.Error(w, "question not found", http.StatusNotFound)
		return
	}

	jsonStr, err := json.Marshal(question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
