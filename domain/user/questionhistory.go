package user

import "github.com/murraycu/go-bigoquiz-server/domain/quiz"

type QuestionHistory struct {
	QuestionId string `json:"questionId,omitempty"`

	AnsweredCorrectlyOnce bool `json:"answeredCorrectlyOnce"`

	//Decrements once for each time the user answers it correctly.
	//Increments once for each time the user answers it wrongly.
	CountAnsweredWrong int `json:"countAnsweredWrong"`

	// These are in the JSON for the convenience of the caller,
	// but they should not be in the datastore:
	QuestionTitle *quiz.Text `json:"questionTitle,omitempty"`
	SectionId     string     `json:"sectionId,omitempty"`
	// The caller doesn't need the SectionTitle because these are already stored within the stats for a particular section.
	SubSectionTitle string `json:"subSectionTitle,omitempty"`
}