package restserver

import (
	"fmt"
	"github.com/murraycu/go-bigoquiz-server/repositories"
	"github.com/murraycu/go-bigoquiz-server/repositories/db"
	"github.com/murraycu/go-bigoquiz-server/usersessionstore"
)

const QUERY_PARAM_QUIZ_ID = "quiz-id"
const QUERY_PARAM_SECTION_ID = "section-id"
const QUERY_PARAM_QUESTION_ID = "question-id"
const QUERY_PARAM_LIST_ONLY = "list-only"
const QUERY_PARAM_NEXT_QUESTION_SECTION_ID = "next-question-section-id"
const PATH_PARAM_QUIZ_ID = "quizId"
const PATH_PARAM_QUESTION_ID = "questionId"

type RestServer struct {
	Quizzes *repositories.QuizzesRepository

	UserDataClient *db.UserDataRepository

	// Session cookie store.
	UserSessionStore *usersessionstore.UserSessionStore
}

func NewRestServer(quizzes *repositories.QuizzesRepository, userSessionStore *usersessionstore.UserSessionStore) (*RestServer, error) {
	result := &RestServer{}

	var err error
	result.UserDataClient, err = db.NewUserDataRepository()
	if err != nil {
		return nil, fmt.Errorf("NewUserDataRepository() failed: %v", err)
	}

	result.Quizzes = quizzes
	result.UserSessionStore = userSessionStore

	return result, nil
}
