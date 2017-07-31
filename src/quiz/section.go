package quiz

type Section struct {
	HasIdAndTitle
	Questions   []QuestionAndAnswer `json:"questions,omitempty" xml:"question"`
	SubSections []SubSection        `json:"subsections,omitempty" xml:"subsection"`
}
