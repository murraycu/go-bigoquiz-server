package quiz

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Quiz struct {
	HasIdAndTitle
	IsPrivate bool `json:"isPrivate" xml:"isPrivate"`

	Sections  []*Section           `json:"sections,omitempty" xml:"section"`
	Questions []*QuestionAndAnswer `json:"questions,omitempty" xml:"question"`

	// A map of all questions in all sections and at the top-level.
	questionsMap map[string]*QuestionAndAnswer `json:"-" xml:"-"`
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

	q.buildQuestionsMap()

	return &q, nil
}

/** Build a map of all questions in all sections and at the top-level.
 */
func (self *Quiz) buildQuestionsMap() {
	self.questionsMap = make(map[string]*QuestionAndAnswer)
	for _, q := range self.Questions {
		self.questionsMap[q.Id] = q
	}

	for _, s := range self.Sections {
		for _, q := range s.Questions {
			self.questionsMap[q.Id] = q
		}
	}
}

func (self *Quiz) GetQuestionAndAnswer(questionId string) *QuestionAndAnswer {
	if self.questionsMap == nil {
		return nil
	}

	return self.questionsMap[questionId]
}
