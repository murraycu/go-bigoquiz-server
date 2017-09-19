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
		t.Error("LoadQuiz() did not faile with an error.")
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
}
