package restserver

import (
	"github.com/murraycu/go-bigoquiz-server/repositories/quizzes"
	"github.com/stretchr/testify/assert"
	"path/filepath"
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

func loadRealRestQuizzes(t *testing.T) restQuizMap {
	directoryFilepath, err := filepath.Abs("../../quizzes")
	assert.Nil(t, err)
	assert.NotNil(t, directoryFilepath)

	quizzesStore, err := quizzes.NewQuizzesRepository(directoryFilepath)
	assert.Nil(t, err)
	assert.NotNil(t, quizzesStore)

	quizzes, err := quizzesStore.LoadQuizzes()
	assert.Nil(t, err)
	assert.NotNil(t, quizzes)

	restQuizzes, err := convertDomainQuizzesToRestQuizzes(quizzes)
	assert.Nil(t, err)
	assert.NotNil(t, quizzes)

	return restQuizzes
}

func TestQuizCacheNewWithRealQuizzes(t *testing.T) {
	restQuizzes := loadRealRestQuizzes(t)

	for _, quiz := range restQuizzes {
		quizCache, err := NewQuizCache(quiz)
		assert.Nil(t, err)
		assert.NotNil(t, quizCache)
	}
}

func TestQuizCacheRealQuizHasQuestions(t *testing.T) {
	restQuizzes := loadRealRestQuizzes(t)

	for _, quiz := range restQuizzes {
		quizCache, err := NewQuizCache(quiz)
		assert.Nil(t, err)
		assert.NotNil(t, quizCache)

		assert.NotZero(t, quizCache.GetQuestionsCount())
	}
}

/* TODO: Check if there are reverse sections, if there should be reverse sections.
func TestQuizCacheRealQuizHasReverseSections(t *testing.T) {
	restQuizzes := loadRealRestQuizzes(t)

	for _, quiz := range restQuizzes {
		quizCache, err := NewQuizCache(quiz)
		assert.Nil(t, err)
		assert.NotNil(t, quizCache)

		reverseSectionFound := false
		for _, section := range quiz.Sections {
			if strings.HasPrefix(section.Id, "reverse-") {
				reverseSectionFound = true
				break
			}
		}

		assert.True(t, reverseSectionFound)
	}
}
*/

func TestQuizCacheRealQuizFound(t *testing.T) {
	restQuizzes := loadRealRestQuizzes(t)
	graph, ok := restQuizzes["graphs"]
	assert.True(t, ok)
	assert.NotNil(t, graph)
}

func TestQuizCacheRealQuizShouldNotHaveReverseSections(t *testing.T) {
	restQuizzes := loadRealRestQuizzes(t)

	// TODO: This depends on knowledge of the real quiz.
	quiz, ok := restQuizzes["compilers"]
	assert.True(t, ok)
	assert.NotNil(t, quiz)

	quizCache, err := NewQuizCache(quiz)
	assert.Nil(t, err)
	assert.NotNil(t, quizCache)

	section, err := quizCache.GetSection("reverse-compilers-structure")
	assert.Nil(t, section)
}

func TestQuizCacheRealQuizShouldHaveReverseSections(t *testing.T) {
	restQuizzes := loadRealRestQuizzes(t)

	// TODO: This depends on knowledge of the real quiz.
	quiz, ok := restQuizzes["graphs"]
	assert.True(t, ok)
	assert.NotNil(t, quiz)

	quizCache, err := NewQuizCache(quiz)
	assert.Nil(t, err)
	assert.NotNil(t, quizCache)

	section, err := quizCache.GetSection("reverse-graph-algorithm-definitions-shortest-path")
	assert.Nil(t, err)
	assert.NotNil(t, section)
}

func TestQuizCacheRealQuizHasChoices(t *testing.T) {
	restQuizzes := loadRealRestQuizzes(t)

	for _, quiz := range restQuizzes {
		quizCache, err := NewQuizCache(quiz)
		assert.Nil(t, err)
		assert.NotNil(t, quizCache)

		err = fillRestQuizExtrasFromQuizCache(quiz, quizCache)
		assert.Nil(t, err)

		for _, section := range quiz.Sections {
			if section.AnswersAsChoices {
				for _, question := range section.Questions {
					assert.NotEmpty(t, question.Question.Choices)
				}
			}

			for _, subSection := range section.SubSections {
				if subSection.AnswersAsChoices {
					for _, question := range subSection.Questions {
						assert.NotEmpty(t, question.Question.Choices)
					}
				}
			}
		}
	}
}
