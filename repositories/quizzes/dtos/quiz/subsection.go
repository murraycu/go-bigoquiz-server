package quiz

type SubSection struct {
	HasIdAndTitle
	Questions        []*QuestionAndAnswer `json:"question"`
	AnswersAsChoices bool                 `json:"answersAsChoices,omitempty"`
}
