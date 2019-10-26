package user

// Statistics for the user, either per quiz (SectionId is then empty), or per section in a quiz.
type Stats struct {
	QuizId    string `json:"quizId,omitEmpty"`
	SectionId string `json:"sectionId,omitEmpty"`

	Answered int `json:"answered"`
	Correct  int `json:"correct"`

	CountQuestionsAnsweredOnce int `json:"countQuestionsAnsweredOnce"`
	CountQuestionsCorrectOnce  int `json:"countQuestionsCorrectOnce"`

	QuestionHistories []QuestionHistory `json:"questionHistories,omitEmpty"`

	// These are from the quiz, for convenience
	// so they don't need to be in the database.
	CountQuestions int    `json:"countQuestions"`
	QuizTitle      string `json:"quizTitle,omitEmpty"`
	SectionTitle   string `json:"sectionTitle,omitEmpty"`
}
