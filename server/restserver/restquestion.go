package restserver

import (
	"github.com/julienschmidt/httprouter"
	"github.com/murraycu/go-bigoquiz-server/domain/quiz"
	"github.com/murraycu/go-bigoquiz-server/repositories/db"
	"net/http"
)

func (s *RestServer) HandleQuestionNext(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		handleErrorAsHttpError(w, http.StatusInternalServerError, "No quiz-id specified")
		return
	}

	q := s.getQuiz(quizId)
	if q == nil {
		handleErrorAsHttpError(w, http.StatusNotFound, "quiz not found")
		return
	}

	userId, err := s.getUserIdFromSessionAndDb(r, w)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "logged-in check failed. getUserIdFromSessionAndDb() failed: %v", err)
		return
	}

	var question *quiz.Question
	if len(userId) == 0 {
		//The user is not logged in,
		//so just return a random question:
		question = q.GetRandomQuestion(sectionId)
	} else {
		c := r.Context()

		dbClient, err := db.NewUserDataRepository()
		if err != nil {
			handleErrorAsHttpError(w, http.StatusInternalServerError, "failed getting stats for user (failed to connect to user data repository). NewUserDataRepository() failed: %v", err)
			return
		}

		if len(sectionId) == 0 {
			mapUserStats, err := dbClient.GetUserStatsForQuiz(c, userId, quizId)
			if err != nil {
				handleErrorAsHttpError(w, http.StatusInternalServerError, "failed getting stats for user. GetUserStatsForQuiz() failed: %v", err)
				return
			}

			question = s.getNextQuestionFromUserStats("", q, mapUserStats)
		} else {
			//This special case is a bit copy-and-pasty of the general case with the
			//map, but it seems more efficient to avoid an unnecessary Map.
			userStats, err := dbClient.GetUserStatsForSection(c, userId, sectionId, quizId)
			if err != nil {
				handleErrorAsHttpError(w, http.StatusInternalServerError, "failed getting stats for user for section. GetUserStatsForSection() failed: %v", err)
				return
			}

			question = s.getNextQuestionFromUserStatsForSection(sectionId, q, userStats)
		}
	}

	if question == nil {
		handleErrorAsHttpError(w, http.StatusNotFound, "question not found")
		return
	}

	question.SetQuestionExtras(q)

	marshalAndWriteOrHttpError(w, question)
}
