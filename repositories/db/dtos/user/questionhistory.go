package user

type QuestionHistory struct {
	QuestionId string `datastore:"questionId"`

	AnsweredCorrectlyOnce bool `datastore:"answeredCorrectlyOnce"`

	//Decrements once for each time the user answers it correctly.
	//Increments once for each time the user answers it wrongly.
	CountAnsweredWrong int `datastore:"countAnsweredWrong"`
}
