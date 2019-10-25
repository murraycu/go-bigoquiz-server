package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const TEST_QUESTION_ID = "test-question-id"

func TestStatsCountAnsweredWrong(t *testing.T) {
	var stats Stats

	assert.Equal(t, 0, stats.GetQuestionCountAnsweredWrong(TEST_QUESTION_ID))

	stats.UpdateStatsForAnswerCorrectness(TEST_QUESTION_ID, false)
	assert.Equal(t, 1, stats.GetQuestionCountAnsweredWrong(TEST_QUESTION_ID))

	stats.UpdateStatsForAnswerCorrectness(TEST_QUESTION_ID, false)
	assert.Equal(t, 2, stats.GetQuestionCountAnsweredWrong(TEST_QUESTION_ID))

	stats.UpdateStatsForAnswerCorrectness(TEST_QUESTION_ID, true)
	assert.Equal(t, 1, stats.GetQuestionCountAnsweredWrong(TEST_QUESTION_ID))
}

func TestStatsGetQuestionWasAnsweredWithWrong(t *testing.T) {
	var stats Stats

	assert.False(t, stats.GetQuestionWasAnswered(TEST_QUESTION_ID))

	stats.UpdateStatsForAnswerCorrectness(TEST_QUESTION_ID, false)
	assert.True(t, stats.GetQuestionWasAnswered(TEST_QUESTION_ID))
}

func TestStatsGetQuestionWasAnsweredWithCorrect(t *testing.T) {
	var stats Stats

	assert.False(t, stats.GetQuestionWasAnswered(TEST_QUESTION_ID))

	stats.UpdateStatsForAnswerCorrectness(TEST_QUESTION_ID, true)
	assert.True(t, stats.GetQuestionWasAnswered(TEST_QUESTION_ID))
}

func TestStatsGetQuestionCountAnsweredWrong(t *testing.T) {
	var stats Stats

	assert.Equal(t, 0, stats.GetQuestionCountAnsweredWrong(TEST_QUESTION_ID))

	stats.UpdateStatsForAnswerCorrectness(TEST_QUESTION_ID, true)
	assert.Equal(t, -1, stats.GetQuestionCountAnsweredWrong(TEST_QUESTION_ID))

	stats.UpdateStatsForAnswerCorrectness(TEST_QUESTION_ID, false)
	assert.Equal(t, 0, stats.GetQuestionCountAnsweredWrong(TEST_QUESTION_ID))
}

func TestStatsIncrementAnswered(t *testing.T) {
	var stats Stats

	assert.Equal(t, 0, stats.Answered)

	stats.UpdateStatsForAnswerCorrectness(TEST_QUESTION_ID, true)
	assert.Equal(t, 1, stats.Answered)

	stats.UpdateStatsForAnswerCorrectness(TEST_QUESTION_ID, true)
	assert.Equal(t, 2, stats.Answered)

	stats.UpdateStatsForAnswerCorrectness(TEST_QUESTION_ID, false)
	assert.Equal(t, 3, stats.Answered)
}

func TestStatsCorrect(t *testing.T) {
	var stats Stats

	assert.Equal(t, 0, stats.Correct)

	stats.UpdateStatsForAnswerCorrectness(TEST_QUESTION_ID, true)
	assert.Equal(t, 1, stats.Correct)

	stats.UpdateStatsForAnswerCorrectness(TEST_QUESTION_ID, true)
	assert.Equal(t, 2, stats.Correct)

	stats.UpdateStatsForAnswerCorrectness(TEST_QUESTION_ID, false)
	assert.Equal(t, 2, stats.Correct)
}

func TestStatsIncrementAnsweredOnce(t *testing.T) {
	var stats Stats

	assert.Equal(t, 0, stats.CountQuestionsAnsweredOnce)

	stats.UpdateStatsForAnswerCorrectness("some-question-1", true)
	assert.Equal(t, 1, stats.CountQuestionsAnsweredOnce)

	stats.UpdateStatsForAnswerCorrectness("some-question-2", true)
	assert.Equal(t, 2, stats.CountQuestionsAnsweredOnce)

	stats.UpdateStatsForAnswerCorrectness("some-question-3", false)
	assert.Equal(t, 3, stats.CountQuestionsAnsweredOnce)
}

func TestStatsIncrementAnsweredCorrectOnce(t *testing.T) {
	var stats Stats

	assert.Equal(t, 0, stats.CountQuestionsCorrectOnce)

	// This should only increment for correct answers, and only for questions that have not been correct before.
	stats.UpdateStatsForAnswerCorrectness("some-question-1", true)
	assert.Equal(t, 1, stats.CountQuestionsCorrectOnce)

	stats.UpdateStatsForAnswerCorrectness("some-question-2", true)
	assert.Equal(t, 2, stats.CountQuestionsCorrectOnce)

	stats.UpdateStatsForAnswerCorrectness("some-question-2", true)
	assert.Equal(t, 2, stats.CountQuestionsCorrectOnce)

	stats.UpdateStatsForAnswerCorrectness("some-question-3", false)
	assert.Equal(t, 2, stats.CountQuestionsCorrectOnce)

	stats.UpdateStatsForAnswerCorrectness("some-question-4", true)
	assert.Equal(t, 3, stats.CountQuestionsCorrectOnce)
}
