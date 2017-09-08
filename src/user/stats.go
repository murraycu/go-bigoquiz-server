package user

import (
	"google.golang.org/appengine/datastore"
	"quiz"
)

// Statistics for the user, either per quiz (SectionId is then empty), or per section in a quiz.
type Stats struct {
	Key *datastore.Key `json:"-"`

	UserId *datastore.Key `json:"-"`

	QuizId    string `json:"quizId,omitEmpty"`
	SectionId string `json:"sectionId,omitEmpty"`

	Answered int `json:"answered"`
	Correct  int `json:"correct"`

	CountQuestionsAnsweredOnce int `json:"countQuestionsAnsweredOnce"`
	CountQuestionsCorrectOnce  int `json:"countQuestionsCorrectOnce"`

	// These are from the quiz, for convenience
	// so they don't need to be in the database.
	// TODO: Make sure they are set for per-section stats.
	CountQuestions int `json:"countQuestionsAnsweredOnce" datastore:"-"`
	QuizTitle    string `json:"quizTitle,omitEmpty" datastore:"-"`
	SectionTitle string `json:"sectionTitle,omitEmpty" datastore:"-"`
}

func (self *Stats) IncrementAnswered() {
	self.Answered += 1
}

func (self *Stats) IncrementCorrect() {
	self.Correct += 1
}

func (self *Stats) UpdateProblemQuestion(question *quiz.Question, answerIsCorrect bool) {
	// TODO
}
