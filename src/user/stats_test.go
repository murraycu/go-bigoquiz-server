package user

import (
	"quiz"
	"testing"
)

const TEST_QUESTION_ID = "test-question-id"

func TestStatsCountAnsweredWrong(t *testing.T) {
	var question quiz.Question
	question.Id = TEST_QUESTION_ID
	var stats Stats
	if stats.GetQuestionCountAnsweredWrong(TEST_QUESTION_ID) != 0 {
		t.Error("stats.GetQuestionCountAnsweredWrong() did not default to 0.")
	}

	stats.UpdateProblemQuestion(&question, false)
	if stats.GetQuestionCountAnsweredWrong(TEST_QUESTION_ID) != 1 {
		t.Error("stats.GetQuestionCountAnsweredWrong() did not increase to 1.")
	}

	stats.UpdateProblemQuestion(&question, false)
	if stats.GetQuestionCountAnsweredWrong(TEST_QUESTION_ID) != 2 {
		t.Errorf("stats.GetQuestionCountAnsweredWrong() did not increase to 2. Instead it is %v", stats.GetQuestionCountAnsweredWrong(TEST_QUESTION_ID))
	}

	stats.UpdateProblemQuestion(&question, true)
	if stats.GetQuestionCountAnsweredWrong(TEST_QUESTION_ID) != 1 {
		t.Error("stats.GetQuestionCountAnsweredWrong() did not decrease to 1.")
	}
}
