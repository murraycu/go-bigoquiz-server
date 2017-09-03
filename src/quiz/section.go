package quiz

type Section struct {
	HasIdAndTitle
	Questions   []*QuestionAndAnswer `json:"questions,omitempty" xml:"question"`
	SubSections []*SubSection        `json:"subSections,omitempty" xml:"subsection"`

	DefaultChoices []*Text           `json:"defaultChoices,omitempty" xml:"default_choices"`
	AnswersAsChoices bool            `json:"answersAsChoices" xml:"answers_as_choices,attr"`

}
