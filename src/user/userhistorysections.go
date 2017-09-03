package user

type UserHistorySections struct {
	LoginInfo LoginInfo `json:"loginInfo,omitempty"`
	QuizTitle string    `json:"quizTitle,omitempty"`

	SectionStats []*UserStats `json:"userStats,omitempty"`
}
