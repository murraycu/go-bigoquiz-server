package quiz

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

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
	absFilePath, err := filepath.Abs("../../quizzes/" + quizId + ".xml")
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
	section := q.GetSection(SECTION_ID)
	if section == nil {
		t.Error("The quiz does not have the expected section.")
	}

	assert.Equal(t, SECTION_ID, section.Id)
	assert.Equal(t, "Data Structure Operations", section.Title)
	assert.Equal(t, "", section.Link)

	const QUESTION_ID = "b-tree-search-worst"
	qa := q.GetQuestionAndAnswer(QUESTION_ID)
	if qa == nil {
		t.Error("The quiz does not have the expected question.")
	}

	assert.Equal(t, QUESTION_ID, qa.Id)
	assert.Equal(t, SECTION_ID, qa.SectionId)
	assert.Equal(t, SECTION_ID, qa.Section.Id)
	assert.Equal(t, "Data Structure Operations", qa.Section.Title)

	const SUB_SECTION_ID = "b-tree"
	assert.Equal(t, SUB_SECTION_ID, qa.SubSectionId)
	assert.Equal(t, SUB_SECTION_ID, qa.SubSection.Id)
	assert.Equal(t, "B-Tree", qa.SubSection.Title)

	assert.Equal(t, "Search (Worst)", qa.Text.Text)
	assert.Equal(t, false, qa.Text.IsHtml)

	assert.Equal(t, "O(log(n))", qa.Answer.Text)

	assert.NotEmpty(t, qa.Choices)
}

func TestLoadQuizWithReverseSection(t *testing.T) {
	q := loadQuiz(t, "datastructures")

	if q.Sections == nil {
		t.Error("The quiz has no sections.")
	}

	const SECTION_ID = "reverse-datastructures-hash-tables"
	section := q.GetSection(SECTION_ID)
	if section == nil {
		t.Error("The quiz does not have the expected reverse section.")
	}

	assert.Equal(t, SECTION_ID, section.Id)
	assert.Equal(t, "Reverse: Hash Tables", section.Title)
	assert.Equal(t, "https://en.wikipedia.org/wiki/Hash_table", section.Link)

	const QUESTION_ID = "reverse-datastructures-hash-tables-open-addressing-strategy-probe-sequence"
	qa := q.GetQuestionAndAnswer(QUESTION_ID)
	if qa == nil {
		t.Error("The quiz does not have the expected reverse question.")
	}

	assert.Equal(t, QUESTION_ID, qa.Id)
	assert.Equal(t, SECTION_ID, qa.SectionId)
	assert.Equal(t, SECTION_ID, qa.Section.Id)
	assert.Equal(t, "Reverse: Hash Tables", qa.Section.Title)

	const SUB_SECTION_ID = "datastructures-hash-tables-open-addressing-strategies"
	assert.Equal(t, SUB_SECTION_ID, qa.SubSectionId)
	assert.Equal(t, SUB_SECTION_ID, qa.SubSection.Id)

	assert.Equal(t, "Open addressing strategies", qa.SubSection.Title)
	assert.Equal(t, "The buckets are examined, starting with the hashed-to slot and proceeding in some probe sequence, until an unoccupied slot is found.", qa.Text.Text)
	assert.Equal(t, false, qa.Text.IsHtml)

	assert.Equal(t, "Probe sequence", qa.Answer.Text)

	assert.NotEmpty(t, qa.Choices)
}
