package quizzes

import (
	"fmt"
	domainquiz "github.com/murraycu/go-bigoquiz-server/domain/quiz"
	dtoquiz "github.com/murraycu/go-bigoquiz-server/repositories/quizzes/dtos/quiz"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type QuizzesRepository struct {
}

// Map of quiz IDs to Quiz.
type MapQuizzes map[string]*domainquiz.Quiz

type QuizzesAndCaches struct {
	Quizzes MapQuizzes
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

func loadQuiz(id string) (*dtoquiz.Quiz, error) {
	absFilePath, err := filepath.Abs("quizzes/" + id + ".xml")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return dtoquiz.LoadQuiz(absFilePath, id)
}

type MapDtoQuizzes map[string]*dtoquiz.Quiz

func loadQuizzes() (MapDtoQuizzes, error) {
	quizzes := make(MapDtoQuizzes, 0)

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

func (q *QuizzesRepository) GetQuizzes() (MapQuizzes, error) {
	quizzes, err := loadQuizzes()
	if err != nil {
		return nil, fmt.Errorf("could not load quiz files: %v", err)
	}

	result, err := convertDtoQuizzesToDomainQuizzes(quizzes)
	if err != nil {
		return nil, fmt.Errorf("could not convert DTO quizzes to domain quizzes: %v", err)
	}

	return result, nil
}
