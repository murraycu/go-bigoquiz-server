package quizzes

import (
	"fmt"
	domainquiz "github.com/murraycu/go-bigoquiz-server/domain/quiz"
	dtoquiz "github.com/murraycu/go-bigoquiz-server/repositories/quizzes/dtos/quiz"
)

func convertDtoQuizzesToDomainQuizzes(dtos map[string]*dtoquiz.Quiz) (MapQuizzes, error) {
	result := make(MapQuizzes)

	for _, dto := range dtos {
		quiz, err := convertDtoQuizToDomainQuiz(dto)
		if err != nil {
			return nil, fmt.Errorf("convertDtoQuizToDomainQuiz() failed: %v", err)
		}

		result[quiz.Id] = quiz
	}

	return result, nil
}

func convertDtoHasIdAndTitleToDomainHasIdAndTitle(dto *dtoquiz.HasIdAndTitle) (*domainquiz.HasIdAndTitle, error) {
	var result domainquiz.HasIdAndTitle

	result.Id = dto.Id
	result.Title = dto.Title
	result.Link = dto.Link

	return &result, nil
}

func convertDtoQuizToDomainQuiz(dto *dtoquiz.Quiz) (*domainquiz.Quiz, error) {
	var result domainquiz.Quiz

	hasIdAndTitle, err := convertDtoHasIdAndTitleToDomainHasIdAndTitle(&dto.HasIdAndTitle)
	if err != nil {
		return nil, fmt.Errorf("convertDtoHasIdAndTitleToDomainHasIdAndTitle() failed: %v", err)
	}

	result.HasIdAndTitle = *hasIdAndTitle

	result.IsPrivate = dto.IsPrivate

	for _, dtoSection := range dto.Sections {
		section, err := convertDtoSectionToDomainSection(dtoSection)
		if err != nil {
			return nil, fmt.Errorf("convertDtoSectionToDomainSection() failed: %v", err)
		}

		result.Sections = append(result.Sections, section)
	}

	result.Questions, err = convertDtoQuestionsToDomainQuestions(dto.Questions)
	if err != nil {
		return nil, fmt.Errorf("convertDtoQuestionsToDomainQuestions() failed: %v", err)
	}

	result.UsesMathML = dto.UsesMathML
	result.AnswersAsChoices = dto.AnswersAsChoices

	return &result, nil
}

func convertDtoQuestionsToDomainQuestions(dtos []*dtoquiz.QuestionAndAnswer) ([]*domainquiz.QuestionAndAnswer, error) {
	var result []*domainquiz.QuestionAndAnswer

	for _, dtoQA := range dtos {
		qa, err := convertDtoQAToDomainQA(dtoQA)
		if err != nil {
			return nil, fmt.Errorf("convertDtoQAToDomainQA() failed: %v", err)
		}

		result = append(result, qa)
	}

	return result, nil
}

func convertDtoTextToDomainText(dto *dtoquiz.Text) (*domainquiz.Text, error) {
	var result domainquiz.Text
	result.Text = dto.Text
	result.IsHtml = dto.IsHtml

	return &result, nil
}

func convertDtoQuestionToDomainQuestion(dto *dtoquiz.Question) (*domainquiz.Question, error) {
	var result domainquiz.Question
	result.Id = dto.Id
	result.Link = dto.Link

	text, err := convertDtoTextToDomainText(&dto.TextDetail)
	if err != nil {
		return nil, fmt.Errorf("convertDtoTextToDomainText() failed: %v", err)
	}

	result.Text = *text

	return &result, nil
}

func convertDtoQAToDomainQA(dto *dtoquiz.QuestionAndAnswer) (*domainquiz.QuestionAndAnswer, error) {
	var result domainquiz.QuestionAndAnswer

	question, err := convertDtoQuestionToDomainQuestion(&dto.Question)
	if err != nil {
		return nil, fmt.Errorf("convertDtoSectionToDomainSection() failed: %v", err)
	}

	result.Question = *question

	answer, err := convertDtoTextToDomainText(&dto.Answer)
	if err != nil {
		return nil, fmt.Errorf("convertDtoTextToDomainText() failed for answer: %v", err)
	}

	result.Answer = *answer

	return &result, nil
}

func convertDtoSubSectionToDomainSubSection(dto *dtoquiz.SubSection) (*domainquiz.SubSection, error) {
	var result domainquiz.SubSection

	hasIdAndTitle, err := convertDtoHasIdAndTitleToDomainHasIdAndTitle(&dto.HasIdAndTitle)
	if err != nil {
		return nil, fmt.Errorf("convertDtoHasIdAndTitleToDomainHasIdAndTitle() failed: %v", err)
	}

	result.HasIdAndTitle = *hasIdAndTitle

	result.Questions, err = convertDtoQuestionsToDomainQuestions(dto.Questions)
	if err != nil {
		return nil, fmt.Errorf("convertDtoQuestionsToDomainQuestions() failed for answer: %v", err)
	}

	result.AnswersAsChoices = dto.AnswersAsChoices

	return &result, nil
}

func convertDtoSectionToDomainSection(dto *dtoquiz.Section) (*domainquiz.Section, error) {
	var result domainquiz.Section

	hasIdAndTitle, err := convertDtoHasIdAndTitleToDomainHasIdAndTitle(&dto.HasIdAndTitle)
	if err != nil {
		return nil, fmt.Errorf("convertDtoHasIdAndTitleToDomainHasIdAndTitle() failed: %v", err)
	}

	result.HasIdAndTitle = *hasIdAndTitle

	result.Questions, err = convertDtoQuestionsToDomainQuestions(dto.Questions)
	if err != nil {
		return nil, fmt.Errorf("convertDtoQuestionsToDomainQuestions() failed: %v", err)
	}

	for _, dtoSubSection := range dto.SubSections {
		subSection, err := convertDtoSubSectionToDomainSubSection(dtoSubSection)
		if err != nil {
			return nil, fmt.Errorf("convertDtoSubSectionToDomainSubSection() failed: %v", err)
		}

		result.SubSections = append(result.SubSections, subSection)
	}

	for _, dtoText := range dto.DefaultChoices {
		defaultChoice, err := convertDtoTextToDomainText(dtoText)
		if err != nil {
			return nil, fmt.Errorf("convertDtoTextToDomainText() failed for DefaultChoices: %v", err)
		}

		result.DefaultChoices = append(result.DefaultChoices, defaultChoice)
	}

	result.AnswersAsChoices = dto.AnswersAsChoices

	return &result, nil
}
