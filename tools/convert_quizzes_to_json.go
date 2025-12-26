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

		err := saveQuizAsJsonDto(*quiz, destDir)
		if err != nil {
			log.Fatalf("saveQuizAsJsonDto() failed: %v\n", err)
		}

	}
}
