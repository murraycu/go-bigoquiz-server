package quizzes

import (
	"fmt"
	dtoquiz "github.com/murraycu/go-bigoquiz-server/repositories/quizzes/dtos/quiz"
	"github.com/stretchr/testify/assert"
	"testing"
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
		Id:   "some-id",
		Link: "some-link",
		Text: dtoquiz.Text{
			Text:   "some-text",
			IsHtml: true,
		},
	}

	result, err := convertDtoQuestionToDomainQuestion(&dto)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, dto.Id, result.Id)
	assert.Equal(t, dto.Link, result.Link)
	assert.Equal(t, dto.Text.Text, result.Text.Text)
	assert.Equal(t, dto.Text.IsHtml, result.Text.IsHtml)
}

func TestConvertDtoQAToDomainQA(t *testing.T) {
	dto := dtoquiz.QuestionAndAnswer{
		Question: dtoquiz.Question{
			Id:   "some-id",
			Link: "some-link",
			Text: dtoquiz.Text{
				Text:   "some-text",
				IsHtml: true,
			},
		},
		Answer: dtoquiz.Text{
			Text:   "some-text",
			IsHtml: true,
		},
	}

	result, err := convertDtoQAToDomainQA(&dto)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, dto.Question.Id, result.Question.Id)
	assert.Equal(t, dto.Question.Link, result.Question.Link)
	assert.Equal(t, dto.Question.Text.Text, result.Question.Text.Text)
	assert.Equal(t, dto.Question.Text.IsHtml, result.Question.Text.IsHtml)
	assert.Equal(t, dto.Answer.Text, result.Answer.Text)
	assert.Equal(t, dto.Answer.IsHtml, result.Answer.IsHtml)
}

func testQuestions(suffix string) []*dtoquiz.QuestionAndAnswer {
	return []*dtoquiz.QuestionAndAnswer{
		testQuestion(suffix + "0"),
		{
			Question: dtoquiz.Question{
				Id:   "some-question-id-1",
				Link: "some-question-link-1",
				Text: dtoquiz.Text{
					Text:   "some-text-1",
					IsHtml: true,
				},
			},
			Answer: dtoquiz.Text{
				Text:   "some-text-1",
				IsHtml: true,
			},
		},
		{
			Question: dtoquiz.Question{
				Id:   "some-id-2",
				Link: "some-link-2",
				Text: dtoquiz.Text{
					Text:   "some-text-2",
					IsHtml: false,
				},
			},
			Answer: dtoquiz.Text{
				Text:   "some-text-2",
				IsHtml: true,
			},
		},
	}
}

func testQuestion(suffix string) *dtoquiz.QuestionAndAnswer {
	return &dtoquiz.QuestionAndAnswer{
		Question: dtoquiz.Question{
			Id:   fmt.Sprintf("some-question-id-%v", suffix),
			Link: fmt.Sprintf("some-question-link-%v", suffix),
			Text: dtoquiz.Text{
				Text:   fmt.Sprintf("some-question-text-%v", suffix),
				IsHtml: true,
			},
		},
		Answer: dtoquiz.Text{
			Text:   fmt.Sprintf("some-answer-text-%v", suffix),
			IsHtml: true,
		},
	}
}

func testDefaultChoices(suffix string) []*dtoquiz.Text {
	return []*dtoquiz.Text{
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

func testSubSections(suffix string) []*dtoquiz.SubSection {
	return []*dtoquiz.SubSection{
		testSubSection(suffix + "0"),
		testSubSection(suffix + "1"),
	}
}

func testSubSection(suffix string) *dtoquiz.SubSection {
	return &dtoquiz.SubSection{
		HasIdAndTitle: dtoquiz.HasIdAndTitle{
			Id:    fmt.Sprintf("some-subsection-id-%v", suffix),
			Title: fmt.Sprintf("some-subsection-title-%v", suffix),
			Link:  fmt.Sprintf("some-subsection-link-%v", suffix),
		},
		Questions:        testQuestions(suffix),
		AnswersAsChoices: false,
	}
}

func testSection(suffix string) *dtoquiz.Section {
	return &dtoquiz.Section{
		HasIdAndTitle: dtoquiz.HasIdAndTitle{
			Id:    fmt.Sprintf("some-id-%v", suffix),
			Title: fmt.Sprintf("some-title-%v", suffix),
			Link:  fmt.Sprintf("some-link-%v", suffix),
		},
		Questions:        testQuestions(suffix),
		SubSections:      testSubSections(suffix),
		DefaultChoices:   testDefaultChoices(suffix),
		AnswersAsChoices: true,
		AndReverse:       true,
	}
}

func testSections(suffix string) []*dtoquiz.Section {
	return []*dtoquiz.Section{
		testSection(suffix + "0"),
		testSection(suffix + "1"),
		testSection(suffix + "2"),
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

	assert.Equal(t, dto.AnswersAsChoices, result.AnswersAsChoices)

	// Test extras:

	dtoSubSection := dto.SubSections[0]
	subSection := result.GetSubSection(dtoSubSection.Id)
	assert.NotNil(t, subSection)
	assert.Equal(t, dtoSubSection.Id, subSection.Id)
	assert.Equal(t, dtoSubSection.Title, subSection.Title)

	assert.Equal(t, 9, result.CountQuestions)
}

func testQuiz(suffix string) *dtoquiz.Quiz {
	return &dtoquiz.Quiz{
		HasIdAndTitle: dtoquiz.HasIdAndTitle{
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
	assert.Equal(t, dto.AnswersAsChoices, result.AnswersAsChoices)

	assert.Len(t, result.Sections, len(dto.Sections))
	assert.Equal(t, dto.Sections[0].Id, result.Sections[0].Id)

	assert.Len(t, result.Questions, len(dto.Questions))
	assert.Equal(t, dto.Questions[0].Id, result.Questions[0].Id)

	assert.Equal(t, dto.UsesMathML, result.UsesMathML)

	// Test extras:
	randomQuestion := result.GetRandomQuestion(dto.Sections[0].Id)
	assert.NotNil(t, randomQuestion)
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
