package quiz

type QuestionAndAnswer struct {
	Question
	Answer Text `xml:"answer" json:"answer"`
}

func (self *QuestionAndAnswer) createReverse() *QuestionAndAnswer {
	var result QuestionAndAnswer
	result.Id = "reverse-" + self.Id
	result.TextDetail = self.Answer

	if self.TextSimple != "" {
		result.Answer.Text = self.TextSimple
	} else {
		result.Answer.Text = self.TextDetail.Text
	}

	result.Answer.IsHtml = false
	return &result
}
