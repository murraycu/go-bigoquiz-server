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
