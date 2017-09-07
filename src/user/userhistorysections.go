package user

import (
	"quiz"
)

type HistorySections struct {
	LoginInfo LoginInfo `json:"loginInfo,omitempty"`
	QuizTitle string    `json:"quizTitle,omitempty"`

	Sections []*quiz.Section `json:"sections,omitempty"`
	Stats    []*Stats    `json:"stats,omitempty"`
}
