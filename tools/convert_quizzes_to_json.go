package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/murraycu/go-bigoquiz-server/repositories/quizzes"
	dtoquiz "github.com/murraycu/go-bigoquiz-server/repositories/quizzes/dtos/quiz"
)

func simplifyQuestionText(question *dtoquiz.Question) error {
	if question.TextSimple != "" {
		// It's already simple.
		return nil
	}

	if question.TextDetail.Text == "" {
		// There's no detailed text to simplify.
		return nil
	}

	if !question.TextDetail.IsHtml {
		question.TextSimple = question.TextDetail.Text
		question.TextDetail = dtoquiz.Text{}
	}

	return nil
}

func simplifyQuestionsAndAnswersText(questions []*dtoquiz.QuestionAndAnswer) error {
	for _, qa := range questions {
		err := simplifyQuestionText(&qa.Question)
		if err != nil {
			return fmt.Errorf("simplifyQuestionText() failed: %v", err)
		}

		if qa.AnswerSimple != "" {
			// It's already simple.
			continue
		}

		if qa.AnswerDetail.Text == "" {
			// There's no detailed text to simplify.
			continue
		}

		if !qa.AnswerDetail.IsHtml {
			qa.AnswerSimple = qa.AnswerDetail.Text
			qa.AnswerDetail = dtoquiz.Text{}
		}
	}

	return nil
}

func simplifyQuizText(quiz *dtoquiz.Quiz) error {
	err := simplifyQuestionsAndAnswersText(quiz.Questions)
	if err != nil {
		return fmt.Errorf("simplifyQuestionText() failed: %v", err)
	}

	for _, section := range quiz.Sections {
		err := simplifyQuestionsAndAnswersText(section.Questions)
		if err != nil {
			return fmt.Errorf("simplifyQuestionText() failed: %v", err)
		}

		for _, subSection := range section.SubSections {
			err := simplifyQuestionsAndAnswersText(subSection.Questions)
			if err != nil {
				return fmt.Errorf("simplifyQuestionText() failed: %v", err)
			}
		}
	}

	return nil
}

func saveQuizAsJsonDto(quiz dtoquiz.Quiz, directoryFilepath string) error {
	id := quiz.Id
	fullPath := filepath.Join(directoryFilepath, id+".json")
	absFilePath, err := filepath.Abs(fullPath)
	if err != nil {
		return fmt.Errorf("Failed to get absolute filepath for json quiz: %v\n", err)
	}

	data, err := json.MarshalIndent(quiz, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to marshal quiz: %v\n", err)
	}

	err = os.WriteFile(absFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("Failed to write quiz to file: %v\n", err)
	}

	return nil
}

func main() {
	sourceDir, err := filepath.Abs("quizzes")
	if err != nil {
		log.Fatalf("could not get absolute filepath for quizzes: %v", err)
	}

	dtoQuizzes, err := quizzes.LoadQuizzesAsDto(sourceDir /*asJson=*/, false /*addReverses=*/, false)
	if err != nil {
		log.Fatalf("LoadQuizzes() failed: %v", err)
	}

	destDir, err := filepath.Abs("quizzes")
	if err != nil {
		log.Fatalf("could not get absolute filepath for json quizzes: %v", err)
	}

	for _, quiz := range dtoQuizzes {
		print(fmt.Sprintf("Processing quiz: %s\n", quiz.Title))

		err := simplifyQuizText(quiz)
		if err != nil {
			log.Fatalf("simplifyText() failed: %v\n", err)
		}

		err = saveQuizAsJsonDto(*quiz, destDir)
		if err != nil {
			log.Fatalf("saveQuizAsJsonDto() failed: %v\n", err)
		}

	}
}
