package restserver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testQuizCache(t *testing.T) *QuizCache {
	quiz := testRestQuiz()
	quizCache, err := NewQuizCache(quiz)
	assert.Nil(t, err)
	return quizCache
}

func TestQuizCacheRandomQuestion(t *testing.T) {
	quizCache := testQuizCache(t)

	randomQuestion := quizCache.GetRandomQuestion("")
	assert.NotNil(t, randomQuestion)
}

func TestQuizCacheRandomQuestionForSection(t *testing.T) {
	quizCache := testQuizCache(t)

	sectionId := quizCache.Quiz.Sections[1].Id

	randomQuestion := quizCache.GetRandomQuestion(sectionId)
	assert.NotNil(t, randomQuestion)
}

func TestQuizCacheQuestionCount(t *testing.T) {
	quizCache := testQuizCache(t)

	assert.NotZero(t, quizCache.GetQuestionsCount())
}

func TestQuizCacheGetQuestionAndAnswer(t *testing.T) {
	quizCache := testQuizCache(t)

	questionId := quizCache.Quiz.Sections[0].Questions[1].Id

	qa := quizCache.GetQuestionAndAnswer(questionId)
	assert.NotNil(t, qa)
}

func TestQuizCacheGetQuestionAndAnswerFromSection(t *testing.T) {
	quizCache := testQuizCache(t)

	questionId := quizCache.Quiz.Sections[1].Questions[1].Id

	qa := quizCache.GetQuestionAndAnswer(questionId)
	assert.NotNil(t, qa)
}

func TestQuizCacheGetAnswer(t *testing.T) {
	quizCache := testQuizCache(t)

	questionId := quizCache.Quiz.Sections[0].Questions[1].Id

	answer := quizCache.GetAnswer(questionId)
	assert.NotNil(t, answer)
}

func TestQuizCacheGetAnswerFromSection(t *testing.T) {
	quizCache := testQuizCache(t)

	questionId := quizCache.Quiz.Sections[1].Questions[1].Id

	answer := quizCache.GetAnswer(questionId)
	assert.NotNil(t, answer)
}

func TestQuizCacheGetSection(t *testing.T) {
	quizCache := testQuizCache(t)

	sectionId := quizCache.Quiz.Sections[1].Id

	section, err := quizCache.GetSection(sectionId)
	assert.Nil(t, err)
	assert.NotNil(t, section)
}

func TestQuizCacheGetSubSection(t *testing.T) {
	quizCache := testQuizCache(t)

	sectionId := quizCache.Quiz.Sections[1].Id
	subSectionId := quizCache.Quiz.Sections[1].SubSections[1].Id

	subSection := quizCache.GetSubSection(sectionId, subSectionId)
	assert.NotNil(t, subSection)
}
