package user

type HistorySections struct {
	LoginInfo LoginInfo `json:"loginInfo,omitempty"`
	QuizTitle string    `json:"quizTitle,omitempty"`

	Stats    []*Stats    `json:"stats,omitempty"`
}
