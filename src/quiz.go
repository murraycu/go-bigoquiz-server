package bigoquiz

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type quiz struct {
	Id        string `json:"id" xml:"id"`
	Title     string `json:"title" xml:"title"`
	IsPrivate bool   `json:"isPrivate" xml:"isPrivate"`
}

func loadQuiz(id string) (*quiz, error) {
	var q quiz

	absFilePath, err := filepath.Abs("quizzes/" + id + ".xml")
	if err != nil {
		fmt.Println(err)
		return &q, err
	}

	file, err := os.Open(absFilePath)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return &q, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return &q, err
	}

	err = xml.Unmarshal(data, &q)
	if err != nil {
		fmt.Println(err)
		return &q, err
	}

	q.Id = id
	return &q, nil
}
