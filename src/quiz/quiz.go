package quiz

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Quiz struct {
	Id        string `json:"id" xml:"id"`
	Title     string `json:"title" xml:"title"`
	IsPrivate bool   `json:"isPrivate" xml:"isPrivate"`
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
