package quiz

type QuestionAndAnswer struct {
	HasIdAndTitle
	Text Text `json:"text,omitempty" xml:"text"`
	Answer Text `json:"answer,omitempty" xml:"answer"`

	// These are not in the XML.
	SectionId    Text `json:"sectionId,omitempty"`
	SubSectionId Text `json:"subSectionId,omitempty"`
}
