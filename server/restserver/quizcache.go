package restserver

import (
	"fmt"
	restquiz "github.com/murraycu/go-bigoquiz-server/server/restserver/quiz"
	"math/rand"
)

type restQuestionAndAnswerArray []*restquiz.QuestionAndAnswer

// Map of sub section IDs to SubSections.
type restSubSectionsMap map[string]*restquiz.SubSection

type QuizCache struct {
	Quiz *restquiz.Quiz

	questionsMap map[string]*restquiz.QuestionAndAnswer

	// An array of all questions in all sections and at the top-level.
	questionsArray restQuestionAndAnswerArray

	// A map of all sections by section ID.
	sectionsMap map[string]*restquiz.Section

	sectionsSubSectionsMap map[string]restSubSectionsMap

	// A map of all sections IDs to an array of questions in that section.
	sectionsQuestionsArrayMap map[string]restQuestionAndAnswerArray
}

func NewQuizCache(quiz *restquiz.Quiz) (*QuizCache, error) {
	result := &QuizCache{}
	result.Quiz = quiz

	err := result.initMaps(quiz)
	if err != nil {
		return nil, fmt.Errorf("initMaps() failed: %v", err)
	}

	return result, nil
}

// Add to the map and array.
func (self *QuizCache) addQuestionToMapAndArray(qa *restquiz.QuestionAndAnswer, section *restquiz.Section, subSection *restquiz.SubSection) error {
	_, exists := self.questionsMap[qa.Id]
	if exists {
		return fmt.Errorf("questionsMap already contains a question with this ID: %v", qa.Id)
	}

	self.questionsMap[qa.Id] = qa

	self.questionsArray = append(self.questionsArray, qa)

	if section != nil {
		sectionQuestionArray, _ := self.sectionsQuestionsArrayMap[section.Id]
		sectionQuestionArray = append(sectionQuestionArray, qa)

		self.sectionsQuestionsArrayMap[section.Id] = sectionQuestionArray
	}

	return nil
}

func (self *QuizCache) initMaps(quiz *restquiz.Quiz) error {
	self.questionsMap = make(map[string]*restquiz.QuestionAndAnswer)
	self.sectionsQuestionsArrayMap = make(map[string]restQuestionAndAnswerArray)
	self.sectionsMap = make(map[string]*restquiz.Section)
	self.sectionsSubSectionsMap = make(map[string]restSubSectionsMap)

	for _, s := range quiz.Sections {
		self.sectionsMap[s.Id] = s

		for _, qa := range s.Questions {
			err := self.addQuestionToMapAndArray(qa, s, nil)
			if err != nil {
				return fmt.Errorf("addQuestionToMapAndArray() failed: %v", err)
			}
		}

		subSectionsMap := make(restSubSectionsMap)
		for _, ss := range s.SubSections {
			subSectionsMap[ss.Id] = ss

			for _, qa := range ss.Questions {
				err := self.addQuestionToMapAndArray(qa, s, ss)
				if err != nil {
					return fmt.Errorf("addQuestionToMapAndArray() failed: %v", err)
				}
			}
		}

		self.sectionsSubSectionsMap[s.Id] = subSectionsMap
	}

	return nil
}

func getRandomQuestionFromSlice(questions []*restquiz.QuestionAndAnswer) *restquiz.Question {
	count := len(questions)
	if count == 0 {
		return nil
	}

	i := rand.Intn(count - 1)
	qa := questions[i]
	return &(qa.Question)
}

func (self *QuizCache) GetRandomQuestion(sectionId string) *restquiz.Question {
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

		questions := self.getQuestionsArrayForSection(section.Id)
		return getRandomQuestionFromSlice(questions)
	}
}

// GetQuestionsCount returns the number of questions in the entire quiz.
func (self *QuizCache) GetQuestionsCount() int {
	return len(self.questionsArray)
}

// GetSectionQuestionsCount returns the number of questions in the section and all its sub-sections.
func (self *QuizCache) GetSectionQuestionsCount(sectionId string) int {
	section, err := self.GetSection(sectionId)
	if err != nil || section == nil {
		return 0
	}

	result := len(section.Questions)

	for _, subSection := range section.SubSections {
		if subSection == nil {
			continue
		}

		result += len(subSection.Questions)
	}

	return result
}

func (self *QuizCache) GetAnswer(questionId string) *restquiz.Text {
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

func (self *QuizCache) GetQuestionAndAnswer(questionId string) *restquiz.QuestionAndAnswer {
	if self.questionsMap == nil {
		return nil
	}

	return self.questionsMap[questionId]
}

func (self *QuizCache) getQuestionsArrayForSection(sectionId string) restQuestionAndAnswerArray {
	if self.sectionsQuestionsArrayMap == nil {
		return nil
	}

	result, ok := self.sectionsQuestionsArrayMap[sectionId]
	if !ok {
		return nil
	}

	return result
}

// Returns nil if sectionId is nil (not an error)
func (self *QuizCache) GetSection(sectionID string) (*restquiz.Section, error) {
	if len(sectionID) == 0 {
		return nil, nil
	}

	section, ok := self.sectionsMap[sectionID]
	if !ok {
		return nil, fmt.Errorf("Section not found with section ID: %v", sectionID)
	}

	return section, nil
}

func (self *QuizCache) GetSubSection(sectionId string, subSectionId string) *restquiz.SubSection {
	m, ok := self.sectionsSubSectionsMap[sectionId]
	if !ok {
		return nil
	}

	subSection, ok := m[subSectionId]
	if !ok {
		return nil
	}

	return subSection
}
