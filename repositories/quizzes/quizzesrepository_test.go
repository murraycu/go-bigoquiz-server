package quizzes

import (
	"path/filepath"
	"testing"

	domainquiz "github.com/murraycu/go-bigoquiz-server/domain/quiz"
	"github.com/stretchr/testify/assert"
)

// TODO: Add tests that don't use the real quizzes.

func loadRealQuizzes(t *testing.T) MapQuizzes {
	directoryFilepath, err := filepath.Abs("../../quizzes")
	assert.Nil(t, err)
	assert.NotNil(t, directoryFilepath)

	quizzesStore, err := NewQuizzesRepository(directoryFilepath)
	assert.Nil(t, err)
	assert.NotNil(t, quizzesStore)

	quizzes, err := quizzesStore.LoadQuizzes()
	assert.Nil(t, err)
	assert.NotNil(t, quizzes)

	return quizzes
}

func TestQuizzesRepositoryNewWithRealQuizzes(t *testing.T) {
	restQuizzes := loadRealQuizzes(t)

	for quizId, quiz := range restQuizzes {
		assert.NotNil(t, quizId)
		assert.NotNil(t, quiz)
	}
}

func TestQuizzesRepositoryRealQuizHasQuestions(t *testing.T) {
	restQuizzes := loadRealQuizzes(t)

	for _, quiz := range restQuizzes {
		assert.NotNil(t, quiz)

		questionsFound := false
		if len(quiz.Questions) > 0 {
			questionsFound = true
		} else {
			for _, section := range quiz.Sections {
				if len(section.Questions) > 0 {
					questionsFound = true
					break
				}

				for _, subSection := range section.SubSections {
					if len(subSection.Questions) > 0 {
						questionsFound = true
						break
					}
				}
			}
		}

		assert.True(t, questionsFound)
	}
}

func questionHasAnswer(question *domainquiz.QuestionAndAnswer) bool {
	return len(question.Answer.Text) > 0
}

func TestQuizzesRepositoryRealQuizHasAnswers(t *testing.T) {
	restQuizzes := loadRealQuizzes(t)

	for _, quiz := range restQuizzes {
		assert.NotNil(t, quiz)

		for _, question := range quiz.Questions {
			assert.True(t, questionHasAnswer(question))
		}

		for _, section := range quiz.Sections {
			for _, question := range section.Questions {
				assert.True(t, questionHasAnswer(question))
			}

			for _, subSection := range section.SubSections {
				for _, question := range subSection.Questions {
					assert.True(t, questionHasAnswer(question))
				}
			}
		}
	}
}

func TestQuizzesRepositoryRealQuizHasSections(t *testing.T) {
	restQuizzes := loadRealQuizzes(t)

	for _, quiz := range restQuizzes {
		assert.NotNil(t, quiz)

		assert.NotEmpty(t, quiz.Sections)
	}
}

func TestQuizzesRepositoryRealQuizHasTitle(t *testing.T) {
	restQuizzes := loadRealQuizzes(t)

	for _, quiz := range restQuizzes {
		assert.NotNil(t, quiz)

		assert.NotEmpty(t, quiz.Title)
	}
}
