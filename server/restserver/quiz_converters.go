package restserver

import (
	"fmt"
	domainquiz "github.com/murraycu/go-bigoquiz-server/domain/quiz"
	restquiz "github.com/murraycu/go-bigoquiz-server/server/restserver/quiz"
	"math/rand"
	"time"
)

type restQuizMap map[string]*restquiz.Quiz

// Then use fillRestQuizExtrasFromQuizCache() on each returned quiz.
func convertDomainQuizzesToRestQuizzes(objs map[string]*domainquiz.Quiz) (restQuizMap, error) {
	result := make(restQuizMap)

	for _, obj := range objs {
		quiz, err := convertDomainQuizToRestQuiz(obj)
		if err != nil {
			return nil, fmt.Errorf("convertDomainQuizToRestQuiz() failed: %v", err)
		}

		result[quiz.Id] = quiz
	}

	return result, nil
}

func convertDomainHasIdAndTitleToRestHasIdAndTitle(obj *domainquiz.HasIdAndTitle) (*restquiz.HasIdAndTitle, error) {
	var result restquiz.HasIdAndTitle

	result.Id = obj.Id
	result.Title = obj.Title
	result.Link = obj.Link

	return &result, nil
}

/** fillRestQuizExtrasFromQuizCache() fills fields that are specific to the REST datastructure,
 * such as the Question.Choices, based on other question answers.
 */
func fillRestQuizExtrasFromQuizCache(quiz *restquiz.Quiz, quizCache *QuizCache) error {
	for _, s := range quiz.Sections {
		for _, qa := range s.Questions {
			err := setQuestionExtras(&qa.Question, s, nil, quizCache)
			if err != nil {
				return fmt.Errorf("setQuestionExtras() failed: %v", err)
			}
		}

		for _, sub := range s.SubSections {
			for _, qa := range sub.Questions {
				err := setQuestionExtras(&qa.Question, s, sub, quizCache)
				if err != nil {
					return fmt.Errorf("setQuestionExtras() failed: %v", err)
				}
			}

			if sub.AnswersAsChoices && !s.AnswersAsChoices {
				setQuestionsChoicesFromAnswers(sub.Questions)
			}
		}

		//Make sure that we set sub-section choices from the answers from all questions in the whole section:
		// TODO: Test this.
		if s.AnswersAsChoices {
			questionsIncludingSubSections := make([]*restquiz.QuestionAndAnswer, 0)
			questionsIncludingSubSections = append(questionsIncludingSubSections, s.Questions...)

			for _, sub := range s.SubSections {
				questionsIncludingSubSections = append(questionsIncludingSubSections, sub.Questions...)
			}

			setQuestionsChoicesFromAnswers(questionsIncludingSubSections)
		}
	}

	return nil
}

const maxChoicesFromAnswers = 6

/** Get the index of an item in the array by comparig only the strings in the TextDetail struct.
 */
func getIndexInArray(array []*restquiz.Text, str *restquiz.Text) (int, bool) {
	for i, s := range array {
		if s.Text == str.Text {
			return i, true
		}
	}

	return 1, false
}

func shuffle(array []*restquiz.Text) {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	for i := len(array) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
}

/**
 * Create a smallenough set of choices which
 * always contains the correct answer.
 * This is slow.
 */
func reduceChoices(choices []*restquiz.Text, answer *restquiz.Text) []*restquiz.Text {
	result := make([]*restquiz.Text, len(choices))
	copy(result, choices)
	shuffle(result)

	answerIndex, ok := getIndexInArray(result, answer)
	if !ok {
		return nil
	}

	if answerIndex >= maxChoicesFromAnswers {
		result = result[0 : maxChoicesFromAnswers-1]
		result = append(result, answer)
		shuffle(result)
	} else {
		result = result[0:maxChoicesFromAnswers]
	}

	return result
}

func setQuestionsChoicesFromAnswers(questions []*restquiz.QuestionAndAnswer) {
	// Build the list of answers, avoiding duplicates:
	choices := make([]*restquiz.Text, 0, len(questions))
	used := make(map[string]bool)

	for _, q := range questions {
		t := q.Answer.Text
		_, ok := used[t]
		if ok {
			continue
		}

		used[t] = true
		choices = append(choices, &(q.Answer))
	}

	tooManyChoices := len(choices) > maxChoicesFromAnswers

	for _, q := range questions {
		if !tooManyChoices {
			q.Choices = choices
		} else {
			reduced := reduceChoices(choices, &(q.Answer))
			q.Choices = reduced
		}
	}
}

// Then use fillRestQuizExtrasFromQuizCache()
func convertDomainQuizToRestQuiz(obj *domainquiz.Quiz) (*restquiz.Quiz, error) {
	var result restquiz.Quiz

	hasIdAndTitle, err := convertDomainHasIdAndTitleToRestHasIdAndTitle(&obj.HasIdAndTitle)
	if err != nil {
		return nil, fmt.Errorf("convertDomainHasIdAndTitleToRestHasIdAndTitle() failed: %v", err)
	}

	result.HasIdAndTitle = *hasIdAndTitle

	result.IsPrivate = obj.IsPrivate

	for _, dtoSection := range obj.Sections {
		section, err := convertDomainSectionToRestSection(dtoSection)
		if err != nil {
			return nil, fmt.Errorf("convertDomainSectionToRestSection() failed: %v", err)
		}

		result.Sections = append(result.Sections, section)
	}

	result.UsesMathML = obj.UsesMathML

	return &result, nil
}

func convertDomainQuestionsToRestQuestions(objs []*domainquiz.QuestionAndAnswer) ([]*restquiz.QuestionAndAnswer, error) {
	var result []*restquiz.QuestionAndAnswer

	for _, dtoQA := range objs {
		qa, err := convertDomainQAToRestQA(dtoQA)
		if err != nil {
			return nil, fmt.Errorf("convertDomainQAToRestQA() failed: %v", err)
		}

		result = append(result, qa)
	}

	return result, nil
}

func convertDomainTextToRestText(obj *domainquiz.Text) (*restquiz.Text, error) {
	var result restquiz.Text
	result.Text = obj.Text
	result.IsHtml = obj.IsHtml

	return &result, nil
}

func convertDomainQuestionToRestQuestion(obj *domainquiz.Question) (*restquiz.Question, error) {
	var result restquiz.Question
	result.Id = obj.Id
	result.Link = obj.Link

	text, err := convertDomainTextToRestText(&obj.Text)
	if err != nil {
		return nil, fmt.Errorf("convertDomainTextToRestText() failed: %v", err)
	}

	result.Text = *text

	return &result, nil
}

func convertDomainQAToRestQA(obj *domainquiz.QuestionAndAnswer) (*restquiz.QuestionAndAnswer, error) {
	var result restquiz.QuestionAndAnswer

	question, err := convertDomainQuestionToRestQuestion(&obj.Question)
	if err != nil {
		return nil, fmt.Errorf("convertDomainSectionToRestSection() failed: %v", err)
	}

	result.Question = *question

	answer, err := convertDomainTextToRestText(&obj.Answer)
	if err != nil {
		return nil, fmt.Errorf("convertDomainTextToRestText() failed for answer: %v", err)
	}

	result.Answer = *answer

	return &result, nil
}

func convertDomainSubSectionToRestSubSection(obj *domainquiz.SubSection) (*restquiz.SubSection, error) {
	var result restquiz.SubSection

	hasIdAndTitle, err := convertDomainHasIdAndTitleToRestHasIdAndTitle(&obj.HasIdAndTitle)
	if err != nil {
		return nil, fmt.Errorf("convertDomainHasIdAndTitleToRestHasIdAndTitle() failed: %v", err)
	}

	result.HasIdAndTitle = *hasIdAndTitle

	result.Questions, err = convertDomainQuestionsToRestQuestions(obj.Questions)
	if err != nil {
		return nil, fmt.Errorf("convertDomainQuestionsToRestQuestions() failed for answer: %v", err)
	}

	result.AnswersAsChoices = obj.AnswersAsChoices

	return &result, nil
}

func convertDomainSectionToRestSection(obj *domainquiz.Section) (*restquiz.Section, error) {
	var result restquiz.Section

	hasIdAndTitle, err := convertDomainHasIdAndTitleToRestHasIdAndTitle(&obj.HasIdAndTitle)
	if err != nil {
		return nil, fmt.Errorf("convertDomainHasIdAndTitleToRestHasIdAndTitle() failed: %v", err)
	}

	result.HasIdAndTitle = *hasIdAndTitle

	result.Questions, err = convertDomainQuestionsToRestQuestions(obj.Questions)
	if err != nil {
		return nil, fmt.Errorf("convertDomainQuestionsToRestQuestions() failed: %v", err)
	}

	for _, dtoSubSection := range obj.SubSections {
		subSection, err := convertDomainSubSectionToRestSubSection(dtoSubSection)
		if err != nil {
			return nil, fmt.Errorf("convertDomainSubSectionToRestSubSection() failed: %v", err)
		}

		result.SubSections = append(result.SubSections, subSection)
	}

	for _, dtoText := range obj.DefaultChoices {
		defaultChoice, err := convertDomainTextToRestText(dtoText)
		if err != nil {
			return nil, fmt.Errorf("convertDomainTextToRestText() failed for DefaultChoices: %v", err)
		}

		result.DefaultChoices = append(result.DefaultChoices, defaultChoice)
	}

	result.AnswersAsChoices = obj.AnswersAsChoices

	return &result, nil
}
