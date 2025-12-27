package quiz

type QuestionAndAnswer struct {
	Question

	// AnswerDetail is an alternative to AnswerSimple.
	// Only one of these should be set.
	AnswerDetail Text `xml:"answer" json:"answerDetail,omitempty,omitzero"`

	// AnswerSimple is an alternative to AnswerDetail.
	// Only one of these should be set.
	AnswerSimple string `json:"answer,omitempty"`
}

func (self *QuestionAndAnswer) createReverse() *QuestionAndAnswer {
	var result QuestionAndAnswer
	result.Id = "reverse-" + self.Id

	// Copy the answer to the question.
	if self.AnswerSimple != "" {
		result.TextDetail.Text = self.AnswerSimple
	} else {
		result.TextDetail = self.AnswerDetail
	}

	// Copy the question to the answer.
	if self.TextSimple != "" {
		result.AnswerDetail.Text = self.TextSimple
	} else {
		result.AnswerDetail.Text = self.TextDetail.Text
	}

	result.AnswerDetail.IsHtml = false
	return &result
}
