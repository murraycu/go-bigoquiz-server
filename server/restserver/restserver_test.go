package restserver

import (
	"context"
	"log"
	"path/filepath"

	"github.com/gorilla/sessions"
	domainquiz "github.com/murraycu/go-bigoquiz-server/domain/quiz"
	domainuser "github.com/murraycu/go-bigoquiz-server/domain/user"
	"github.com/murraycu/go-bigoquiz-server/server/loginserver/oauthparsers"
	"golang.org/x/oauth2"

	"net/http"
	"testing"

	"github.com/murraycu/go-bigoquiz-server/repositories/db"
	"github.com/murraycu/go-bigoquiz-server/repositories/quizzes"
	"github.com/murraycu/go-bigoquiz-server/server/usersessionstore"
	"github.com/stretchr/testify/assert"
)

type MockUserSessionStore struct {
}

func (m MockUserSessionStore) GetSession(r *http.Request) (*sessions.Session, error) {
	panic("Unimplemented")
}

func (m MockUserSessionStore) GetUserIdAndOAuthTokenFromSession(r *http.Request) (*usersessionstore.UserIdAndOAuthToken, error) {
	panic("Unimplemented")
}

type MockUserDataRepository struct{}

func (m MockUserDataRepository) GetUserProfileById(c context.Context, strUserId string) (*domainuser.Profile, error) {
	panic("Unimplemented")
}

func (m MockUserDataRepository) GetUserStats(c context.Context, strUserId string) (map[string]*domainuser.Stats, error) {
	panic("Unimplemented")
}

func (m MockUserDataRepository) GetUserStatsForQuiz(c context.Context, strUserId string, quizId string) (map[string]*domainuser.Stats, error) {
	panic("Unimplemented")
}

func (m MockUserDataRepository) GetUserStatsForSection(c context.Context, strUserId string, quizId string, sectionId string) (*domainuser.Stats, error) {
	panic("Unimplemented")
}

func (m MockUserDataRepository) StoreUserStats(c context.Context, userID string, stats *domainuser.Stats) error {
	panic("Unimplemented")
}

func (m MockUserDataRepository) DeleteUserStatsForQuiz(c context.Context, strUserId string, quizId string) error {
	panic("Unimplemented")
}

func (m MockUserDataRepository) StoreGoogleLoginInUserProfile(c context.Context, userInfo oauthparsers.GoogleUserInfo, strUserId string, token *oauth2.Token) (string, error) {
	panic("Unimplemented")
}

func (m MockUserDataRepository) StoreGitHubLoginInUserProfile(c context.Context, userInfo oauthparsers.GitHubUserInfo, strUserId string, token *oauth2.Token) (string, error) {
	panic("Unimplemented")
}

func (m MockUserDataRepository) StoreFacebookLoginInUserProfile(c context.Context, userInfo oauthparsers.FacebookUserInfo, strUserId string, token *oauth2.Token) (string, error) {
	panic("Unimplemented")
}

type MockQuizzesRepository struct{}

func (m MockQuizzesRepository) LoadQuizzes() (quizzes.MapQuizzes, error) {
	return map[string]*domainquiz.Quiz{
		"id1": &domainquiz.Quiz{},
		"id2": &domainquiz.Quiz{},
	}, nil
}

func TestNewRestServer(t *testing.T) {
	userSessionStore := &MockUserSessionStore{}
	userDataRepository := &MockUserDataRepository{}
	quizzesStore := &MockQuizzesRepository{}

	restServer, err := NewRestServer(quizzesStore, userSessionStore, userDataRepository)
	assert.Nil(t, err)
	assert.NotNil(t, restServer)
}

func TestHasQuizzes(t *testing.T) {
	userSessionStore := &MockUserSessionStore{}
	userDataRepository := &MockUserDataRepository{}
	quizzesStore := &MockQuizzesRepository{}

	restServer, err := NewRestServer(quizzesStore, userSessionStore, userDataRepository)
	assert.Nil(t, err)

	assert.NotEmpty(t, restServer.quizzesListSimple)
	assert.NotEmpty(t, restServer.quizzesListFull)
}

func TestNewRestServerWithDataStore(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test which requires more setup.")
	}

	userSessionStore, err := usersessionstore.NewUserSessionStore("some-test-value")
	assert.Nil(t, err)
	assert.NotNil(t, userSessionStore)

	userDataClient, err := db.NewUserDataRepository()
	assert.Nil(t, err)
	assert.NotNil(t, userDataClient)

	// TODO: Mock the QuizzesRepository.
	directoryFilepath, err := filepath.Abs("../../quizzes")
	if err != nil {
		log.Fatalf("Couldn't get absolute filepath for quizzes: %v", err)
		return
	}

	quizzesStore, err := quizzes.NewQuizzesRepository(directoryFilepath)
	assert.Nil(t, err)
	assert.NotNil(t, quizzesStore)

	restServer, err := NewRestServer(quizzesStore, userSessionStore, userDataClient)
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
