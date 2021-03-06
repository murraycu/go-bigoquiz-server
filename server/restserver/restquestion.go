package restserver

import (
	"github.com/julienschmidt/httprouter"
	restquiz "github.com/murraycu/go-bigoquiz-server/server/restserver/quiz"
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

	quizCache, err := s.getQuizCache(q.Id)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusNotFound, "quiz cache not found")
		return
	}

	userId, err := s.getUserIdFromSessionAndDb(r)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "logged-in check failed. getUserIdFromSessionAndDb() failed: %v", err)
		return
	}

	var question *restquiz.Question
	if len(userId) == 0 {
		//The user is not logged in,
		//so just return a random question:
		question = quizCache.GetRandomQuestion(sectionId)
	} else {
		c := r.Context()

		if len(sectionId) == 0 {
			mapUserStats, err := s.userDataClient.GetUserStatsForQuiz(c, userId, quizId)
			if err != nil {
				handleErrorAsHttpError(w, http.StatusInternalServerError, "failed getting stats for user. GetUserStatsForQuiz() failed: %v", err)
				return
			}

			question, err = s.getNextQuestionFromUserStats("", q, mapUserStats)
			if err != nil {
				handleErrorAsHttpError(w, http.StatusInternalServerError, "getNextQuestionFromUserStats() failed")
				return
			}
		} else {
			//This special case is a bit copy-and-pasty of the general case with the
			//map, but it seems more efficient to avoid an unnecessary Map.
			userStats, err := s.userDataClient.GetUserStatsForSection(c, userId, sectionId, quizId)
			if err != nil {
				handleErrorAsHttpError(w, http.StatusInternalServerError, "failed getting stats for user for section. GetUserStatsForSection() failed: %v", err)
				return
			}

			question, err = s.getNextQuestionFromUserStatsForSection(sectionId, q, userStats)
			if err != nil {
				handleErrorAsHttpError(w, http.StatusInternalServerError, "getNextQuestionFromUserStatsForSection() failed")
				return
			}
		}
	}

	if question == nil {
		handleErrorAsHttpError(w, http.StatusNotFound, "question not found")
		return
	}

	marshalAndWriteOrHttpError(w, &question)
}
