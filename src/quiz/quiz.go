package quiz

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
)

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
 */
func (self *Quiz) buildQuestionsMapAndArray() {
	self.questionsMap = make(map[string]*QuestionAndAnswer)
	self.questionsArray = make([]*QuestionAndAnswer, 0, len(self.Questions))

	for _, q := range self.Questions {
		self.questionsMap[q.Id] = q

		self.questionsArray = append(self.questionsArray, q);
	}

	for _, s := range self.Sections {
		for _, q := range s.Questions {
			self.questionsMap[q.Id] = q

			self.questionsArray = append(self.questionsArray, q);
		}

		for _, sub := range s.SubSections {
		  for _, q := range sub.Questions {
			  self.questionsMap[q.Id] = q

			  self.questionsArray = append(self.questionsArray, q);
		  }
		}
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
	return &(qa.Question);
}
