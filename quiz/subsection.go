package quiz

type SubSection struct {
	HasIdAndTitle
	Questions        []*QuestionAndAnswer `json:"questions,omitempty" xml:"question"`
	AnswersAsChoices bool                 `json:"answersAsChoices" xml:"answers_as_choices,attr"`
}
