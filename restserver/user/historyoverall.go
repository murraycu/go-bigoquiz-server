package user

import (
	"github.com/murraycu/go-bigoquiz-server/domain/user"
)

type HistoryOverall struct {
	LoginInfo *LoginInfo `json:"loginInfo,omitempty"`

	// Stats per quiz
	Stats []*user.Stats `json:"stats,omitempty"`
}

func (self *HistoryOverall) SetQuizStats(quizId string, stats *user.Stats) {
	self.Stats = append(self.Stats, stats)
}
