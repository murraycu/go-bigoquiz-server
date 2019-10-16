package quiz

type SubSection struct {
	HasIdAndTitle
	Questions []*QuestionAndAnswer `json:"questions,omitempty"`

	// TODO: We only need this until we have called setQuestionsChoicesFromAnswers().
	AnswersAsChoices bool `json:"-"`
}
