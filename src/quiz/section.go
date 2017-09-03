package quiz

type Section struct {
	HasIdAndTitle
	Questions   []*QuestionAndAnswer `json:"questions,omitempty" xml:"question"`
	SubSections []*SubSection        `json:"subSections,omitempty" xml:"subsection"`
}
