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
	QuizTitle string `json:"quizTitle"`
	Section    *HasIdAndTitle `json:"section,omitempty"`
	SubSection *HasIdAndTitle `json:"subSection,omitempty"`

	// These are not in the XML.
	// But we want to show them in the JSON.
	Choices []*Text `json:"choices,omitempty"`
}

// Set extra titles for convenience.
func (self *Question) SetTitles(quizTitle string, briefSection *HasIdAndTitle, subSection *SubSection) {
	self.QuizTitle = quizTitle
	self.Section = briefSection

	if subSection != nil {
		self.SubSection = &(subSection.HasIdAndTitle)
	}
}
