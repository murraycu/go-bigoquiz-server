package quiz

import (
	"path/filepath"
	"testing"
)

func TestLoadQuizWithBadFilepath(t *testing.T) {
	id := "doesnotexist"
	absFilePath, err := filepath.Abs("../bigoquiz/quizzes/" + id + ".xml")
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

func TestLoadQuiz(t *testing.T) {
	id := "bigo"
	absFilePath, err := filepath.Abs("../bigoquiz/quizzes/" + id + ".xml")
	if err != nil {
		t.Error("Could not find file.", err)
	}

	q, err := LoadQuiz(absFilePath, id)
	if err != nil {
		t.Error("LoadQuiz() failed.", err)
	}

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

	if len(qa.Question.Choices) == 0 {
		t.Error("The question does not have any choices")
	}
}
