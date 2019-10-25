package quiz

type Section struct {
	HasIdAndTitle
	Questions   []*QuestionAndAnswer `xml:"question"`
	SubSections []*SubSection        `xml:"subsection"`

	DefaultChoices   []*Text `xml:"default_choices"`
	AnswersAsChoices bool    `xml:"answers_as_choices,attr"`

	// Whether the quiz should contain an extra generated section,
	// with the answers as questions, and the questions as the answers.
	AndReverse bool `xml:"and_reverse,attr"`
}

func (self *Section) createReverse() *Section {
	var result Section

	result.Id = "reverse-" + self.Id
	result.Title = "Reverse: " + self.Title
	result.Link = self.Link
	result.AnswersAsChoices = self.AnswersAsChoices

	for _, sub := range self.SubSections {
		var reverseSub SubSection
		reverseSub.Id = sub.Id
		reverseSub.Title = sub.Title
		reverseSub.Link = sub.Link
		reverseSub.AnswersAsChoices = sub.AnswersAsChoices

		for _, q := range sub.Questions {
			reverseSub.Questions = append(reverseSub.Questions, q.createReverse())
		}

		result.SubSections = append(result.SubSections, &reverseSub)
	}

	for _, q := range self.Questions {
		result.Questions = append(result.Questions, q.createReverse())
	}

	return &result
}
