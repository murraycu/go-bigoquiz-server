package user

import (
	"quiz"
)

type UserHistorySections struct {
	LoginInfo LoginInfo `json:"loginInfo,omitempty"`
	QuizTitle string    `json:"quizTitle,omitempty"`

	Sections []*quiz.Section `json:"sections,omitempty"`
	Stats    []*UserStats    `json:"stats,omitempty"`
}
