package user

import (
	"github.com/murraycu/go-bigoquiz-server/domain/user"
)

type HistorySections struct {
	LoginInfo LoginInfo `json:"loginInfo,omitempty"`
	QuizId    string    `json:"quizId,omitempty"`
	QuizTitle string    `json:"quizTitle,omitempty"`

	Stats []*user.Stats `json:"stats,omitempty"`
}
