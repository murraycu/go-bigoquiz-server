package quiz

type Section struct {
	HasIdAndTitle
	Questions   []*QuestionAndAnswer `json:"questions,omitempty" xml:"question"`
	SubSections []*SubSection        `json:"subSections,omitempty" xml:"subsection"`

	DefaultChoices   []*Text `json:"defaultChoices,omitempty" xml:"default_choices"`
	AnswersAsChoices bool    `json:"answersAsChoices" xml:"answers_as_choices,attr"`

	// Whether the quiz should contain an extra generated section,
	// with the answers as questions, and the questions as the answers.
	AndReverse bool `json:"-" xml:"and_reverse,attr"`

	// These do not appear in the JSON.
	subSectionsMap map[string]*SubSection `json:"-" xml:"-"`
	CountQuestions int                    `json:"-" xml:"-"`

	// An array of all questions in all sub-sections and in this section directly.
	QuestionsArray []*QuestionAndAnswer `json:"-" xml:"-"`
}

func (self *Section) GetSubSection(subSectionId string) *SubSection {
	if self.subSectionsMap == nil {
		return nil
	}

	s, ok := self.subSectionsMap[subSectionId]
	if !ok {
		return nil
	}

	if s == nil {
		return nil
	}

	return s
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
