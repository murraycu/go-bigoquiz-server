package restserver

import (
	"fmt"
	domainuser "github.com/murraycu/go-bigoquiz-server/domain/user"
	restquiz "github.com/murraycu/go-bigoquiz-server/server/restserver/quiz"
	restuser "github.com/murraycu/go-bigoquiz-server/server/restserver/user"
)

func convertDomainStatsToRestStats(stats *domainuser.Stats, quizCache *QuizCache) (*restuser.Stats, error) {
	questionHistories, err := convertDomainQuestionHistoriesToRestQuestionHistories(stats.QuestionHistories, quizCache)
	if err != nil {
		return nil, fmt.Errorf("convertDomainQuestionHistoriesToRestQuestionHistories() failed: %v", err)
	}

	section, err := quizCache.GetSection(stats.SectionId)
	if err != nil {
		return nil, fmt.Errorf("GetSection() failed: %v", err)
	}

	var sectionTitle string
	if section != nil {
		sectionTitle = section.Title
	}

	questionsCount := 0
	if len(stats.SectionId) == 0 {
		// This should be stats for a whole quiz.
		questionsCount = quizCache.GetQuestionsCount()
	} else {
		// This should be stats for just one section.
		questionsCount = quizCache.GetSectionQuestionsCount(stats.SectionId)
	}

	return &restuser.Stats{
		QuizId:                     stats.QuizId,
		SectionId:                  stats.SectionId,
		Answered:                   stats.Answered,
		Correct:                    stats.Correct,
		CountQuestionsAnsweredOnce: stats.CountQuestionsAnsweredOnce,
		CountQuestionsCorrectOnce:  stats.CountQuestionsCorrectOnce,
		QuestionHistories:          questionHistories,

		CountQuestions: questionsCount,
		QuizTitle:      quizCache.Quiz.Title,
		SectionTitle:   sectionTitle,
	}, nil
}

func convertDomainQuestionHistoriesToRestQuestionHistories(questionHistories []domainuser.QuestionHistory, quizCache *QuizCache) ([]restuser.QuestionHistory, error) {
	var result []restuser.QuestionHistory

	for _, qh := range questionHistories {
		qa := quizCache.GetQuestionAndAnswer(qh.QuestionId)
		if qa == nil {
			continue
		}

		restQH, err := convertDomainQuestionHistoryToRestQuestionHistory(qh, &qa.Question)
		if err != nil {
			return nil, fmt.Errorf("convertDomainQuestionHistoryToRestQuestionHistory() failed: %v", err)
		}

		result = append(result, restQH)
	}

	return result, nil
}

func convertDomainQuestionHistoryToRestQuestionHistory(obj domainuser.QuestionHistory, question *restquiz.Question) (restuser.QuestionHistory, error) {
	var subSectionTitle string
	if question.SubSection != nil {
		subSectionTitle = question.SubSection.Title
	}

	return restuser.QuestionHistory{
		QuestionId:            obj.QuestionId,
		AnsweredCorrectlyOnce: obj.AnsweredCorrectlyOnce,
		CountAnsweredWrong:    obj.CountAnsweredWrong,

		// Extras, which are in the REST struct, but not in the domain struct.
		QuestionTitle:   &question.Text,
		SectionId:       question.SectionId,
		SubSectionTitle: subSectionTitle,
	}, nil
}
