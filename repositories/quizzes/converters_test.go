package quizzes

import (
	"testing"

	dtoquiz "github.com/murraycu/go-bigoquiz-server/repositories/quizzes/dtos/quiz"
	"github.com/stretchr/testify/assert"
)

func TestConvertDtoHasIdAndTitleToDomainHasIdAndTitle(t *testing.T) {
	dto := dtoquiz.HasIdAndTitle{
		Id:    "some-id",
		Title: "some-title",
		Link:  "some-link",
	}

	result, err := convertDtoHasIdAndTitleToDomainHasIdAndTitle(&dto)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, dto.Id, result.Id)
	assert.Equal(t, dto.Title, result.Title)
	assert.Equal(t, dto.Link, result.Link)
}

func TestConvertDtoTextToDomainText(t *testing.T) {
	dto := dtoquiz.Text{
		Text:   "some-text",
		IsHtml: true,
	}

	result, err := convertDtoTextToDomainText(&dto)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, dto.Text, result.Text)
	assert.Equal(t, dto.IsHtml, result.IsHtml)
}

func TestConvertDtoQuestionToDomainQuestion(t *testing.T) {
	dto := dtoquiz.Question{
		Id:         "some-id",
		Link:       "some-link",
		TextSimple: "some-text",
	}

	result, err := convertDtoQuestionToDomainQuestion(&dto)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, dto.Id, result.Id)
	assert.Equal(t, dto.Link, result.Link)
	assert.Equal(t, dto.TextSimple, result.Text.Text)
	assert.Equal(t, false, result.Text.IsHtml)
}

func TestConvertDtoHtmlQuestionToDomainQuestion(t *testing.T) {
	dto := dtoquiz.Question{
		Id:   "some-id",
		Link: "some-link",
		TextDetail: dtoquiz.Text{
			Text:   "some-text",
			IsHtml: true,
		},
	}

	result, err := convertDtoQuestionToDomainQuestion(&dto)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, dto.Id, result.Id)
	assert.Equal(t, dto.Link, result.Link)
	assert.Equal(t, dto.TextDetail.Text, result.Text.Text)
	assert.Equal(t, dto.TextDetail.IsHtml, result.Text.IsHtml)
}

func TestConvertDtoQAToDomainQA(t *testing.T) {
	dto := dtoquiz.QuestionAndAnswer{
		Question: dtoquiz.Question{
			Id:   "some-id",
			Link: "some-link",
			TextDetail: dtoquiz.Text{
				Text:   "some-text",
				IsHtml: true,
			},
		},
		AnswerDetail: dtoquiz.Text{
			Text:   "some-text",
			IsHtml: true,
		},
	}

	result, err := convertDtoQAToDomainQA(&dto)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, dto.Question.Id, result.Question.Id)
	assert.Equal(t, dto.Question.Link, result.Question.Link)
	assert.Equal(t, dto.Question.TextDetail.Text, result.Question.Text.Text)
	assert.Equal(t, dto.Question.TextDetail.IsHtml, result.Question.Text.IsHtml)
	assert.Equal(t, dto.AnswerDetail.Text, result.Answer.Text)
	assert.Equal(t, dto.AnswerDetail.IsHtml, result.Answer.IsHtml)
}

func testQuestions(prefix string) []*dtoquiz.QuestionAndAnswer {
	subPrefix := prefix + "_some-questions_question-"
	return []*dtoquiz.QuestionAndAnswer{
		testQuestion(subPrefix + "0"),
		testQuestion(subPrefix + "1"),
		testQuestion(subPrefix + "2"),
		testQuestion(subPrefix + "3"),
	}
}

func testQuestion(prefix string) *dtoquiz.QuestionAndAnswer {
	subPrefix := prefix + "_qa_"
	return &dtoquiz.QuestionAndAnswer{
		Question: dtoquiz.Question{
			Id:   subPrefix + "some-question-id",
			Link: subPrefix + "some-question-link",
			TextDetail: dtoquiz.Text{
				Text:   subPrefix + "some-question-text",
				IsHtml: true,
			},
		},
		AnswerDetail: dtoquiz.Text{
			Text:   subPrefix + "some-answer-text",
			IsHtml: true,
		},
	}
}

func testDefaultChoices(prefix string) []*dtoquiz.Text {
	subPrefix := prefix + "_default_choices_some-default-choice-text"
	return []*dtoquiz.Text{
		{
			Text:   subPrefix + "-1",
			IsHtml: true,
		},
		{
			Text:   subPrefix + "-2",
			IsHtml: false,
		},
		{
			Text:   subPrefix + "-3",
			IsHtml: true,
		},
	}
}

func testSubSections(prefix string) []*dtoquiz.SubSection {
	subPrefix := prefix + "_some-sub-sections-"
	return []*dtoquiz.SubSection{
		testSubSection(subPrefix + "0"),
		testSubSection(subPrefix + "1"),
	}
}

func testSubSection(prefix string) *dtoquiz.SubSection {
	subPrefix := prefix + "_some-sub-section_"
	return &dtoquiz.SubSection{
		HasIdAndTitle: dtoquiz.HasIdAndTitle{
			Id:    subPrefix + "some-subsection-id",
			Title: subPrefix + "some-subsection-title",
			Link:  subPrefix + "some-subsection-link",
		},
		Questions:        testQuestions(prefix),
		AnswersAsChoices: false,
	}
}

func testSection(prefix string) *dtoquiz.Section {
	subPrefix := prefix + "_some-section_"

	return &dtoquiz.Section{
		HasIdAndTitle: dtoquiz.HasIdAndTitle{
			Id:    subPrefix + "some-section-id",
			Title: subPrefix + "some-section-title",
			Link:  subPrefix + "some-section-link",
		},
		Questions:        testQuestions(subPrefix + "some-questions"),
		SubSections:      testSubSections(prefix + "some-sub-sections"),
		DefaultChoices:   testDefaultChoices(prefix + "some-default-choices"),
		AnswersAsChoices: true,
		AndReverse:       true,
	}
}

func testSections(prefix string) []*dtoquiz.Section {
	subPrefix := prefix + "_some-section"
	return []*dtoquiz.Section{
		testSection(subPrefix + "-0"),
		testSection(subPrefix + "-1"),
		testSection(subPrefix + "-2"),
	}
}

func TestConvertDtoQuestionsToDomainQuestions(t *testing.T) {
	dto := testQuestions("foo")

	result, err := convertDtoQuestionsToDomainQuestions(dto)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Len(t, result, len(dto))

	// Don't test every field - leave that to the QuestionAndAnswer conversion test.
	assert.Equal(t, result[0].Question.Id, dto[0].Question.Id)
	assert.Equal(t, result[1].Question.Id, dto[1].Question.Id)
}

func TestConvertDtoSubSectionToDomainSubSection(t *testing.T) {
	dto := dtoquiz.SubSection{
		HasIdAndTitle: dtoquiz.HasIdAndTitle{
			Id:    "some-id",
			Title: "some-title",
			Link:  "some-link",
		},
		Questions:        testQuestions("foo"),
		AnswersAsChoices: true,
	}

	result, err := convertDtoSubSectionToDomainSubSection(&dto)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	// Don't test every field - leave that to the individual conversion tests.
	assert.Equal(t, dto.HasIdAndTitle.Id, result.HasIdAndTitle.Id)
	assert.Equal(t, dto.HasIdAndTitle.Title, result.HasIdAndTitle.Title)
	assert.Equal(t, dto.HasIdAndTitle.Link, result.HasIdAndTitle.Link)

	assert.Len(t, result.Questions, len(dto.Questions))
	assert.Equal(t, dto.Questions[0].Id, result.Questions[0].Id)

	assert.Equal(t, dto.AnswersAsChoices, result.AnswersAsChoices)
}

func TestConvertDtoSectionToDomainSection(t *testing.T) {
	dto := testSection("foo")

	result, err := convertDtoSectionToDomainSection(dto)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	// Don't test every field - leave that to the individual conversion tests.
	assert.Equal(t, dto.HasIdAndTitle.Id, result.HasIdAndTitle.Id)
	assert.Equal(t, dto.HasIdAndTitle.Title, result.HasIdAndTitle.Title)
	assert.Equal(t, dto.HasIdAndTitle.Link, result.HasIdAndTitle.Link)

	assert.Len(t, result.Questions, len(dto.Questions))
	assert.Equal(t, dto.Questions[0].Id, result.Questions[0].Id)

	assert.Len(t, result.SubSections, len(dto.SubSections))
	assert.Equal(t, dto.SubSections[0].Id, result.SubSections[0].Id)

	assert.Len(t, result.DefaultChoices, len(dto.DefaultChoices))
	assert.Equal(t, dto.DefaultChoices[0].Text, result.DefaultChoices[0].Text)
	assert.Equal(t, dto.DefaultChoices[0].IsHtml, result.DefaultChoices[0].IsHtml)
	assert.Equal(t, dto.DefaultChoices[1].Text, result.DefaultChoices[1].Text)
	assert.Equal(t, dto.DefaultChoices[1].IsHtml, result.DefaultChoices[1].IsHtml)
}

func testQuiz(prefix string) *dtoquiz.Quiz {
	subPrefix := prefix + "_some-quiz-"
	return &dtoquiz.Quiz{
		HasIdAndTitle: dtoquiz.HasIdAndTitle{
			Id:    subPrefix + "some-quiz-id",
			Title: subPrefix + "some-quiz-title",
			Link:  subPrefix + "some-quiz-link",
		},
		IsPrivate:        true,
		AnswersAsChoices: true,

		Sections:  testSections(subPrefix + "-some-quiz-sections"),
		Questions: testQuestions(subPrefix + "-some-quiz-questions"),
	}
}

func TestConvertDtoQuizToDomainQuiz(t *testing.T) {
	dto := testQuiz("foo")

	result, err := convertDtoQuizToDomainQuiz(dto)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	// Don't test every field - leave that to the individual conversion tests.
	assert.Equal(t, dto.HasIdAndTitle.Id, result.HasIdAndTitle.Id)
	assert.Equal(t, dto.HasIdAndTitle.Title, result.HasIdAndTitle.Title)
	assert.Equal(t, dto.HasIdAndTitle.Link, result.HasIdAndTitle.Link)

	assert.Equal(t, dto.IsPrivate, result.IsPrivate)

	assert.Len(t, result.Sections, len(dto.Sections))
	assert.Equal(t, dto.Sections[0].Id, result.Sections[0].Id)

	assert.Len(t, result.Questions, len(dto.Questions))
	assert.Equal(t, dto.Questions[0].Id, result.Questions[0].Id)

	assert.Equal(t, dto.UsesMathML, result.UsesMathML)
	assert.Equal(t, dto.AnswersAsChoices, result.AnswersAsChoices)
}

func TestConvertDtoQuizzesToDomainQuizzes(t *testing.T) {
	q0 := testQuiz("0")
	q1 := testQuiz("1")
	q2 := testQuiz("2")

	dto := map[string]*dtoquiz.Quiz{
		q0.Id: q0,
		q1.Id: q1,
		q2.Id: q2,
	}

	result, err := convertDtoQuizzesToDomainQuizzes(dto)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Len(t, result, len(dto))

	// Don't test every field - leave that to the individual conversion tests.
	q0Id := q0.HasIdAndTitle.Id
	assert.Equal(t, dto[q0Id].HasIdAndTitle.Id, result[q0Id].HasIdAndTitle.Id)
	assert.Equal(t, dto[q0Id].HasIdAndTitle.Title, result[q0Id].HasIdAndTitle.Title)
}
