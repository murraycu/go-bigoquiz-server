package quiz

type QuestionAndAnswer struct {
	Question
	Answer Text `xml:"answer" json:"answer"`
}

func (self *QuestionAndAnswer) createReverse() *QuestionAndAnswer {
	var result QuestionAndAnswer
	result.Id = "reverse-" + self.Id
	result.Text = self.Answer
	result.Answer = self.Text
	result.Answer.IsHtml = false
	return &result
}
