package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const TEST_QUESTION_ID = "test-question-id"

func TestStatsCountAnsweredWrong(t *testing.T) {
	var stats Stats

	assert.Equal(t, 0, stats.GetQuestionCountAnsweredWrong(TEST_QUESTION_ID))

	stats.UpdateProblemQuestion(TEST_QUESTION_ID, false)
	assert.Equal(t, 1, stats.GetQuestionCountAnsweredWrong(TEST_QUESTION_ID))

	stats.UpdateProblemQuestion(TEST_QUESTION_ID, false)
	assert.Equal(t, 2, stats.GetQuestionCountAnsweredWrong(TEST_QUESTION_ID))

	stats.UpdateProblemQuestion(TEST_QUESTION_ID, true)
	assert.Equal(t, 1, stats.GetQuestionCountAnsweredWrong(TEST_QUESTION_ID))
}