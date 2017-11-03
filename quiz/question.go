package quiz

type Question struct {
	Id    string `json:"id" xml:"id,attr"`

	// A URL.
	Link string `json:"link,omitempty" xml:"link"`

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
	QuizTitle  string         `json:"quizTitle"`
	Section    *HasIdAndTitle `json:"section,omitempty"`
	SubSection *HasIdAndTitle `json:"subSection,omitempty"`

	// These are not in the XML.
	// But we want to show them in the JSON.
	Choices []*Text `json:"choices,omitempty"`
}

/** Set extra titles for convenience.
 */
func (self *Question) SetTitles(quizTitle string, briefSection *HasIdAndTitle, subSection *SubSection) {
	self.QuizTitle = quizTitle
	self.Section = briefSection

	if subSection != nil {
		self.SubSection = &(subSection.HasIdAndTitle)
	}
}

/** Set extra details for the question,
 * such as quiz and section titles,
 * so the client doesn't need to look these up separately.
 */
func (self *Question) SetQuestionExtras(q *Quiz) {
	var briefSection *HasIdAndTitle
	var subSection *SubSection

	section := q.GetSection(self.SectionId)
	if section != nil {
		// Create a simpler version of the Section information:
		briefSection = new(HasIdAndTitle)
		briefSection.Id = section.Id
		briefSection.Title = section.Title
		briefSection.Link = section.Link

		subSection = section.GetSubSection(self.SubSectionId)
	}

	self.SetTitles(q.Title, briefSection, subSection)

	self.QuizUsesMathML = q.UsesMathML
}
