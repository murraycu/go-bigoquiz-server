package quiz

type Quiz struct {
	HasIdAndTitle
	IsPrivate bool

	Sections  []*Section
	Questions []*QuestionAndAnswer

	UsesMathML bool

	// TODO: We only need this until we have called setQuestionsChoicesFromAnswers().
	AnswersAsChoices bool
}
