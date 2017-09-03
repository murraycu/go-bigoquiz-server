package quiz

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
)

const maxChoicesFromAnswers = 6

type Quiz struct {
	HasIdAndTitle
	IsPrivate bool `json:"isPrivate" xml:"isPrivate"`

	Sections  []*Section           `json:"sections,omitempty" xml:"section"`
	Questions []*QuestionAndAnswer `json:"questions,omitempty" xml:"question"`

	// A map of all questions in all sections and at the top-level.
	questionsMap map[string]*QuestionAndAnswer `json:"-" xml:"-"`

	// An array of all questions in all sections and at the top-level.
	questionsArray []*QuestionAndAnswer `json:"-" xml:"-"`
}

func LoadQuiz(absFilePath string, id string) (*Quiz, error) {
	var q Quiz

	file, err := os.Open(absFilePath)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = xml.Unmarshal(data, &q)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	q.Id = id

	q.buildQuestionsMapAndArray()

	return &q, nil
}

/** Build a map of all questions in all sections and at the top-level.
 * And also make sure that questions have their section ID and sub-section IDs,
 * and that the questions have the correct choices.
 */
func (self *Quiz) buildQuestionsMapAndArray() {
	self.questionsMap = make(map[string]*QuestionAndAnswer)
	self.questionsArray = make([]*QuestionAndAnswer, 0, len(self.Questions))

	// Build the map and array:
	for _, q := range self.Questions {
		self.questionsMap[q.Id] = q

		self.questionsArray = append(self.questionsArray, q)
	}

	for _, s := range self.Sections {
		for _, q := range s.Questions {
			self.addQuestionToMapAndArray(q, s, nil)
		}

		if s.AnswersAsChoices {
			setQuestionsChoicesFromAnswers(s.Questions)
		}

		for _, sub := range s.SubSections {
			for _, q := range sub.Questions {
				self.addQuestionToMapAndArray(q, s, sub)
			}

			if sub.AnswersAsChoices {
				setQuestionsChoicesFromAnswers(sub.Questions)
			}
		}
	}

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

/** Add to the map and array.
 * And also make sure that questions have their section ID and sub-section IDs.
 */
func (self *Quiz) addQuestionToMapAndArray(qa *QuestionAndAnswer, section *Section, subSection *SubSection) {
	self.setQuestionDetails(qa, section, subSection)

	self.questionsMap[qa.Id] = qa
	self.questionsArray = append(self.questionsArray, qa)
}

func setQuestionsChoicesFromAnswers(questions []*QuestionAndAnswer) {
	// Build the list of answers, avoiding duplicates:
	choices := make([]*Text, 0, len(questions))
	used := make(map[string]bool)

	for _, q := range questions {
		t := q.Answer.Text
		_, ok := used[t]
		if ok {
			continue
		}

		used[t] = true
		choices = append(choices, &(q.Answer))
	}

	tooManyChoices := len(choices) > maxChoicesFromAnswers

	for _, q := range questions {
		if !tooManyChoices {
			q.Choices = choices
		} else {
			reduced := reduceChoices(choices, &(q.Answer))
			q.Choices = reduced
		}
	}
}

/**
 * Create a small-enough set of choices which
 * always contains the correct answer.
 * This is slow.
 */
func reduceChoices(choices []*Text, answer *Text) []*Text {

	result := make([]*Text, len(choices))
	copy(result, choices)
	shuffle(result)

	answerIndex, ok := getIndexInArray(choices, answer)
	if !ok {
		return nil
	}

	if answerIndex >= maxChoicesFromAnswers {
		result = choices[0 : maxChoicesFromAnswers-1]
		result = append(result, answer)
		shuffle(result)
	} else {
		result = choices[0:maxChoicesFromAnswers]
	}

	return result
}

func getIndexInArray(array []*Text, str *Text) (int, bool) {
	for i, s := range array {
		if s == str {
			return i, true
		}
	}

	return -1, false
}

func shuffle(array []*Text) {

	for i := len(array) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
}

func (self *Quiz) GetQuestionAndAnswer(questionId string) *QuestionAndAnswer {
	if self.questionsMap == nil {
		return nil
	}

	return self.questionsMap[questionId]
}

func (self *Quiz) GetRandomQuestion() *Question {
	if self.questionsMap == nil {
		return nil
	}

	count := len(self.questionsArray)
	i := rand.Intn(count - 1)
	var qa *QuestionAndAnswer = self.questionsArray[i]
	return &(qa.Question)
}
