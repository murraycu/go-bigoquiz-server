package quiz

type SubSection struct {
	HasIdAndTitle
	Questions        []*QuestionAndAnswer `json:"questions,omitempty"`
	AnswersAsChoices bool                 `json:"answersAsChoices"`
}
