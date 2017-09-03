package quiz

type Question struct {
	HasIdAndTitle
	Text Text `json:"text,omitempty" xml:"text"`

	// These are not in the XML.
	SectionId    string `json:"sectionId,omitempty"`
	SubSectionId string `json:"subSectionId,omitempty"`
}
