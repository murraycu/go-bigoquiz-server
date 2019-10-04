package quiz

import (
	"path/filepath"
	"testing"
)

func TestLoadQuizWithBadFilepath(t *testing.T) {
	id := "doesnotexist"
	absFilePath, err := filepath.Abs("../quizzes/" + id + ".xml")
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
	absFilePath, err := filepath.Abs("../quizzes/" + quizId + ".xml")
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

	if section.Id != SECTION_ID {
		t.Error("The section does not have the expected ID.")
	}

	if section.Title != "Data Structure Operations" {
		t.Error("The section does not have the expected title.")
	}

	if section.Link != "" {
		t.Error("The section does not have the expected link.")
	}

	const QUESTION_ID = "b-tree-search-worst"
	qa := q.GetQuestionAndAnswer(QUESTION_ID)
	if qa == nil {
		t.Error("The quiz does not have the expected question.")
	}

	if qa.Question.Id != QUESTION_ID {
		t.Error("The question does not have the expected ID.")
	}

	if qa.Question.SectionId != SECTION_ID {
		t.Error("The question does not have the expected SectionId.")
	}

	if qa.Question.Section.Id != SECTION_ID {
		t.Error("The question does not have the expected section ID.")
	}

	if qa.Question.Section.Title != "Data Structure Operations" {
		t.Error("The question does not have the expected section ID.")
	}

	const SUB_SECTION_ID = "b-tree"
	if qa.Question.SubSectionId != SUB_SECTION_ID {
		t.Error("The question does not have the expected sub-section ID.")
	}

	if qa.Question.SubSection.Id != SUB_SECTION_ID {
		t.Error("The question does not have the expected sub-section ID.")
	}

	if qa.Question.SubSection.Title != "B-Tree" {
		t.Error("The question does not have the expected sub-section Title.")
	}

	if qa.Question.Text.Text != "Search (Worst)" {
		t.Error("The question does not have the expected text.")
	}

	if qa.Question.Text.IsHtml {
		t.Error("The question does not have the expected isHtml value.")
	}

	if qa.Answer.Text != "O(log(n))" {
		t.Error("The question does not have the expected answer text.")
	}

	if len(qa.Question.Choices) == 0 {
		t.Error("The question does not have any choices")
	}
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

	if section.Id != SECTION_ID {
		t.Error("The reverse section does not have the expected ID.")
	}

	if section.Title != "Reverse: Hash Tables" {
		t.Error("The reverse section does not have the expected title.")
	}

	if section.Link != "https://en.wikipedia.org/wiki/Hash_table" {
		t.Error("The reverse section does not have the expected link.")
	}

	const QUESTION_ID = "reverse-datastructures-hash-tables-open-addressing-strategy-probe-sequence"
	qa := q.GetQuestionAndAnswer(QUESTION_ID)
	if qa == nil {
		t.Error("The quiz does not have the expected reverse question.")
	}

	if qa.Question.Id != QUESTION_ID {
		t.Error("The reverse question does not have the expected ID.")
	}

	if qa.Question.SectionId != SECTION_ID {
		t.Error("The reverse question does not have the expected SectionId.")
	}

	if qa.Question.Section.Id != SECTION_ID {
		t.Error("The reverse question does not have the expected section ID.")
	}

	if qa.Question.Section.Title != "Reverse: Hash Tables" {
		t.Error("The reverse question does not have the expected section ID.")
	}

	const SUB_SECTION_ID = "datastructures-hash-tables-open-addressing-strategies"
	if qa.Question.SubSectionId != SUB_SECTION_ID {
		t.Error("The question does not have the expected subSectionId.")
	}

	if qa.Question.SubSection.Id != SUB_SECTION_ID {
		t.Error("The question does not have the expected sub-section ID.")
	}

	if qa.Question.SubSection.Title != "Open addressing strategies" {
		t.Error("The question does not have the expected sub-section title.")
	}

	if qa.Question.Text.Text != "The buckets are examined, starting with the hashed-to slot and proceeding in some probe sequence, until an unoccupied slot is found." {
		t.Error("The question does not have the expected text.", qa.Question.Text.Text)
	}

	if qa.Question.Text.IsHtml {
		t.Error("The question does not have the expected isHtml value.")
	}

	if qa.Answer.Text != "Probe sequence" {
		t.Error("The question does not have the expected answer text.", qa.Answer.Text)
	}

	if len(qa.Question.Choices) == 0 {
		t.Error("The question does not have any choices")
	}
}
