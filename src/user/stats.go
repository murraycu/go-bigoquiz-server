package user

import (
	"google.golang.org/appengine/datastore"
	"quiz"
)

// Statistics for the user, either per quiz (SectionId is then empty), or per section in a quiz.
type Stats struct {
	Key *datastore.Key `json:"-"`

	UserId *datastore.Key `json:"-" datastore:"userId"`

	QuizId    string `json:"quizId,omitEmpty" datastore:"quizId"`
	SectionId string `json:"sectionId,omitEmpty" datastore:"sectionId"`

	Answered int `json:"answered" datastore:"answered"`
	Correct  int `json:"correct" datastore:"correct"`

	CountQuestionsAnsweredOnce int `json:"countQuestionsAnsweredOnce" datastore:"countQuestionsAnsweredOnce"`
	CountQuestionsCorrectOnce  int `json:"countQuestionsCorrectOnce" datastore:"countQuestionsCorrectOnce"`

	// These are from the quiz, for convenience
	// so they don't need to be in the database.
	// TODO: Make sure they are set for per-section stats.
	CountQuestions int `json:"countQuestions" datastore:"-"`
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
