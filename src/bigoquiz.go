package bigoquiz

import (
	"encoding/json"
	"net/http"
)

type quiz struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	IsPrivate bool   `json:"isPrivate"`
}

func init() {
	http.HandleFunc("/api/quiz", rest_quizzes)
}

func rest_quizzes(w http.ResponseWriter, r *http.Request) {
	var quizzes [2]quiz
	testQuiz := &quizzes[0]
	testQuiz.Id = "test-id"
	testQuiz.Title = "Test Title"
	testQuiz.IsPrivate = false

	testQuiz = &quizzes[1]
	testQuiz.Id = "test-id2"
	testQuiz.Title = "Test Title 2"
	testQuiz.IsPrivate = true

	w.Header().Set("Content-Type", "application/json") // normal header
	w.WriteHeader(http.StatusOK)

	jsonStr, err := json.Marshal(quizzes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonStr)
}
