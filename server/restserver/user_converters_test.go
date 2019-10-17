package restserver

import (
	"fmt"
	domainquiz "github.com/murraycu/go-bigoquiz-server/domain/quiz"
	domainuser "github.com/murraycu/go-bigoquiz-server/domain/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testQuestions(suffix string) []*domainquiz.QuestionAndAnswer {
	return []*domainquiz.QuestionAndAnswer{
		testQuestion(suffix + "0"),
		{
			Question: domainquiz.Question{
				Id:   "some-question-id-1",
				Link: "some-question-link-1",
				Text: domainquiz.Text{
					Text:   "some-text-1",
					IsHtml: true,
				},
			},
			Answer: domainquiz.Text{
				Text:   "some-text-1",
				IsHtml: true,
			},
		},
		{
			Question: domainquiz.Question{
				Id:   "some-id-2",
				Link: "some-link-2",
				Text: domainquiz.Text{
					Text:   "some-text-2",
					IsHtml: false,
				},
			},
			Answer: domainquiz.Text{
				Text:   "some-text-2",
				IsHtml: true,
			},
		},
	}
}

func testQuestion(suffix string) *domainquiz.QuestionAndAnswer {
	return &domainquiz.QuestionAndAnswer{
		Question: domainquiz.Question{
			Id:   fmt.Sprintf("some-question-id-%v", suffix),
			Link: fmt.Sprintf("some-question-link-%v", suffix),
			Text: domainquiz.Text{
				Text:   fmt.Sprintf("some-question-text-%v", suffix),
				IsHtml: true,
			},
		},
		Answer: domainquiz.Text{
			Text:   fmt.Sprintf("some-answer-text-%v", suffix),
			IsHtml: true,
		},
	}
}

func testDefaultChoices(suffix string) []*domainquiz.Text {
	return []*domainquiz.Text{
		{
			Text:   fmt.Sprintf("some-default-choice-1-%v", suffix),
			IsHtml: true,
		},
		{
			Text:   fmt.Sprintf("some-default-choice-2-%v", suffix),
			IsHtml: false,
		},
		{
			Text:   fmt.Sprintf("some-default-choice-3-%v", suffix),
			IsHtml: true,
		},
	}
}

func testSubSections(suffix string) []*domainquiz.SubSection {
	return []*domainquiz.SubSection{
		testSubSection(suffix + "0"),
		testSubSection(suffix + "1"),
	}
}

func testSubSection(suffix string) *domainquiz.SubSection {
	return &domainquiz.SubSection{
		HasIdAndTitle: domainquiz.HasIdAndTitle{
			Id:    fmt.Sprintf("some-subsection-id-%v", suffix),
			Title: fmt.Sprintf("some-subsection-title-%v", suffix),
			Link:  fmt.Sprintf("some-subsection-link-%v", suffix),
		},
		Questions:        testQuestions(suffix),
		AnswersAsChoices: false,
	}
}

func testSection(suffix string) *domainquiz.Section {
	return &domainquiz.Section{
		HasIdAndTitle: domainquiz.HasIdAndTitle{
			Id:    fmt.Sprintf("some-id-%v", suffix),
			Title: fmt.Sprintf("some-title-%v", suffix),
			Link:  fmt.Sprintf("some-link-%v", suffix),
		},
		Questions:        testQuestions(suffix),
		SubSections:      testSubSections(suffix),
		DefaultChoices:   testDefaultChoices(suffix),
		AnswersAsChoices: true,
	}
}

func testSections(suffix string) []*domainquiz.Section {
	return []*domainquiz.Section{
		testSection(suffix + "0"),
		testSection(suffix + "1"),
		testSection(suffix + "2"),
	}
}

func testQuiz(suffix string) *domainquiz.Quiz {
	return &domainquiz.Quiz{
		HasIdAndTitle: domainquiz.HasIdAndTitle{
			Id:    fmt.Sprintf("some-quiz-id-%v", suffix),
			Title: fmt.Sprintf("some-quiz-title-%v", suffix),
			Link:  fmt.Sprintf("some-quiz-link-%v", suffix),
		},
		IsPrivate:        true,
		AnswersAsChoices: true,

		Sections:  testSections(suffix),
		Questions: testQuestions(suffix),
	}
}

func TestConvertDomainQuestionHistoryToRestQuestionHistory(t *testing.T) {
	quiz := testQuiz("")
	question := quiz.Questions[1].Question

	obj := domainuser.QuestionHistory{
		QuestionId:            question.Id,
		AnsweredCorrectlyOnce: true,
		CountAnsweredWrong:    3,
	}

	result, err := convertDomainQuestionHistoryToRestQuestionHistory(obj, &question)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, obj.QuestionId, result.QuestionId)
	assert.Equal(t, obj.AnsweredCorrectlyOnce, result.AnsweredCorrectlyOnce)
	assert.Equal(t, obj.CountAnsweredWrong, result.CountAnsweredWrong)
}

func TestConvertDomainStatsToRestStats(t *testing.T) {
	quiz := testQuiz("")

	section := quiz.Sections[1]
	questionID0 := section.Questions[0].Id
	questionID1 := section.Questions[1].Id

	obj := domainuser.Stats{
		QuizId:                     quiz.Id,
		SectionId:                  section.Id,
		Answered:                   11,
		Correct:                    9,
		CountQuestionsAnsweredOnce: 5,
		CountQuestionsCorrectOnce:  4,
		QuestionHistories: []domainuser.QuestionHistory{
			{
				QuestionId:            questionID0,
				AnsweredCorrectlyOnce: true,
				CountAnsweredWrong:    1,
			},
			{
				QuestionId:            questionID1,
				AnsweredCorrectlyOnce: false,
				CountAnsweredWrong:    2,
			},
		},
	}

	result, err := convertDomainStatsToRestStats(&obj, quiz)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, obj.QuizId, result.QuizId)
	assert.Equal(t, obj.SectionId, result.SectionId)
	assert.Equal(t, obj.Answered, result.Answered)
	assert.Equal(t, obj.Correct, result.Correct)
	assert.Equal(t, obj.CountQuestionsAnsweredOnce, result.CountQuestionsAnsweredOnce)
	assert.Equal(t, obj.CountQuestionsCorrectOnce, result.CountQuestionsCorrectOnce)

	assert.Len(t, obj.QuestionHistories, 2)

	qh0 := obj.QuestionHistories[0]
	assert.NotNil(t, qh0)
	assert.Equal(t, obj.QuestionHistories[0].QuestionId, qh0.QuestionId)
	assert.Equal(t, obj.QuestionHistories[0].AnsweredCorrectlyOnce, qh0.AnsweredCorrectlyOnce)
	assert.Equal(t, obj.QuestionHistories[0].CountAnsweredWrong, qh0.CountAnsweredWrong)

	qh1 := obj.QuestionHistories[1]
	assert.NotNil(t, qh0)
	assert.Equal(t, obj.QuestionHistories[1].QuestionId, qh1.QuestionId)
	assert.Equal(t, obj.QuestionHistories[1].AnsweredCorrectlyOnce, qh1.AnsweredCorrectlyOnce)
	assert.Equal(t, obj.QuestionHistories[1].CountAnsweredWrong, qh1.CountAnsweredWrong)
}
