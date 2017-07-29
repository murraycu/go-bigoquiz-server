package bigoquiz

import (
	"encoding/json"
	"encoding/xml"
	"path/filepath"
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
)

type quiz struct {
  Id        string `json:"id" xml:"id"`
  Title     string `json:"title" xml:"title"`
  IsPrivate bool   `json:"isPrivate" xml:"isPrivate"`
}

func init() {
	http.HandleFunc("/api/quiz", rest_quizzes)
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

func loadQuizzes() ([]*quiz, error) {
	quizzes := make([]*quiz, 0)

	quizNames := []string {"algorithms", "bigo"}
	for _, name := range quizNames {
	  q, err := loadQuiz(name)
	  if err != nil {
		  fmt.Println(err)
		  return quizzes, err
	  }

	  quizzes = append(quizzes, q)
	}

	return quizzes, nil
}

func rest_quizzes(w http.ResponseWriter, r *http.Request) {
	quizzes, err := loadQuizzes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // normal header
	w.WriteHeader(http.StatusOK)

	jsonStr, err := json.Marshal(quizzes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonStr)
}
