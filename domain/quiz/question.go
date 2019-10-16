package quiz

type Question struct {
	Id string

	// A URL.
	Link string

	Text Text

	QuizUsesMathML bool
}
