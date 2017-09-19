package quiz

import (
	"testing"
)

func TestCreateReverse(t *testing.T) {
	const TEST_ID = "someid"
	const TEST_QUESTION = "somequestion"
	const TEST_ANSWER = "someanswer"

	var qa QuestionAndAnswer
	qa.Id = TEST_ID
	qa.Text.Text = TEST_QUESTION
	qa.Answer.Text = TEST_ANSWER

	reverse := qa.createReverse()
	if reverse == nil {
		t.Error("createReverse() returned nil.")
	}

	if reverse.Id != ("reverse-" + TEST_ID) {
		t.Error("The reverse QuestionAndAnswer did not have the expected ID.")
	}

	if reverse.Text.Text != TEST_ANSWER {
		t.Error("The reverse QuestionAndAnswer did not have the expected title.")
	}

	if reverse.Answer.Text != TEST_QUESTION {
		t.Error("The reverse QuestionAndAnswer did not have the expected answer.")
	}
}

