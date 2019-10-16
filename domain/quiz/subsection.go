package quiz

type SubSection struct {
	HasIdAndTitle
	Questions        []*QuestionAndAnswer
	AnswersAsChoices bool
}
