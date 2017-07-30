package bigoquiz

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"quiz"
	"strconv"
	"strings"
)

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

func loadQuiz(id string) (*quiz.Quiz, error) {
	absFilePath, err := filepath.Abs("quizzes/" + id + ".xml")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return quiz.LoadQuiz(absFilePath, id)
}

func loadQuizzes() ([]*quiz.Quiz, error) {
	quizzes := make([]*quiz.Quiz, 0)

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

// TODO: Is there instead some way to output just the top-level of the JSON,
// and only some of the fields?
// TODO: Avoid outputing an empty section list in the JSON.
func buildQuizzesSimple(quizzes []*quiz.Quiz) []*quiz.Quiz {
	// Create a slice with the same capacity.
	result := make([]*quiz.Quiz, 0, len(quizzes))

	for _, q := range quizzes {
		var simple quiz.Quiz
		simple.Id = q.Id
		simple.Title = q.Title
		simple.IsPrivate = q.IsPrivate

		result = append(result, &simple)
	}

	return result
}

func restHandleQuiz(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	listOnly := false
	queryValues := r.URL.Query()
	if queryValues != nil {
		listOnlyStr := queryValues.Get("list_only")
		listOnly, _ = strconv.ParseBool(listOnlyStr)
	}

	// TODO: Cache this.
	quizzes, err := loadQuizzes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if listOnly {
		// TODO: Cache this.
		quizzes = buildQuizzesSimple(quizzes)
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
