package quiz

type Section struct {
	HasIdAndTitle
	Questions   []*QuestionAndAnswer `json:"questions,omitempty"`
	SubSections []*SubSection        `json:"subSections,omitempty"`

	DefaultChoices   []*Text `json:"defaultChoices,omitempty"`
	AnswersAsChoices bool    `json:"answersAsChoices"`

	// These do not appear in the JSON.
	subSectionsMap map[string]*SubSection `json:"-"`
	CountQuestions int                    `json:"-"`

	// An array of all questions in all sub-sections and in this section directly.
	QuestionsArray []*QuestionAndAnswer `json:"-"`
}

func (self *Section) InitMaps() error {
	self.subSectionsMap = make(map[string]*SubSection)
	for _, s := range self.SubSections {
		self.subSectionsMap[s.HasIdAndTitle.Id] = s

		for _, qa := range s.Questions {
			self.addQuestionArray(qa)
		}
	}

	for _, qa := range self.Questions {
		self.addQuestionArray(qa)
	}

	return nil
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

func (self *Section) addQuestionArray(qa *QuestionAndAnswer) {
	self.QuestionsArray = append(self.QuestionsArray, qa)
	self.CountQuestions++
}
