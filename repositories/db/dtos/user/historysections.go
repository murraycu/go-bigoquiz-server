package user

type HistorySections struct {
	LoginInfo LoginInfo `json:"loginInfo,omitempty"`
	QuizId    string    `json:"quizId,omitempty"`
	QuizTitle string    `json:"quizTitle,omitempty"`

	Stats []*Stats `json:"stats,omitempty"`
}
