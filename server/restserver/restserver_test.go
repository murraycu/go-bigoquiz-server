package restserver

import (
	"log"
	"path/filepath"
	"testing"

	"github.com/murraycu/go-bigoquiz-server/repositories/quizzes"
	"github.com/murraycu/go-bigoquiz-server/server/usersessionstore"
	"github.com/stretchr/testify/assert"
)

func TestNewRestServer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test which requires more setup.")
	}

	// TODO: Mock the UserSessionStore.
	userSessionStore, err := usersessionstore.NewUserSessionStore("some-test-value")
	assert.Nil(t, err)
	assert.NotNil(t, userSessionStore)

	// TODO: Mock the QuizzesRepository.
	directoryFilepath, err := filepath.Abs("../../quizzes")
	if err != nil {
		log.Fatalf("Couldn't get absolute filepath for quizzes: %v", err)
		return
	}

	quizzesStore, err := quizzes.NewQuizzesRepository(directoryFilepath)
	assert.Nil(t, err)
	assert.NotNil(t, quizzesStore)

	restServer, err := NewRestServer(quizzesStore, userSessionStore)
	assert.Nil(t, err)
	assert.NotNil(t, restServer)

	// TODO: Don't use private API.
	assert.NotEmpty(t, restServer.quizzesListSimple)
	assert.NotEmpty(t, restServer.quizzesListFull)
}

func testRestQuizzes() restQuizMap {
	quiz1 := testRestQuiz()
	quiz2 := testRestQuiz()

	return restQuizMap{
		quiz1.Id: quiz1,
		quiz2.Id: quiz2,
	}
}

func TestBuildQuizzesSimple(t *testing.T) {
	quizzes := testRestQuizzes()
	result := buildQuizzesSimple(quizzes)
	assert.NotEmpty(t, result)
}

func TestBuildQuizzesFull(t *testing.T) {
	quizzes := testRestQuizzes()
	result := buildQuizzesFull(quizzes)
	assert.NotEmpty(t, result)
}

func TestBuildQuizzesSimpleWithRealQuizzes(t *testing.T) {
	quizzes := loadRealRestQuizzes(t)
	result := buildQuizzesSimple(quizzes)
	assert.NotEmpty(t, result)
}

func TestBuildQuizzesFullWithRealQuizzes(t *testing.T) {
	quizzes := loadRealRestQuizzes(t)
	result := buildQuizzesFull(quizzes)
	assert.NotEmpty(t, result)
}
