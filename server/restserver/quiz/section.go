package quiz

type Section struct {
	HasIdAndTitle
	Questions   []*QuestionAndAnswer `json:"questions,omitempty"`
	SubSections []*SubSection        `json:"subSections,omitempty"`

	DefaultChoices []*Text `json:"defaultChoices,omitempty"`

	// TODO: We only need this until we have called setQuestionsChoicesFromAnswers().
	AnswersAsChoices bool `json:"-"`
}
