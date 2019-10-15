package quiz

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Quiz struct {
	HasIdAndTitle
	IsPrivate        bool `xml:"is_private,attr"`
	AnswersAsChoices bool `xml:"answers_as_choices,attr"`

	Sections  []*Section           `xml:"section"`
	Questions []*QuestionAndAnswer `xml:"question"`

	UsesMathML bool `xml:"uses_mathml,attr"`
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

	for _, s := range self.Sections {
		//Make sure that we set sub-section choices from the answers from all questions in the whole section:
		if s.AnswersAsChoices {
			questionsIncludingSubSections := make([]*QuestionAndAnswer, 0)
			questionsIncludingSubSections = append(questionsIncludingSubSections, s.Questions...)

			for _, sub := range s.SubSections {
				questionsIncludingSubSections = append(questionsIncludingSubSections, sub.Questions...)
			}
		}
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
