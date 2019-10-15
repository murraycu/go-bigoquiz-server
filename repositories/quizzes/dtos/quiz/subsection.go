package quiz

type SubSection struct {
	HasIdAndTitle
	Questions        []*QuestionAndAnswer `xml:"question"`
	AnswersAsChoices bool                 `xml:"answers_as_choices,attr"`
}
