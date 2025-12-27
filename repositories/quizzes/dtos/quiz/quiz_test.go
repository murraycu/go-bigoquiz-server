package quiz

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getQuestionAndAnswer(q *Quiz, questionId string) *QuestionAndAnswer {
	for _, qa := range q.Questions {
		if qa.Id == questionId {
			return qa
		}
	}

	// Look in Sections:
	for _, s := range q.Sections {
		for _, qa := range s.Questions {
			if qa.Id == questionId {
				return qa
			}
		}

		// Look in SubSections:
		for _, ss := range s.SubSections {
			for _, qa := range ss.Questions {
				if qa.Id == questionId {
					return qa
				}
			}
		}
	}

	return nil
}

func getSection(q *Quiz, sectionId string) *Section {
	for _, s := range q.Sections {
		if s.Id == sectionId {
			return s
		}
	}

	return nil
}

func TestLoadQuizWithBadFilepath(t *testing.T) {
	id := "doesnotexist"
	absFilePath, err := filepath.Abs("../../quizzes/" + id + ".xml")
	assert.Nil(t, err)
	assert.NotNil(t, absFilePath)

	q, err := LoadQuiz(absFilePath /*asJson=*/, false /*addReverse=*/, true, id)
	assert.NotNil(t, err)
	assert.Nil(t, q)
}

func loadQuiz(t *testing.T, quizId string) *Quiz {
	absFilePath, err := filepath.Abs("../../../../quizzes/" + quizId + ".xml")
	assert.Nil(t, err)
	assert.NotNil(t, absFilePath)

	q, err := LoadQuiz(absFilePath /*asJson=*/, false /*addReverses=*/, true, quizId)
	assert.Nil(t, err)
	assert.NotNil(t, q)

	return q
}

func TestLoadQuiz(t *testing.T) {
	q := loadQuiz(t, "bigo")

	assert.NotNil(t, q.Sections)

	const SECTION_ID = "data-structure-operations"
	section := getSection(q, SECTION_ID)
	assert.NotNil(t, section)

	assert.Equal(t, SECTION_ID, section.Id)
	assert.Equal(t, "Data Structure Operations", section.Title)
	assert.Equal(t, "", section.Link)

	const QUESTION_ID = "b-tree-search-worst"
	qa := getQuestionAndAnswer(q, QUESTION_ID) // (Gets it from a sub-section)
	assert.NotNil(t, qa)

	assert.Equal(t, QUESTION_ID, qa.Id)

	assert.Equal(t, "Search (Worst)", qa.TextDetail.Text)
	assert.Equal(t, false, qa.TextDetail.IsHtml)

	assert.Equal(t, "O(log(n))", qa.Answer.Text)
}

func TestLoadQuizWithReverseSection(t *testing.T) {
	q := loadQuiz(t, "datastructures")

	assert.NotNil(t, q.Sections)

	const SECTION_ID = "reverse-datastructures-hash-tables"
	section := getSection(q, SECTION_ID)
	assert.NotNil(t, section)

	assert.Equal(t, SECTION_ID, section.Id)
	assert.Equal(t, "Reverse: Hash Tables", section.Title)
	assert.Equal(t, "https://en.wikipedia.org/wiki/Hash_table", section.Link)

	const QUESTION_ID = "reverse-datastructures-hash-tables-open-addressing-strategy-probe-sequence"
	qa := getQuestionAndAnswer(q, QUESTION_ID)
	assert.NotNil(t, qa)

	assert.Equal(t, QUESTION_ID, qa.Id)

	assert.Equal(t, "The buckets are examined, starting with the hashed-to slot and proceeding in some probe sequence, until an unoccupied slot is found.", qa.TextDetail.Text)
	assert.Equal(t, false, qa.TextDetail.IsHtml)

	assert.Equal(t, "Probe sequence", qa.Answer.Text)
}

func TestLoadQuizWithNoSections(t *testing.T) {
	q := loadQuiz(t, "bitwise")

	// loadQuiz() creates a section for the top-level questions.
	// (TODO: Just require this in the original .xml?)
	assert.NotNil(t, q.Sections)

	section := q.Sections[0]
	assert.NotNil(t, section)
	assert.Equal(t, q.Title, section.Title)

	const QUESTION_ID = "bitwise-divide-by-2"
	qa := getQuestionAndAnswer(q, QUESTION_ID)
	assert.NotNil(t, qa)

	assert.Equal(t, QUESTION_ID, qa.Id)
}
