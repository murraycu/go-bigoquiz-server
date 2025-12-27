package quiz

type QuestionAndAnswer struct {
	Question
	Answer Text `xml:"answer" json:"answer"`
}

func (self *QuestionAndAnswer) createReverse() *QuestionAndAnswer {
	var result QuestionAndAnswer
	result.Id = "reverse-" + self.Id
	result.TextDetail = self.Answer
	result.Answer = self.TextDetail
	result.Answer.IsHtml = false
	return &result
}
