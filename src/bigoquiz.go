package bigoquiz

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func filesWithExtension(dirPath string, ext string) ([]string, error) {
	result := make([]string, 0)

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	dotSuffix := "." + ext
	suffixLen := len(dotSuffix)
	for _, f := range files {

		name := f.Name()
		if strings.HasSuffix(name, dotSuffix) {
			prefix := name[0 : len(name)-suffixLen]
			result = append(result, prefix)
		}
	}

	return result, nil
}

func loadQuizzes() ([]*quiz, error) {
	quizzes := make([]*quiz, 0)

	absFilePath, err := filepath.Abs("quizzes")
	if err != nil {
		fmt.Println(err)
		return quizzes, err
	}

	quizNames, err := filesWithExtension(absFilePath, "xml")
	if err != nil {
		fmt.Println(err)
		return quizzes, err
	}

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
