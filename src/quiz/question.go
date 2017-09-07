package quiz

type Question struct {
	HasIdAndTitle
	Text Text `json:"text,omitempty" xml:"text"`

	// These are not in the XML.
	SectionId    string `json:"sectionId,omitempty"`
	SubSectionId string `json:"subSectionId,omitempty"`

	QuizUsesMathML bool `json:"quizUsesMathML"`

	// These are not in the XML.
	// But we want to show them in the JSON.
	// We don't use the Section and SubSection types here,
	// because we don't want to recurse infinitely into questions and
	// their sections and then their questions again.
	Section    *HasIdAndTitle `json:"section,omitempty"`
	SubSection *HasIdAndTitle `json:"subSection,omitempty"`

	// These are not in the XML.
	// But we want to show them in the JSON.
	Choices []*Text `json:"choices,omitempty"`
}
