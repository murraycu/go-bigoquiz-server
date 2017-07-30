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

func loadQuizzes() (map[string]*quiz.Quiz, error) {
	quizzes := make(map[string]*quiz.Quiz, 0)

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

		quizzes[q.Id] = q
	}

	return quizzes, nil
}

// TODO: Is there instead some way to output just the top-level of the JSON,
// and only some of the fields?
func buildQuizzesSimple(quizzes map[string]*quiz.Quiz) map[string]*quiz.Quiz {
	// Create a slice with the same capacity.
	result := make(map[string]*quiz.Quiz)

	for _, q := range quizzes {
		var simple quiz.Quiz
		simple.Id = q.Id
		simple.Title = q.Title
		simple.IsPrivate = q.IsPrivate

		result[simple.Id] = &simple
	}

	return result
}

func restHandleQuizAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func restHandleQuizById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	quizId := ps.ByName("quizId")
	if quizId == "" {
		// This makes no sense. restHandleQuizAll() should have been called.
		http.Error(w, "Empty quiz ID", http.StatusInternalServerError)
		return
	}

	// TODO: Cache this.
	quizzes, err := loadQuizzes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := quizzes[quizId]
	if q == nil {
		http.Error(w, "quiz not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // normal header
	w.WriteHeader(http.StatusOK)

	jsonStr, err := json.Marshal(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonStr)
}
