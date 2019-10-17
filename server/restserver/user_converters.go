package restserver

import (
	"fmt"
	domainquiz "github.com/murraycu/go-bigoquiz-server/domain/quiz"
	domainuser "github.com/murraycu/go-bigoquiz-server/domain/user"
	restuser "github.com/murraycu/go-bigoquiz-server/server/restserver/user"
)

func convertDomainStatsToRestStats(stats *domainuser.Stats, quiz *domainquiz.Quiz) (*restuser.Stats, error) {
	questionHistories, err := convertDomainQuestionHistoriesToRestQuestionHistories(stats.QuestionHistories, quiz)
	if err != nil {
		return nil, fmt.Errorf("convertDomainQuestionHistoriesToRestQuestionHistories() failed: %v", err)
	}

	return &restuser.Stats{
		QuizId:                     stats.QuizId,
		SectionId:                  stats.SectionId,
		Answered:                   stats.Answered,
		Correct:                    stats.Correct,
		CountQuestionsAnsweredOnce: stats.CountQuestionsAnsweredOnce,
		CountQuestionsCorrectOnce:  stats.CountQuestionsCorrectOnce,
		QuestionHistories:          questionHistories,
	}, nil
}

func convertDomainQuestionHistoriesToRestQuestionHistories(questionHistories []domainuser.QuestionHistory, quiz *domainquiz.Quiz) ([]restuser.QuestionHistory, error) {
	var result []restuser.QuestionHistory

	for _, qh := range questionHistories {
		qa := quiz.GetQuestionAndAnswer(qh.QuestionId)
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

func convertDomainQuestionHistoryToRestQuestionHistory(obj domainuser.QuestionHistory, question *domainquiz.Question) (restuser.QuestionHistory, error) {
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
