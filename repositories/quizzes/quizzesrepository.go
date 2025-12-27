package quizzes

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	domainquiz "github.com/murraycu/go-bigoquiz-server/domain/quiz"
	dtoquiz "github.com/murraycu/go-bigoquiz-server/repositories/quizzes/dtos/quiz"
)

type QuizzesRepository struct {
	directoryPath string
}

// Map of quiz IDs to Quiz.
type MapQuizzes map[string]*domainquiz.Quiz

type QuizzesAndCaches struct {
	Quizzes MapQuizzes
}

// NewQuizzesRepository creates a new quizzes repository.
//
// directoryPath is the path to a directory containing quizzes in JSON files.
func NewQuizzesRepository(directoryPath string) (*QuizzesRepository, error) {
	result := &QuizzesRepository{}
	result.directoryPath = directoryPath

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

func loadQuizAsDto(directoryFilepath string, id string) (*dtoquiz.Quiz, error) {
	fullPath := filepath.Join(directoryFilepath, id+".json")
	absFilePath, err := filepath.Abs(fullPath)
	if err != nil {
		return nil, err
	}

	return dtoquiz.LoadQuiz(absFilePath, id)
}

type MapDtoQuizzes map[string]*dtoquiz.Quiz

func loadQuizzesAsDto(directoryFilepath string) (MapDtoQuizzes, error) {
	quizzes := make(MapDtoQuizzes, 0)

	quizNames, err := filesWithExtension(directoryFilepath, "json")
	if err != nil {
		fmt.Println(err)
		return quizzes, err
	}

	for _, name := range quizNames {
		q, err := loadQuizAsDto(directoryFilepath, name)
		if err != nil {
			fmt.Println(err)
			return quizzes, err
		}

		quizzes[q.Id] = q
	}

	return quizzes, nil
}

func (q *QuizzesRepository) LoadQuizzes() (MapQuizzes, error) {
	quizzes, err := loadQuizzesAsDto(q.directoryPath)
	if err != nil {
		return nil, fmt.Errorf("could not load quiz files: %v", err)
	}

	result, err := convertDtoQuizzesToDomainQuizzes(quizzes)
	if err != nil {
		return nil, fmt.Errorf("could not convert DTO quizzes to domain quizzes: %v", err)
	}

	return result, nil
}
