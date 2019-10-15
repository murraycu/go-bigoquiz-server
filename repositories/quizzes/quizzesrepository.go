package quizzes

import (
	"fmt"
	"github.com/murraycu/go-bigoquiz-server/domain/quiz"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

// See https://gobyexample.com/sorting-by-functions
type quizListByTitle []*quiz.Quiz

func (s quizListByTitle) Len() int {
	return len(s)
}

func (s quizListByTitle) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s quizListByTitle) Less(i, j int) bool {
	return s[i].Title < s[j].Title
}

type QuizzesRepository struct {
}

// Map of quiz IDs to Quiz.
type mapQuizzes map[string]*quiz.Quiz

type mapList []*quiz.Quiz

type QuizzesAndCaches struct {
	Quizzes           mapQuizzes
	QuizzesListSimple mapList
	QuizzesListFull   mapList
}

func NewQuizzesRepository() (*QuizzesRepository, error) {
	result := &QuizzesRepository{}

	return result, nil
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

func loadQuiz(id string) (*quiz.Quiz, error) {
	absFilePath, err := filepath.Abs("quizzes/" + id + ".xml")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return quiz.LoadQuiz(absFilePath, id)
}

func loadQuizzes() (mapQuizzes, error) {
	quizzes := make(mapQuizzes, 0)

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
func buildQuizzesSimple(quizzes mapQuizzes) mapList {
	// Create a slice with the same capacity.
	result := make(mapList, 0, len(quizzes))

	for _, q := range quizzes {
		var simple quiz.Quiz
		q.CopyHasIdAndTitle(&simple.HasIdAndTitle)
		simple.IsPrivate = q.IsPrivate

		result = append(result, &simple)
	}

	sort.Sort(quizListByTitle(result))

	return result
}

func buildQuizzesFull(quizzes mapQuizzes) mapList {
	// Create a slice with the same capacity.
	result := make(mapList, 0, len(quizzes))

	for _, q := range quizzes {
		result = append(result, q)
	}

	sort.Sort(quizListByTitle(result))

	return result
}

func (q *QuizzesRepository) GetQuizzesAndCaches() (*QuizzesAndCaches, error) {
	var err error

	var result QuizzesAndCaches
	result.Quizzes, err = loadQuizzes()
	if err != nil {
		return nil, fmt.Errorf("Could not load quiz files: %v", err)
	}

	result.QuizzesListSimple = buildQuizzesSimple(result.Quizzes)
	result.QuizzesListFull = buildQuizzesFull(result.Quizzes)

	return &result, nil
}
