package user

type QuestionHistory struct {
	QuestionId string

	AnsweredCorrectlyOnce bool

	//Decrements once for each time the user answers it correctly.
	//Increments once for each time the user answers it wrongly.
	CountAnsweredWrong int
}
