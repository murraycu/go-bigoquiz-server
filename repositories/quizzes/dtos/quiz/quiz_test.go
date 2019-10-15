package quiz

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
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
	if err != nil {
		t.Error("Could not find file.", err)
	}

	q, err := LoadQuiz(absFilePath, id)
	if err == nil {
		t.Error("LoadQuiz() did not fail with an error.")
	}

	if q != nil {
		t.Error("LoadQuiz() returned a quiz when it should have returned nil.")
	}
}

func loadQuiz(t *testing.T, quizId string) *Quiz {
	absFilePath, err := filepath.Abs("../../../../quizzes/" + quizId + ".xml")
	if err != nil {
		t.Error("Could not find file.", err)
	}
	q, err := LoadQuiz(absFilePath, quizId)
	if err != nil {
		t.Error("LoadQuiz() failed.", err)
	}
	return q
}

func TestLoadQuiz(t *testing.T) {
	q := loadQuiz(t, "bigo")

	if q.Sections == nil {
		t.Error("The quiz has no sections.")
	}

	const SECTION_ID = "data-structure-operations"
	section := getSection(q, SECTION_ID)
	if section == nil {
		t.Error("The quiz does not have the expected section.")
	}

	assert.Equal(t, SECTION_ID, section.Id)
	assert.Equal(t, "Data Structure Operations", section.Title)
	assert.Equal(t, "", section.Link)

	const QUESTION_ID = "b-tree-search-worst"
	qa := getQuestionAndAnswer(q, QUESTION_ID) // (Gets it from a sub-section)
	if qa == nil {
		t.Error("The quiz does not have the expected question.")
	}

	assert.Equal(t, QUESTION_ID, qa.Id)

	assert.Equal(t, "Search (Worst)", qa.Text.Text)
	assert.Equal(t, false, qa.Text.IsHtml)

	assert.Equal(t, "O(log(n))", qa.Answer.Text)
}

func TestLoadQuizWithReverseSection(t *testing.T) {
	q := loadQuiz(t, "datastructures")

	if q.Sections == nil {
		t.Error("The quiz has no sections.")
	}

	const SECTION_ID = "reverse-datastructures-hash-tables"
	section := getSection(q, SECTION_ID)
	if section == nil {
		t.Error("The quiz does not have the expected reverse section.")
	}

	assert.Equal(t, SECTION_ID, section.Id)
	assert.Equal(t, "Reverse: Hash Tables", section.Title)
	assert.Equal(t, "https://en.wikipedia.org/wiki/Hash_table", section.Link)

	const QUESTION_ID = "reverse-datastructures-hash-tables-open-addressing-strategy-probe-sequence"
	qa := getQuestionAndAnswer(q, QUESTION_ID)
	if qa == nil {
		t.Error("The quiz does not have the expected reverse question.")
	}

	assert.Equal(t, QUESTION_ID, qa.Id)

	assert.Equal(t, "The buckets are examined, starting with the hashed-to slot and proceeding in some probe sequence, until an unoccupied slot is found.", qa.Text.Text)
	assert.Equal(t, false, qa.Text.IsHtml)

	assert.Equal(t, "Probe sequence", qa.Answer.Text)
}
