package db

import (
	"cloud.google.com/go/datastore"
	"fmt"
	domainuser "github.com/murraycu/go-bigoquiz-server/domain/user"
	dtouser "github.com/murraycu/go-bigoquiz-server/repositories/db/dtos/user"
)

func convertDtoQuestionHistoryToDomainQuestionHistory(dto dtouser.QuestionHistory) *domainuser.QuestionHistory {
	return &domainuser.QuestionHistory{
		QuestionId:            dto.QuestionId,
		AnsweredCorrectlyOnce: dto.AnsweredCorrectlyOnce,
		CountAnsweredWrong:    dto.CountAnsweredWrong}
}

func convertDtoStatsToDomainStats(dto *dtouser.Stats) *domainuser.Stats {
	var result domainuser.Stats

	result.QuizId = dto.QuizId
	result.SectionId = dto.SectionId

	result.Answered = dto.Answered
	result.Correct = dto.Correct
	result.CountQuestionsAnsweredOnce = dto.CountQuestionsAnsweredOnce
	result.CountQuestionsCorrectOnce = dto.CountQuestionsCorrectOnce

	for _, qh := range dto.QuestionHistories {
		result.QuestionHistories = append(result.QuestionHistories,
			*convertDtoQuestionHistoryToDomainQuestionHistory(qh))
	}

	// TODO: Fill extras, which are not in the dot?
	/*
		result.CountQuestions = dto.CountQuestions
		result.QuizTitle = dto.QuizTitle
		result.SectionTitle = dto.SectionTitle
	*/

	return &result
}

func convertDomainQuestionHistoryToDtoQuestionHistory(history domainuser.QuestionHistory) *dtouser.QuestionHistory {
	return &dtouser.QuestionHistory{
		QuestionId:            history.QuestionId,
		AnsweredCorrectlyOnce: history.AnsweredCorrectlyOnce,
		CountAnsweredWrong:    history.CountAnsweredWrong,
	}

	// Fill these?
	/*
		result.QuestionTitle = history.QuestionTitle
		result.SectionId = history.SectionId
		result.SubSectionTitle = history.SubSectionTitle
		return &result
	*/
}

func convertDomainStatsToDtoStats(stats *domainuser.Stats, userID string) (*dtouser.Stats, error) {
	var result dtouser.Stats

	userId, err := datastore.DecodeKey(userID)
	if err != nil {
		return nil, fmt.Errorf("datastore,DecodeKey() failed: %v", err)
	}

	result.UserId = userId

	result.QuizId = stats.QuizId
	result.SectionId = stats.SectionId

	result.Answered = stats.Answered
	result.Correct = stats.Correct
	result.CountQuestionsAnsweredOnce = stats.CountQuestionsAnsweredOnce
	result.CountQuestionsCorrectOnce = stats.CountQuestionsCorrectOnce

	for _, qh := range stats.QuestionHistories {
		result.QuestionHistories = append(result.QuestionHistories,
			*convertDomainQuestionHistoryToDtoQuestionHistory(qh))

		// TODO: fill extras?
	}

	// TODO: fill extras?
	/*
		result.CountQuestions = stats.CountQuestions
		result.QuizTitle = stats.QuizTitle
		result.SectionTitle = stats.SectionTitle
	*/

	return &result, nil
}

func convertDtoProfileToDomainProfile(dto *dtouser.Profile) *domainuser.Profile {
	return &domainuser.Profile{
		Name:  dto.Name,
		Email: dto.Email,

		GoogleProfileUrl:   dto.GoogleProfileUrl,
		GitHubProfileUrl:   dto.GitHubProfileUrl,
		FacebookProfileUrl: dto.FacebookProfileUrl,
	}

}
