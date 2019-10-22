package restserver

import (
	"fmt"
	"github.com/murraycu/go-bigoquiz-server/repositories"
	"github.com/murraycu/go-bigoquiz-server/repositories/db"
	"github.com/murraycu/go-bigoquiz-server/server/usersessionstore"
	"log"
	"net/http"
)

const QUERY_PARAM_QUIZ_ID = "quiz-id"
const QUERY_PARAM_SECTION_ID = "section-id"
const QUERY_PARAM_QUESTION_ID = "question-id"
const QUERY_PARAM_LIST_ONLY = "list-only"
const QUERY_PARAM_NEXT_QUESTION_SECTION_ID = "next-question-section-id"
const PATH_PARAM_QUIZ_ID = "quizId"
const PATH_PARAM_QUESTION_ID = "questionId"

type RestServer struct {
	quizzes *repositories.QuizzesAndCaches

	userDataClient *db.UserDataRepository

	// Session cookie store.
	userSessionStore *usersessionstore.UserSessionStore
}

func NewRestServer(quizzesStore *repositories.QuizzesRepository, userSessionStore *usersessionstore.UserSessionStore) (*RestServer, error) {
	result := &RestServer{}

	var err error
	result.userDataClient, err = db.NewUserDataRepository()
	if err != nil {
		return nil, fmt.Errorf("NewUserDataRepository() failed: %v", err)
	}

	result.quizzes, err = quizzesStore.GetQuizzesAndCaches()
	result.userSessionStore = userSessionStore

	return result, nil
}

func handleErrorAsHttpError(w http.ResponseWriter, code int, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	log.Print(msg)

	http.Error(w, msg, code)
}
