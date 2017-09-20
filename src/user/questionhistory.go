package user

type QuestionHistory struct {
	QuestionId string `json:"questionId,omitempty" datastore:"questionId"`

	AnsweredCorrectlyOnce bool `json:"answeredCorrectlyOnce" datastore:"answeredCorrectlyOnce"`

	//Decrements once for each time the user answers it correctly.
	//Increments once for each time the user answers it wrongly.
	CountAnsweredWrong int `json:"countAnsweredWrong" datastore:"countAnsweredWrong"`
}
