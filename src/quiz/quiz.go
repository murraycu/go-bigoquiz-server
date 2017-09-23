package quiz

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

const maxChoicesFromAnswers = 6

type Quiz struct {
	HasIdAndTitle
	IsPrivate        bool `json:"isPrivate" xml:"is_private,attr"`
	AnswersAsChoices bool `json:"answersAsChoices" xml:"answers_as_choices,attr"`

	Sections  []*Section           `json:"sections,omitempty" xml:"section"`
	Questions []*QuestionAndAnswer `json:"questions,omitempty" xml:"question"`

	UsesMathML bool `json:"usesMathML"`

	// A map of all questions in all sections and at the top-level.
	questionsMap map[string]*QuestionAndAnswer `json:"-" xml:"-"`

	// An array of all questions in all sections and at the top-level.
	questionsArray []*QuestionAndAnswer `json:"-" xml:"-"`

	// A map of all sections by section ID.
	sectionsMap map[string]*Section `json:"-" xml:"-"`
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

	// Deal with quizzes that have no sections, with just quizzes at the top-level:
	if len(q.Sections) == 0 {
		// Add a virtual section, so we have somewhere to put the questions.
		// This lets a quiz have just questions with no sections.
		// The generated section will have the same id and title as the quiz itself.
		var section Section
		section.Id = q.Id
		section.Title = q.Title
		section.Questions = q.Questions
		section.AnswersAsChoices = q.AnswersAsChoices
		q.Questions = nil

		q.Sections = append(q.Sections, &section)
	}

	q.process()

	return &q, nil
}

/** Various processing of the quiz after the simple unmarshalling from the XML.
 * For instance, build a map of all questions in all sections and at the top-level,
 * and a map of sections by section ID.
 * And also make sure that questions have their section ID and sub-section IDs,
 * and that the questions have the correct choices.
 */
func (self *Quiz) process() {
	self.addReverseSections()

	self.questionsMap = make(map[string]*QuestionAndAnswer)
	self.questionsArray = make([]*QuestionAndAnswer, 0, len(self.Questions))
	self.sectionsMap = make(map[string]*Section)

	// Build the map and array:
	for _, q := range self.Questions {
		self.questionsMap[q.Id] = q

		self.questionsArray = append(self.questionsArray, q)
	}

	for _, s := range self.Sections {
		self.sectionsMap[s.Id] = s

		s.subSectionsMap = make(map[string]*SubSection)

		var sectionCountQuestions int
		for _, sub := range s.SubSections {
			s.subSectionsMap[sub.Id] = sub

			for _, q := range sub.Questions {
				self.addQuestionToMapAndArray(q, s, sub)
			}

			sectionCountQuestions += len(sub.Questions)
			s.QuestionsArray = append(s.QuestionsArray, sub.Questions...)

			//Don't use subsection answers as choices if the parent section wants answers-as-choices.
			//In that case, all questions will instead share answers from all sub-sections.
			if sub.AnswersAsChoices && !s.AnswersAsChoices {
				setQuestionsChoicesFromAnswers(sub.Questions)
			}
		}

		//Add any Questions that are not in a subsection:
		for _, q := range s.Questions {
			self.addQuestionToMapAndArray(q, s, nil)
		}

		sectionCountQuestions += len(s.Questions)

		//Make sure that we set sub-section choices from the answers from all questions in the whole section:
		if s.AnswersAsChoices {
			questionsIncludingSubSections := make([]*QuestionAndAnswer, 0)
			questionsIncludingSubSections = append(questionsIncludingSubSections, s.Questions...)

			for _, sub := range s.SubSections {
				questionsIncludingSubSections = append(questionsIncludingSubSections, sub.Questions...)
			}

			setQuestionsChoicesFromAnswers(questionsIncludingSubSections)
		}

		s.CountQuestions = sectionCountQuestions
		s.QuestionsArray = append(s.QuestionsArray, s.Questions...)
	}

}

/** Optionally generate reverse sections.
 */
func (self *Quiz) addReverseSections() {
	reverseSections := make([]*Section, 0, 0)
	for _, s := range self.Sections {
		if s.AndReverse {
			reverseSections = append(reverseSections, s.createReverse())
		}
	}

	self.Sections = append(self.Sections, reverseSections...)
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
			q.Question.Choices = choices
		} else {
			reduced := reduceChoices(choices, &(q.Answer))
			q.Question.Choices = reduced
		}
	}
}

func (self *Quiz) GetQuestionsCount() int {
	return len(self.questionsArray)
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

	answerIndex, ok := getIndexInArray(result, answer)
	if !ok {
		return nil
	}

	if answerIndex >= maxChoicesFromAnswers {
		result = result[0 : maxChoicesFromAnswers-1]
		result = append(result, answer)
		shuffle(result)
	} else {
		result = result[0:maxChoicesFromAnswers]
	}

	return result
}

/** Get the index of an item in the array by comparig only the strings in the Text struct.
 */
func getIndexInArray(array []*Text, str *Text) (int, bool) {
	for i, s := range array {
		if s.Text == str.Text {
			return i, true
		}
	}

	return -1, false
}

func shuffle(array []*Text) {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	for i := len(array) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
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
	var qa *QuestionAndAnswer = questions[i]
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
