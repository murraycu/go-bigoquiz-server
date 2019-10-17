package db

import (
	domainuser "github.com/murraycu/go-bigoquiz-server/domain/user"
	dtouser "github.com/murraycu/go-bigoquiz-server/repositories/db/dtos/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertDtoQuestionHistoryToDomainQuestionHistory(t *testing.T) {
	dto := dtouser.QuestionHistory{
		QuestionId:            "question-id-1",
		AnsweredCorrectlyOnce: true,
		CountAnsweredWrong:    3,
	}

	result := convertDtoQuestionHistoryToDomainQuestionHistory(dto)
	assert.NotNil(t, result)

	assert.Equal(t, dto.QuestionId, result.QuestionId)
	assert.Equal(t, dto.AnsweredCorrectlyOnce, result.AnsweredCorrectlyOnce)
	assert.Equal(t, dto.CountAnsweredWrong, result.CountAnsweredWrong)
}

func TestConvertDtoStatsToDomainStats(t *testing.T) {
	dto := dtouser.Stats{
		QuizId:                     "example-quiz-id-1",
		SectionId:                  "example-section-id-2",
		Answered:                   11,
		Correct:                    9,
		CountQuestionsAnsweredOnce: 5,
		CountQuestionsCorrectOnce:  4,
		QuestionHistories: []dtouser.QuestionHistory{
			{
				QuestionId:            "question-id-1",
				AnsweredCorrectlyOnce: true,
				CountAnsweredWrong:    1,
			},
			{
				QuestionId:            "question-id-2",
				AnsweredCorrectlyOnce: false,
				CountAnsweredWrong:    2,
			},
		},
	}

	result := convertDtoStatsToDomainStats(&dto)
	assert.NotNil(t, result)

	assert.Equal(t, dto.QuizId, result.QuizId)
	assert.Equal(t, dto.SectionId, result.SectionId)
	assert.Equal(t, dto.Answered, result.Answered)
	assert.Equal(t, dto.Correct, result.Correct)
	assert.Equal(t, dto.CountQuestionsAnsweredOnce, result.CountQuestionsAnsweredOnce)
	assert.Equal(t, dto.CountQuestionsCorrectOnce, result.CountQuestionsCorrectOnce)

	assert.Len(t, dto.QuestionHistories, 2)

	qh0 := dto.QuestionHistories[0]
	assert.NotNil(t, qh0)
	assert.Equal(t, dto.QuestionHistories[0].QuestionId, qh0.QuestionId)
	assert.Equal(t, dto.QuestionHistories[0].AnsweredCorrectlyOnce, qh0.AnsweredCorrectlyOnce)
	assert.Equal(t, dto.QuestionHistories[0].CountAnsweredWrong, qh0.CountAnsweredWrong)

	qh1 := dto.QuestionHistories[1]
	assert.NotNil(t, qh0)
	assert.Equal(t, dto.QuestionHistories[1].QuestionId, qh1.QuestionId)
	assert.Equal(t, dto.QuestionHistories[1].AnsweredCorrectlyOnce, qh1.AnsweredCorrectlyOnce)
	assert.Equal(t, dto.QuestionHistories[1].CountAnsweredWrong, qh1.CountAnsweredWrong)
}

func TestConvertDomainQuestionHistoryToDtoQuestionHistory(t *testing.T) {
	obj := domainuser.QuestionHistory{
		QuestionId:            "question-id-1",
		AnsweredCorrectlyOnce: true,
		CountAnsweredWrong:    3,
		/* These are not in the DTO:
		QuestionTitle:
		SectionId: "section-id-2",
		SubSectionTitle: "sub-section-title-example",
		*/
	}

	result := convertDomainQuestionHistoryToDtoQuestionHistory(obj)
	assert.NotNil(t, result)

	assert.Equal(t, obj.QuestionId, result.QuestionId)
	assert.Equal(t, obj.AnsweredCorrectlyOnce, result.AnsweredCorrectlyOnce)
	assert.Equal(t, obj.CountAnsweredWrong, result.CountAnsweredWrong)
}

func TestConvertDomainStatsToDtoStats(t *testing.T) {
	obj := domainuser.Stats{
		QuizId:                     "example-quiz-id-1",
		SectionId:                  "example-section-id-2",
		Answered:                   11,
		Correct:                    9,
		CountQuestionsAnsweredOnce: 5,
		CountQuestionsCorrectOnce:  4,
		QuestionHistories: []domainuser.QuestionHistory{
			{
				QuestionId:            "question-id-1",
				AnsweredCorrectlyOnce: true,
				CountAnsweredWrong:    1,
				// QuestionTitle, SectionId, and CountAnsweredWrong will not be in the dto anyway.
			},
			{
				QuestionId:            "question-id-2",
				AnsweredCorrectlyOnce: false,
				CountAnsweredWrong:    2,
			},
		},

		/* These will not be in the dto anyway.
		CountQuestions: 5,
		QuizTitle:      "example-quiz-title",
		SectionTitle:   "example-section-title",
		*/
	}

	userId := "EgsKB0FydGljbGUQAQ"
	result, err := convertDomainStatsToDtoStats(&obj, userId)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// We cannot test this via comparison without being aware of the datastore key encoding:
	// assert.Equal(t, userId, result.UserId)
	// So instead we just check that there is some value.
	assert.NotNil(t, result.UserId)

	assert.Equal(t, obj.QuizId, result.QuizId)
	assert.Equal(t, obj.SectionId, result.SectionId)
	assert.Equal(t, obj.Answered, result.Answered)
	assert.Equal(t, obj.Correct, result.Correct)
	assert.Equal(t, obj.CountQuestionsAnsweredOnce, result.CountQuestionsAnsweredOnce)
	assert.Equal(t, obj.CountQuestionsCorrectOnce, result.CountQuestionsCorrectOnce)

	assert.Len(t, result.QuestionHistories, 2)

	qh0 := result.QuestionHistories[0]
	assert.NotNil(t, qh0)
	assert.Equal(t, obj.QuestionHistories[0].QuestionId, qh0.QuestionId)
	assert.Equal(t, obj.QuestionHistories[0].AnsweredCorrectlyOnce, qh0.AnsweredCorrectlyOnce)
	assert.Equal(t, obj.QuestionHistories[0].CountAnsweredWrong, qh0.CountAnsweredWrong)

	qh1 := result.QuestionHistories[1]
	assert.NotNil(t, qh0)
	assert.Equal(t, obj.QuestionHistories[1].QuestionId, qh1.QuestionId)
	assert.Equal(t, obj.QuestionHistories[1].AnsweredCorrectlyOnce, qh1.AnsweredCorrectlyOnce)
	assert.Equal(t, obj.QuestionHistories[1].CountAnsweredWrong, qh1.CountAnsweredWrong)
}

func TestConvertDtoProfileToDomainProfile(t *testing.T) {
	dto := dtouser.Profile{
		Name:  "example name",
		Email: "example@example.com",

		// GoogleId and GoogleAccessToken don't appear in the domain struct.
		GoogleProfileUrl: "example-google-profile-url",

		// GitHubId and GitHubAccessToken don't appear in the domain struct.
		GitHubProfileUrl: "example-github-profile-url",

		// FacebookID and FacebookbAccessToken don't appear in the domain struct.
		FacebookProfileUrl: "example-facebook-profile-url",
	}

	result := convertDtoProfileToDomainProfile(&dto)
	assert.NotNil(t, result)
	assert.Equal(t, dto.Name, result.Name)
	assert.Equal(t, dto.Email, result.Email)
	assert.Equal(t, dto.GoogleProfileUrl, result.GoogleProfileUrl)
	assert.Equal(t, dto.GitHubProfileUrl, result.GitHubProfileUrl)
	assert.Equal(t, dto.FacebookProfileUrl, result.FacebookProfileUrl)
}
