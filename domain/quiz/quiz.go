package quiz

import (
	"math/rand"
)

type Quiz struct {
	HasIdAndTitle
	IsPrivate        bool `json:"isPrivate"`
	AnswersAsChoices bool `json:"answersAsChoices"`

	Sections  []*Section           `json:"sections,omitempty"`
	Questions []*QuestionAndAnswer `json:"questions,omitempty"`

	UsesMathML bool `json:"usesMathML"`

	// A map of all questions in all sections and at the top-level.
	questionsMap map[string]*QuestionAndAnswer `json:"-"`

	// An array of all questions in all sections and at the top-level.
	questionsArray []*QuestionAndAnswer `json:"-"`

	// A map of all sections by section ID.
	sectionsMap map[string]*Section `json:"-"`
}

/** Add to the map and array.
- * And also make sure that questions have their section ID and sub-section IDs.
- */
func (self *Quiz) addQuestionToMapAndArray(qa *QuestionAndAnswer, section *Section, subSection *SubSection) {
	self.setQuestionDetails(qa, section, subSection)

	self.questionsMap[qa.Id] = qa
	self.questionsArray = append(self.questionsArray, qa)
}

func (self *Quiz) InitMaps() error {
	self.questionsMap = make(map[string]*QuestionAndAnswer)
	for _, qa := range self.Questions {
		self.addQuestionToMapAndArray(qa, nil, nil)
	}

	self.sectionsMap = make(map[string]*Section)
	for _, s := range self.Sections {
		self.sectionsMap[s.Id] = s

		for _, qa := range s.Questions {
			self.addQuestionToMapAndArray(qa, s, nil)

			for _, ss := range s.SubSections {
				for _, qa := range s.Questions {
					self.addQuestionToMapAndArray(qa, s, ss)
				}
			}
		}
	}

	return nil
}

func (self *Quiz) GetAnswer(questionId string) *Text {
	if self.questionsMap == nil {
		return nil
	}

	qa, ok := self.questionsMap[questionId]
	if !ok {
		return nil
	}

	if qa == nil {
		return nil
	}

	return &qa.Answer
}

func (self *Quiz) GetSection(sectionId string) *Section {
	if self.sectionsMap == nil {
		return nil
	}

	s, ok := self.sectionsMap[sectionId]
	if !ok {
		return nil
	}

	if s == nil {
		return nil
	}

	return s
}

/** Make sure that questions have their section ID, sub-section IDs, and choices.
 */
func (self *Quiz) setQuestionDetails(qa *QuestionAndAnswer, section *Section, subSection *SubSection) {
	if qa == nil {
		return
	}

	q := &(qa.Question)

	// Update the section and subSection,
	// so we can return it in the JSON,
	// so the caller of the REST API does not need to discover these details separately.
	if section != nil {
		q.SectionId = section.Id
		q.Section = &(section.HasIdAndTitle)
	}

	if subSection != nil {
		q.SubSectionId = subSection.Id
		q.SubSection = &(subSection.HasIdAndTitle)
	}
}

func (self *Quiz) GetQuestionsCount() int {
	return len(self.questionsArray)
}

func (self *Quiz) GetQuestionAndAnswer(questionId string) *QuestionAndAnswer {
	if self.questionsMap == nil {
		return nil
	}

	return self.questionsMap[questionId]
}

func getRandomQuestionFromSlice(questions []*QuestionAndAnswer) *Question {
	count := len(questions)
	if count == 0 {
		return nil
	}

	i := rand.Intn(count - 1)
	qa := questions[i]
	return &(qa.Question)
}

func (self *Quiz) GetRandomQuestion(sectionId string) *Question {
	if self.questionsMap == nil {
		return nil
	}

	if len(sectionId) == 0 {
		return getRandomQuestionFromSlice(self.questionsArray)
	} else {
		section, ok := self.sectionsMap[sectionId]
		if !ok {
			return nil
		}

		return getRandomQuestionFromSlice(section.QuestionsArray)
	}
}
