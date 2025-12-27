package quiz

type QuestionAndAnswer struct {
	Question
	AnswerDetail Text `xml:"answerdetail" json:"answerdetail"`
}

func (self *QuestionAndAnswer) createReverse() *QuestionAndAnswer {
	var result QuestionAndAnswer
	result.Id = "reverse-" + self.Id
	result.TextDetail = self.AnswerDetail

	if self.TextSimple != "" {
		result.AnswerDetail.Text = self.TextSimple
	} else {
		result.AnswerDetail.Text = self.TextDetail.Text
	}

	result.AnswerDetail.IsHtml = false
	return &result
}
