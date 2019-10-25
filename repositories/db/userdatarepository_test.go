package db

import (
	"context"
	domainuser "github.com/murraycu/go-bigoquiz-server/domain/user"
	"github.com/murraycu/go-bigoquiz-server/server/loginserver/oauthparsers"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"testing"
	"time"
)

func TestNewRestServerInstantiate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test which requires more setup.")
	}

	userDataClient, err := NewUserDataRepository()
	assert.Nil(t, err)
	assert.NotNil(t, userDataClient)
}

func TestNewRestServerGetUserProfileByIdForNonExistantUser(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test which requires more setup.")
	}

	userDataClient, err := NewUserDataRepository()
	assert.Nil(t, err)
	assert.NotNil(t, userDataClient)

	c := context.Background()

	// This must be decodable with datastore.DecodeKey().
	userId := "EhYKC1VzZXJQcm9maWxlEICAgICw2IIK"

	userProfile, err := userDataClient.GetUserProfileById(c, userId)
	assert.Nil(t, err)
	assert.Nil(t, userProfile)
}

func TestNewRestServerStoreAndGetUserProfileById(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test which requires more setup.")
	}

	userDataClient, err := NewUserDataRepository()
	assert.Nil(t, err)
	assert.NotNil(t, userDataClient)

	c := context.Background()

	// This must be decodable with datastore.DecodeKey().
	userId := "EhYKC1VzZXJQcm9maWxlEICAgICw2IIL"

	userInfo := oauthparsers.GoogleUserInfo{
		Sub:           "some-google-user-id",
		Email:         "example@example.com",
		EmailVerified: true,
		Name:          "Example McExample",
	}
	strUserId := "" // Create a new user.
	token := oauth2.Token{
		AccessToken: "some-access-token",
	}
	userId, err = userDataClient.StoreGoogleLoginInUserProfile(c, userInfo, strUserId, &token)
	assert.Nil(t, err)
	assert.NotNil(t, userId)

	userProfile, err := userDataClient.GetUserProfileById(c, userId)
	assert.Nil(t, err)
	assert.NotNil(t, userProfile)

	assert.Equal(t, userInfo.Email, userProfile.Email)
	assert.Equal(t, userInfo.Name, userProfile.Name)
}

func storeUserStatsInStore(t *testing.T, c context.Context, userDataClient *UserDataRepository, userId string) *domainuser.Stats {
	stats := domainuser.Stats{
		QuizId:    "some-quiz-id",
		SectionId: "some-section-id",

		Answered: 100,
		Correct:  90,

		CountQuestionsAnsweredOnce: 10,
		CountQuestionsCorrectOnce:  5,

		QuestionHistories: []domainuser.QuestionHistory{
			{
				QuestionId: "some-question-id",

				AnsweredCorrectlyOnce: true,
				CountAnsweredWrong:    -1,
			},
		},
	}

	err := userDataClient.StoreUserStats(c, userId, &stats)
	assert.Nil(t, err)

	return &stats
}

// Returns the user ID of the created user.
func createUserInStore(t *testing.T, c context.Context, userDataClient *UserDataRepository) string {
	userInfo := oauthparsers.GoogleUserInfo{
		Sub:           "some-google-user-id",
		Email:         "example@example.com",
		EmailVerified: true,
		Name:          "Example McExample",
	}
	// Create a new user.
	token := oauth2.Token{
		AccessToken: "some-access-token",
	}

	userId, err := userDataClient.StoreGoogleLoginInUserProfile(c, userInfo, "", &token)
	assert.Nil(t, err)
	assert.NotNil(t, userId)

	return userId
}

func TestNewRestServerStoreAndGetStatsForSection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test which requires more setup.")
	}

	userDataClient, err := NewUserDataRepository()
	assert.Nil(t, err)
	assert.NotNil(t, userDataClient)

	c := context.Background()

	userId := createUserInStore(t, c, userDataClient)

	stats := storeUserStatsInStore(t, c, userDataClient, userId)

	result, err := userDataClient.GetUserStatsForSection(c, userId, stats.QuizId, stats.SectionId)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, stats.QuizId, result.QuizId)
	assert.Equal(t, stats.SectionId, result.SectionId)
	assert.Equal(t, stats.Answered, result.Answered)
	assert.Equal(t, stats.Correct, result.Correct)
	assert.Equal(t, stats.CountQuestionsAnsweredOnce, result.CountQuestionsAnsweredOnce)
	assert.Equal(t, stats.CountQuestionsCorrectOnce, result.CountQuestionsCorrectOnce)

	qa0 := stats.QuestionHistories[0]
	questionHistoryQuestionId := qa0.QuestionId
	assert.Equal(t, qa0.AnsweredCorrectlyOnce, result.GetQuestionWasAnswered(questionHistoryQuestionId))
	assert.Equal(t, qa0.CountAnsweredWrong, result.GetQuestionCountAnsweredWrong(questionHistoryQuestionId))
}

func TestNewRestServerStoreAndGetStatsForQuiz(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test which requires more setup.")
	}

	userDataClient, err := NewUserDataRepository()
	assert.Nil(t, err)
	assert.NotNil(t, userDataClient)

	c := context.Background()

	userId := createUserInStore(t, c, userDataClient)

	stats := storeUserStatsInStore(t, c, userDataClient, userId)

	result, err := userDataClient.GetUserStatsForQuiz(c, userId, stats.QuizId)
	assert.Nil(t, err)
	assert.NotEmpty(t, result)

	resultSection, ok := result[stats.SectionId]
	assert.True(t, ok)
	assert.NotNil(t, resultSection)

	assert.Equal(t, stats.SectionId, resultSection.SectionId)
}

func TestNewRestServerStoreAndGetStatsForAll(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test which requires more setup.")
	}

	userDataClient, err := NewUserDataRepository()
	assert.Nil(t, err)
	assert.NotNil(t, userDataClient)

	c := context.Background()

	userId := createUserInStore(t, c, userDataClient)

	stats := storeUserStatsInStore(t, c, userDataClient, userId)

	result, err := userDataClient.GetUserStats(c, userId)
	assert.Nil(t, err)
	assert.NotEmpty(t, result)

	resultQuiz, ok := result[stats.QuizId]
	assert.True(t, ok)
	assert.NotNil(t, resultQuiz)
	assert.Equal(t, stats.QuizId, resultQuiz.QuizId)
}

/* TODO:
func TestNewRestServerStoreAndDeleteStatsForSection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test which requires more setup.")
	}

	userDataClient, err := NewUserDataRepository()
	assert.Nil(t, err)
	assert.NotNil(t, userDataClient)

	c := context.Background()

	userId := createUserInStore(t, c, userDataClient)

	stats := storeUserStatsInStore(t, c, userDataClient, userId)

	err = userDataClient.DeleteUserStatsForQuiz(c, userId, stats.QuizId)
	assert.Nil(t, err)

	result, err := userDataClient.GetUserStatsForSection(c, userId, stats.QuizId, stats.SectionId)
	assert.Nil(t, err)
	assert.Nil(t, result)
}
*/


func TestNewRestServerUpdateStatsCorrectly(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test which requires more setup.")
	}

	userDataClient, err := NewUserDataRepository()
	assert.Nil(t, err)
	assert.NotNil(t, userDataClient)

	c := context.Background()

	quizId := "some-quiz-id-1"
	sectionId := "some-section-id"
	questionId := "some-question-id"

	userId := createUserInStore(t, c, userDataClient)

	sectionStats := new(domainuser.Stats)
	sectionStats.QuizId = quizId
	sectionStats.SectionId = sectionId
	sectionStats.UpdateStatsForAnswerCorrectness(questionId, true)

	err = userDataClient.StoreUserStats(c, userId, sectionStats)
	assert.Nil(t, err)

	// This seems to be necessary for the datastore emulator to let us read the data back reliably.
	time.Sleep(time.Millisecond * 500)

	result, err := userDataClient.GetUserStatsForSection(c, userId, quizId, sectionId)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, result.QuizId, quizId)
	assert.Equal(t, result.SectionId, sectionId)

	assert.Equal(t, 1, result.Answered)
	assert.Equal(t, 1, result.Correct)
	assert.Equal(t, 1, result.CountQuestionsAnsweredOnce)
	assert.Equal(t, 1, result.CountQuestionsCorrectOnce)
}