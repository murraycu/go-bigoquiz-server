package quiz

type SubSection struct {
	HasIdAndTitle
	Questions   []*QuestionAndAnswer `json:"questions,omitempty" xml:"question"`
}
