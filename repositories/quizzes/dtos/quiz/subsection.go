package quiz

type SubSection struct {
	HasIdAndTitle
	Questions        []*QuestionAndAnswer `xml:"question" json:"question"`
	AnswersAsChoices bool                 `xml:"answers_as_choices,attr" json:"answersAsChoices,omitempty"`
}
