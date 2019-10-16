package quiz

type Section struct {
	HasIdAndTitle
	Questions   []*QuestionAndAnswer
	SubSections []*SubSection

	DefaultChoices []*Text

	// TODO: We only need this until we have called setQuestionsChoicesFromAnswers().
	AnswersAsChoices bool
}
