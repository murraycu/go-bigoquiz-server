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
	return &q, nil
}
