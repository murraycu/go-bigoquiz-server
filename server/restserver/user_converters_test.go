package restserver

import (
	domainuser "github.com/murraycu/go-bigoquiz-server/domain/user"
	restquiz "github.com/murraycu/go-bigoquiz-server/server/restserver/quiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testRestQuestions(prefix string) []*restquiz.QuestionAndAnswer {
	subPrefix := prefix + "_some_questions"

	return []*restquiz.QuestionAndAnswer{
		testRestQuestion(subPrefix + "-0"),
		testRestQuestion(subPrefix + "-1"),
		testRestQuestion(subPrefix + "-2"),
		testRestQuestion(subPrefix + "-3"),
	}
}

func testRestQuestion(prefix string) *restquiz.QuestionAndAnswer {
	subPrefix := prefix + "_some-qa-"

	return &restquiz.QuestionAndAnswer{
		Question: restquiz.Question{
			Id:   subPrefix + "some-question-id",
			Link: subPrefix + "some-question-link",
			Text: restquiz.Text{
				Text:   subPrefix + "some-question-text",
				IsHtml: true,
			},
		},
		Answer: restquiz.Text{
			Text:   subPrefix + "some-question-answer",
			IsHtml: true,
		},
	}
}

func testRestDefaultChoices(prefix string) []*restquiz.Text {
	subPrefix := prefix + "_some-default-choice"

	return []*restquiz.Text{
		{
			Text:   subPrefix + "-0-some-default-choice-text",
			IsHtml: true,
		},
		{
			Text:   subPrefix + "-1-some-default-choice-text",
			IsHtml: false,
		},
		{
			Text:   subPrefix + "-2-some-default-choice-text",
			IsHtml: true,
		},
	}
}

func testRestSubSections(prefix string) []*restquiz.SubSection {
	subPrefix := prefix + "_some-sub-sections"

	return []*restquiz.SubSection{
		testRestSubSection(subPrefix + "-0"),
		testRestSubSection(subPrefix + "-1"),
	}
}

func testRestSubSection(prefix string) *restquiz.SubSection {
	subPrefix := prefix + "_some-sub-section"

	return &restquiz.SubSection{
		HasIdAndTitle: restquiz.HasIdAndTitle{
			Id:    subPrefix + "-some-sub-section-id",
			Title: subPrefix + "-some-sub-section-title",
			Link:  subPrefix + "-some-sub-section-link",
		},
		Questions:        testRestQuestions(subPrefix + "-some-sub-section-questions"),
		AnswersAsChoices: true,
	}
}

func testRestSection(prefix string) *restquiz.Section {
	subPrefix := prefix + "_some_section"
	return &restquiz.Section{
		HasIdAndTitle: restquiz.HasIdAndTitle{
			Id:    subPrefix + "-id",
			Title: subPrefix + "-title",
			Link:  subPrefix + "_-link",
		},
		Questions:        testRestQuestions(subPrefix + "-questions"),
		SubSections:      testRestSubSections(subPrefix + "-sub-sections"),
		DefaultChoices:   testRestDefaultChoices(subPrefix + "-default-choices"),
		AnswersAsChoices: true,
	}
}

func testRestSections(prefix string) []*restquiz.Section {
	subPrefix := prefix + "_some-section"
	return []*restquiz.Section{
		testRestSection(subPrefix + "-0"),
		testRestSection(subPrefix + "-1"),
		testRestSection(subPrefix + "-2"),
	}
}

func testRestQuiz() *restquiz.Quiz {
	prefix := "some-quiz_"
	return &restquiz.Quiz{
		HasIdAndTitle: restquiz.HasIdAndTitle{
			Id:    prefix + "some-quiz-id",
			Title: prefix + "some-quiz-title",
			Link:  prefix + "some-quiz-link",
		},
		IsPrivate: true,

		Sections: testRestSections(prefix + "some-quiz-sections"),

		UsesMathML: true,
	}
}

func TestConvertDomainQuestionHistoryToRestQuestionHistory(t *testing.T) {
	quiz := testRestQuiz()
	question := quiz.Sections[1].Questions[1].Question

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

func TestConvertDomainStatsToRestStatsPerSection(t *testing.T) {
	quiz := testRestQuiz()

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

	quizCache, err := NewQuizCache(quiz)
	assert.Nil(t, err)
	assert.NotNil(t, quizCache)

	result, err := convertDomainStatsToRestStats(&obj, quizCache)
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

	// Extras:
	assert.Equal(t, quizCache.GetSectionQuestionsCount(section.Id), result.CountQuestions)
	assert.Equal(t, quizCache.Quiz.Title, result.QuizTitle)
	assert.Equal(t, section.Title, result.SectionTitle)
}

func TestConvertDomainStatsToRestStatsPerQuiz(t *testing.T) {
	quiz := testRestQuiz()

	obj := domainuser.Stats{
		QuizId:                     quiz.Id,
		SectionId:                  "",
		Answered:                   11,
		Correct:                    9,
		CountQuestionsAnsweredOnce: 5,
		CountQuestionsCorrectOnce:  4,
		QuestionHistories:          nil,
	}

	quizCache, err := NewQuizCache(quiz)
	assert.Nil(t, err)
	assert.NotNil(t, quizCache)

	result, err := convertDomainStatsToRestStats(&obj, quizCache)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, obj.QuizId, result.QuizId)
	assert.Equal(t, obj.SectionId, result.SectionId)

	// Extras:
	assert.Equal(t, quizCache.GetQuestionsCount(), result.CountQuestions)
	assert.Equal(t, quizCache.Quiz.Title, result.QuizTitle)
	assert.Empty(t, result.SectionTitle)
}
