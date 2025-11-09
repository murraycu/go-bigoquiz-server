package restserver

import (
	"fmt"
	"testing"

	domainquiz "github.com/murraycu/go-bigoquiz-server/domain/quiz"
	"github.com/stretchr/testify/assert"
)

func TestConvertDomainHasIdAndTitleToRestHasIdAndTitle(t *testing.T) {
	obj := domainquiz.HasIdAndTitle{
		Id:    "some-id",
		Title: "some-title",
		Link:  "some-link",
	}

	result, err := convertDomainHasIdAndTitleToRestHasIdAndTitle(&obj)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, obj.Id, result.Id)
	assert.Equal(t, obj.Title, result.Title)
	assert.Equal(t, obj.Link, result.Link)
}

func TestConvertDomainTextToRestText(t *testing.T) {
	obj := domainquiz.Text{
		Text:   "some-text",
		IsHtml: true,
	}

	result, err := convertDomainTextToRestText(&obj)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, obj.Text, result.Text)
	assert.Equal(t, obj.IsHtml, result.IsHtml)
}

func TestConvertDomainQuestionToRestQuestion(t *testing.T) {
	obj := domainquiz.Question{
		Id:   "some-id",
		Link: "some-link",
		Text: domainquiz.Text{
			Text:   "some-text",
			IsHtml: true,
		},
	}

	result, err := convertDomainQuestionToRestQuestion(&obj)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, obj.Id, result.Id)
	assert.Equal(t, obj.Link, result.Link)
	assert.Equal(t, obj.Text.Text, result.Text.Text)
	assert.Equal(t, obj.Text.IsHtml, result.Text.IsHtml)
}

func TestConvertDomainQAToRestQA(t *testing.T) {
	obj := domainquiz.QuestionAndAnswer{
		Question: domainquiz.Question{
			Id:   "some-id",
			Link: "some-link",
			Text: domainquiz.Text{
				Text:   "some-text",
				IsHtml: true,
			},
		},
		Answer: domainquiz.Text{
			Text:   "some-text",
			IsHtml: true,
		},
	}

	result, err := convertDomainQAToRestQA(&obj)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, obj.Question.Id, result.Question.Id)
	assert.Equal(t, obj.Question.Link, result.Question.Link)
	assert.Equal(t, obj.Question.Text.Text, result.Question.Text.Text)
	assert.Equal(t, obj.Question.Text.IsHtml, result.Question.Text.IsHtml)
	assert.Equal(t, obj.Answer.Text, result.Answer.Text)
	assert.Equal(t, obj.Answer.IsHtml, result.Answer.IsHtml)

	// We check extras, such as QuizTitle, separately,
	// because that need the QuizCache
	// (which can't be instantiated until we have the REST Quiz.)
}

func testQuizCacheWithExtras(t *testing.T) *QuizCache {
	quizCache := testQuizCache(t)
	err := fillRestQuizExtrasFromQuizCache(quizCache.Quiz, quizCache)
	assert.Nil(t, err)
	return quizCache
}

func TestConvertDomainQuizToRestQuizWithExtrasSectionQuestion(t *testing.T) {
	quizCache := testQuizCacheWithExtras(t)

	questionId := quizCache.Quiz.Sections[1].Questions[1].Id

	qa := quizCache.GetQuestionAndAnswer(questionId)
	assert.NotNil(t, qa)

	assert.NotEmpty(t, qa.QuizTitle)
	assert.NotNil(t, qa.SectionId)
	assert.NotNil(t, qa.Section)
	assert.Empty(t, qa.SubSectionId)
	assert.Nil(t, qa.SubSection)
	assert.NotEmpty(t, qa.Choices)
}

func TestConvertDomainQuizToRestQuizWithExtrasSubSectionQuestion(t *testing.T) {
	quizCache := testQuizCacheWithExtras(t)

	questionId := quizCache.Quiz.Sections[1].SubSections[1].Questions[1].Id

	qa := quizCache.GetQuestionAndAnswer(questionId)
	assert.NotNil(t, qa)

	assert.NotEmpty(t, qa.QuizTitle)
	assert.NotEmpty(t, qa.SectionId)
	assert.NotNil(t, qa.Section)
	assert.NotEmpty(t, qa.SubSectionId)
	assert.NotNil(t, qa.SubSection)
	assert.NotEmpty(t, qa.Choices)
}

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

func TestConvertDomainQuestionsToRestQuestions(t *testing.T) {
	obj := testQuestions("foo")

	result, err := convertDomainQuestionsToRestQuestions(obj)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Len(t, result, len(obj))

	// Don't test every field - leave that to the QuestionAndAnswer conversion test.
	assert.Equal(t, result[0].Question.Id, obj[0].Question.Id)
	assert.Equal(t, result[1].Question.Id, obj[1].Question.Id)
}

func TestConvertDomainSubSectionToRestSubSection(t *testing.T) {
	obj := domainquiz.SubSection{
		HasIdAndTitle: domainquiz.HasIdAndTitle{
			Id:    "some-id",
			Title: "some-title",
			Link:  "some-link",
		},
		Questions: testQuestions("foo"),
	}

	result, err := convertDomainSubSectionToRestSubSection(&obj)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	// Don't test every field - leave that to the individual conversion tests.
	assert.Equal(t, obj.HasIdAndTitle.Id, result.HasIdAndTitle.Id)
	assert.Equal(t, obj.HasIdAndTitle.Title, result.HasIdAndTitle.Title)
	assert.Equal(t, obj.HasIdAndTitle.Link, result.HasIdAndTitle.Link)

	assert.Len(t, result.Questions, len(obj.Questions))
	assert.Equal(t, obj.Questions[0].Id, result.Questions[0].Id)

	assert.Equal(t, obj.AnswersAsChoices, result.AnswersAsChoices)
}

func TestConvertDomainSectionToRestSection(t *testing.T) {
	obj := testSection("foo")

	result, err := convertDomainSectionToRestSection(obj)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	// Don't test every field - leave that to the individual conversion tests.
	assert.Equal(t, obj.HasIdAndTitle.Id, result.HasIdAndTitle.Id)
	assert.Equal(t, obj.HasIdAndTitle.Title, result.HasIdAndTitle.Title)
	assert.Equal(t, obj.HasIdAndTitle.Link, result.HasIdAndTitle.Link)

	assert.Len(t, result.Questions, len(obj.Questions))
	assert.Equal(t, obj.Questions[0].Id, result.Questions[0].Id)

	assert.Len(t, result.SubSections, len(obj.SubSections))
	assert.Equal(t, obj.SubSections[0].Id, result.SubSections[0].Id)

	assert.Len(t, result.DefaultChoices, len(obj.DefaultChoices))
	assert.Equal(t, obj.DefaultChoices[0].Text, result.DefaultChoices[0].Text)
	assert.Equal(t, obj.DefaultChoices[0].IsHtml, result.DefaultChoices[0].IsHtml)
	assert.Equal(t, obj.DefaultChoices[1].Text, result.DefaultChoices[1].Text)
	assert.Equal(t, obj.DefaultChoices[1].IsHtml, result.DefaultChoices[1].IsHtml)

	assert.Equal(t, obj.AnswersAsChoices, result.AnswersAsChoices)
}

func testQuiz(suffix string) *domainquiz.Quiz {
	return &domainquiz.Quiz{
		HasIdAndTitle: domainquiz.HasIdAndTitle{
			Id:    fmt.Sprintf("some-quiz-id-%v", suffix),
			Title: fmt.Sprintf("some-quiz-title-%v", suffix),
			Link:  fmt.Sprintf("some-quiz-link-%v", suffix),
		},
		IsPrivate: true,

		Sections:  testSections(suffix),
		Questions: testQuestions(suffix),

		UsesMathML: false,

		AnswersAsChoices: true,
	}
}

func TestConvertDomainQuizToRestQuiz(t *testing.T) {
	obj := testQuiz("foo")

	result, err := convertDomainQuizToRestQuiz(obj)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	// Don't test every field - leave that to the individual conversion tests.
	assert.Equal(t, obj.HasIdAndTitle.Id, result.HasIdAndTitle.Id)
	assert.Equal(t, obj.HasIdAndTitle.Title, result.HasIdAndTitle.Title)
	assert.Equal(t, obj.HasIdAndTitle.Link, result.HasIdAndTitle.Link)

	assert.Equal(t, obj.IsPrivate, result.IsPrivate)

	assert.Len(t, result.Sections, len(obj.Sections))
	assert.Equal(t, obj.Sections[0].Id, result.Sections[0].Id)

	assert.Equal(t, obj.UsesMathML, result.UsesMathML)
}

func TestConvertDomainQuizzesToRestQuizzes(t *testing.T) {
	q0 := testQuiz("0")
	q1 := testQuiz("1")
	q2 := testQuiz("2")

	obj := map[string]*domainquiz.Quiz{
		q0.Id: q0,
		q1.Id: q1,
		q2.Id: q2,
	}

	result, err := convertDomainQuizzesToRestQuizzes(obj)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Len(t, result, len(obj))

	// Don't test every field - leave that to the individual conversion tests.
	q0Id := q0.HasIdAndTitle.Id
	assert.Equal(t, obj[q0Id].HasIdAndTitle.Id, result[q0Id].HasIdAndTitle.Id)
	assert.Equal(t, obj[q0Id].HasIdAndTitle.Title, result[q0Id].HasIdAndTitle.Title)
}
