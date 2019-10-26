package user

type HistoryOverall struct {
	LoginInfo *LoginInfo `json:"loginInfo,omitempty"`

	// Stats per quiz
	Stats []*Stats `json:"stats,omitempty"`
}

func (self *HistoryOverall) AddQuizStats(quizId string, stats *Stats) {
	self.Stats = append(self.Stats, stats)
}
