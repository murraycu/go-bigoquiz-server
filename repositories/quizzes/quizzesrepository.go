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

/** directoryPath is the path to a directory containing quizzes in XML files.
 */
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

func loadQuizAsDto(directoryFilepath string, asJson bool, addReverses bool, id string) (*dtoquiz.Quiz, error) {
	ext := ".xml"
	if asJson {
		ext = ".json"
	}

	fullPath := filepath.Join(directoryFilepath, id+ext)
	absFilePath, err := filepath.Abs(fullPath)
	if err != nil {
		return nil, err
	}

	return dtoquiz.LoadQuiz(absFilePath, asJson, addReverses, id)
}

type MapDtoQuizzes map[string]*dtoquiz.Quiz

func LoadQuizzesAsDto(directoryFilepath string, asJson bool, addReverses bool) (MapDtoQuizzes, error) {
	quizzes := make(MapDtoQuizzes, 0)

	quizNames, err := filesWithExtension(directoryFilepath, "xml")
	if err != nil {
		fmt.Println(err)
		return quizzes, err
	}

	for _, name := range quizNames {
		q, err := loadQuizAsDto(directoryFilepath, asJson, addReverses, name)
		if err != nil {
			fmt.Println(err)
			return quizzes, err
		}

		quizzes[q.Id] = q
	}

	return quizzes, nil
}

func (q *QuizzesRepository) LoadQuizzes(asJson bool, addReverses bool) (MapQuizzes, error) {
	quizzes, err := LoadQuizzesAsDto(q.directoryPath, asJson, addReverses)
	if err != nil {
		return nil, fmt.Errorf("could not load quiz files: %v", err)
	}

	result, err := convertDtoQuizzesToDomainQuizzes(quizzes)
	if err != nil {
		return nil, fmt.Errorf("could not convert DTO quizzes to domain quizzes: %v", err)
	}

	return result, nil
}
