package bigoquiz

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"quiz"
	"sort"
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

// See https://gobyexample.com/sorting-by-functions
type QuizListByTitle []*quiz.Quiz

func (s QuizListByTitle) Len() int {
	return len(s)
}

func (s QuizListByTitle) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s QuizListByTitle) Less(i, j int) bool {
	return s[i].Title < s[j].Title
}

// TODO: Is there instead some way to output just the top-level of the JSON,
// and only some of the fields?
func buildQuizzesSimple(quizzes map[string]*quiz.Quiz) []*quiz.Quiz {
	// Create a slice with the same capacity.
	result := make([]*quiz.Quiz, 0, len(quizzes))

	for _, q := range quizzes {
		var simple quiz.Quiz
		q.CopyHasIdAndTitle(&simple.HasIdAndTitle)
		simple.IsPrivate = q.IsPrivate

		result = append(result, &simple)
	}

	sort.Sort(QuizListByTitle(result))

	return result
}

func buildQuizzesFull(quizzes map[string]*quiz.Quiz) []*quiz.Quiz {
	// Create a slice with the same capacity.
	result := make([]*quiz.Quiz, 0, len(quizzes))

	for _, q := range quizzes {
		result = append(result, q)
	}

	sort.Sort(QuizListByTitle(result))

	return result
}

func restHandleQuizAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	listOnly := false
	queryValues := r.URL.Query()
	if queryValues != nil {
		listOnlyStr := queryValues.Get(QUERY_PARAM_LIST_ONLY)
		listOnly, _ = strconv.ParseBool(listOnlyStr)
	}

	var quizArray []*quiz.Quiz = nil
	if listOnly {
		quizArray = quizzesListSimple
	} else {
		quizArray = quizzesListFull
	}

	w.Header().Set("Content-Type", "application/json") // normal header
	w.WriteHeader(http.StatusOK)

	jsonStr, err := json.Marshal(quizArray)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getQuiz(quizId string) *quiz.Quiz {
	return quizzes[quizId]
}

func restHandleQuizById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	quizId := ps.ByName(PATH_PARAM_QUIZ_ID)
	if quizId == "" {
		// This makes no sense. restHandleQuizAll() should have been called.
		http.Error(w, "Empty quiz ID", http.StatusInternalServerError)
		return
	}

	q := getQuiz(quizId)
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

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func restHandleQuizSectionsByQuizId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	listOnly := false
	queryValues := r.URL.Query()
	if queryValues != nil {
		listOnlyStr := queryValues.Get(QUERY_PARAM_LIST_ONLY)
		listOnly, _ = strconv.ParseBool(listOnlyStr)
	}

	quizId := ps.ByName(PATH_PARAM_QUIZ_ID)
	if quizId == "" {
		// This makes no sense. restHandleQuizAll() should have been called.
		http.Error(w, "Empty quiz ID", http.StatusInternalServerError)
		return
	}

	q := getQuiz(quizId)
	if q == nil {
		http.Error(w, "quiz not found", http.StatusInternalServerError)
		return
	}

	sections := q.Sections
	if listOnly {
		simpleSections := make([]*quiz.Section, 0, len(sections))
		for _, s := range sections {
			var simple quiz.Section
			s.CopyHasIdAndTitle(&simple.HasIdAndTitle)
			simpleSections = append(simpleSections, &simple)
		}

		sections = simpleSections
	}

	jsonStr, err := json.Marshal(sections)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func restHandleQuizQuestionById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	quizId := ps.ByName(PATH_PARAM_QUIZ_ID)
	if quizId == "" {
		// This makes no sense. restHandleQuizAll() should have been called.
		http.Error(w, "Empty quiz ID", http.StatusInternalServerError)
		return
	}

	questionId := ps.ByName(PATH_PARAM_QUESTION_ID)
	if questionId == "" {
		// This makes no sense.
		http.Error(w, "Empty question ID", http.StatusInternalServerError)
		return
	}

	q := getQuiz(quizId)
	if q == nil {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	}

	qa := q.GetQuestionAndAnswer(questionId)
	if qa == nil {
		http.Error(w, "question not found", http.StatusInternalServerError)
		return
	}

	qa.Question.SetQuestionExtras(q)

	jsonStr, err := json.Marshal(qa.Question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
