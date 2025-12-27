package quiz

type QuestionAndAnswer struct {
	Question `json:"question,omitempty"`
	Answer   Text `json:"answer,omitempty"`
}
