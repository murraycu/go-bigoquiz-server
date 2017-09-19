package quiz

type QuestionAndAnswer struct {
	Question `json:"question,omitempty"` // TODO: Name this field (as Question) and xtill unmarhal the XML properly.
	Answer   Text                        `json:"answer,omitempty" xml:"answer"`
}

func (self *QuestionAndAnswer) createReverse() *QuestionAndAnswer {
	var result QuestionAndAnswer
	result.Id = "reverse-" + self.Id
	result.Text = self.Answer
	result.Answer = self.Text
	result.Answer.IsHtml = false
	return &result
}
