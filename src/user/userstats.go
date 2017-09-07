package user

import "google.golang.org/appengine/datastore"

type Stats struct {
	UserId    datastore.Key `json:"-"`

	QuizId    string `json:"quizId,omitEmpty"`
	SectionId string `json:"sectionId,omitEmpty"`

	Answered int `json:"answered"`
	Correct  int `json:"correct"`

	CountQuestionsAnsweredOnce int `json:"countQuestionsAnsweredOnce"`
	CountQuestionsCorrectOnce  int `json:"countQuetionsCorrectOnce"`
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
