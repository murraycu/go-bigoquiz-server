package bigoquiz

import (
	"db"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"net/http"
	"quiz"
)

func restHandleQuestionNext(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var quizId string
	var sectionId string
	// var sectionId string
	queryValues := r.URL.Query()
	if queryValues != nil {
		quizId = queryValues.Get(QUERY_PARAM_QUIZ_ID)
		sectionId = queryValues.Get(QUERY_PARAM_SECTION_ID)
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

	userId, err := getUserIdFromSessionAndDb(r, w)
	if err != nil {
		http.Error(w, "logged-in check failed.", http.StatusInternalServerError)
		return
	}

	var question *quiz.Question
	if userId == nil {
		//The user is not logged in,
		//so just return a random question:
		question = q.GetRandomQuestion(sectionId)
	} else {
		c := appengine.NewContext(r)

		if len(sectionId) == 0 {
			mapUserStats, err := db.GetUserStatsForQuiz(c, userId, quizId)
			if err != nil {
				http.Error(w, "failed getting stats for user", http.StatusInternalServerError)
				return
			}

			question = getNextQuestionFromUserStats("", q, mapUserStats)
		} else {
			//This special case is a bit copy-and-pasty of the general case with the
			//map, but it seems more efficient to avoid an unnecessary Map.
			userStats, err := db.GetUserStatsForSection(c, userId, sectionId, quizId)
			if err != nil {
				http.Error(w, "failed getting stats for user for section", http.StatusInternalServerError)
				return
			}

			question = getNextQuestionFromUserStatsForSection(sectionId, q, userStats)
		}
	}

	if question == nil {
		http.Error(w, "question not found", http.StatusNotFound)
		return
	}

	question.SetQuestionExtras(q)

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
