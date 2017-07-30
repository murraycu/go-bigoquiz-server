package quiz

type Section struct {
	HasIdAndTitle
	Questions   []Question   `json:"questions,omitempty" xml:"question"`
	SubSections []SubSection `json:"subsections,omitempty" xml:"subsection"`
}
