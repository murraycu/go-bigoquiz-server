package restserver

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/murraycu/go-bigoquiz-server/domain/quiz"
	"net/http"
	"strconv"
)

func (s *RestServer) HandleQuizAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	listOnly := false
	queryValues := r.URL.Query()
	if queryValues != nil {
		listOnlyStr := queryValues.Get(QUERY_PARAM_LIST_ONLY)
		listOnly, _ = strconv.ParseBool(listOnlyStr)
	}

	var quizArray []*quiz.Quiz = nil
	if listOnly {
		quizArray = s.quizzes.QuizzesListSimple
	} else {
		quizArray = s.quizzes.QuizzesListFull
	}

	w.Header().Set("Content-Type", "application/json") // normal header
	w.WriteHeader(http.StatusOK)

	jsonStr, err := json.Marshal(quizArray)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *RestServer) getQuiz(quizId string) *quiz.Quiz {
	return s.quizzes.Quizzes[quizId]
}

func (s *RestServer) HandleQuizById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	quizId := ps.ByName(PATH_PARAM_QUIZ_ID)
	if quizId == "" {
		// This makes no sense. restHandleQuizAll() should have been called.
		http.Error(w, "Empty quiz ID", http.StatusInternalServerError)
		return
	}

	q := s.getQuiz(quizId)
	if q == nil {
		http.Error(w, "quiz not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // normal header
	w.WriteHeader(http.StatusOK)

	jsonStr, err := json.Marshal(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *RestServer) HandleQuizSectionsByQuizId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	listOnly := false
	queryValues := r.URL.Query()
	if queryValues != nil {
		listOnlyStr := queryValues.Get(QUERY_PARAM_LIST_ONLY)
		listOnly, _ = strconv.ParseBool(listOnlyStr)
	}

	quizId := ps.ByName(PATH_PARAM_QUIZ_ID)
	if quizId == "" {
		// This makes no sense. restHandleQuizAll() should have been called.
		http.Error(w, "Empty quiz ID", http.StatusInternalServerError)
		return
	}

	q := s.getQuiz(quizId)
	if q == nil {
		http.Error(w, "quiz not found", http.StatusInternalServerError)
		return
	}

	sections := q.Sections
	if listOnly {
		simpleSections := make([]*quiz.Section, 0, len(sections))
		for _, s := range sections {
			var simple quiz.Section
			s.CopyHasIdAndTitle(&simple.HasIdAndTitle)
			simpleSections = append(simpleSections, &simple)
		}

		sections = simpleSections
	}

	jsonStr, err := json.Marshal(sections)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *RestServer) HandleQuizQuestionById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	quizId := ps.ByName(PATH_PARAM_QUIZ_ID)
	if quizId == "" {
		// This makes no sense. restHandleQuizAll() should have been called.
		http.Error(w, "Empty quiz ID", http.StatusInternalServerError)
		return
	}

	questionId := ps.ByName(PATH_PARAM_QUESTION_ID)
	if questionId == "" {
		// This makes no sense.
		http.Error(w, "Empty question ID", http.StatusInternalServerError)
		return
	}

	q := s.getQuiz(quizId)
	if q == nil {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	}

	qa := q.GetQuestionAndAnswer(questionId)
	if qa == nil {
		http.Error(w, "question not found", http.StatusInternalServerError)
		return
	}

	qa.Question.SetQuestionExtras(q)

	jsonStr, err := json.Marshal(qa.Question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
