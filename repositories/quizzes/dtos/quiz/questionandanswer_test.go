package quiz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateReverseWithTextDetail(t *testing.T) {
	const TEST_ID = "someid"
	const TEST_QUESTION = "somequestion"
	const TEST_ANSWER = "someanswer"

	var qa QuestionAndAnswer
	qa.Id = TEST_ID
	qa.TextDetail.Text = TEST_QUESTION
	qa.AnswerDetail.Text = TEST_ANSWER

	reverse := qa.createReverse()
	assert.NotNil(t, reverse)

	assert.Equal(t, "reverse-"+TEST_ID, reverse.Id)
	assert.Equal(t, TEST_ANSWER, reverse.TextDetail.Text)
	assert.Equal(t, TEST_QUESTION, reverse.AnswerDetail.Text)
}

func TestCreateReverseWithTextSimple(t *testing.T) {
	const TEST_ID = "someid"
	const TEST_QUESTION = "somequestion"
	const TEST_ANSWER = "someanswer"

	var qa QuestionAndAnswer
	qa.Id = TEST_ID
	qa.TextSimple = TEST_QUESTION
	qa.AnswerDetail.Text = TEST_ANSWER

	reverse := qa.createReverse()
	assert.NotNil(t, reverse)

	assert.Equal(t, "reverse-"+TEST_ID, reverse.Id)
	assert.Equal(t, TEST_ANSWER, reverse.TextDetail.Text)
	assert.Equal(t, TEST_QUESTION, reverse.AnswerDetail.Text)
}
