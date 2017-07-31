package quiz

type QuestionAndAnswer struct {
	Question `json:"question,omitempty"` // TODO: Name this field (as Question) and till unmarhal the XML properly.
	Answer   Text                        `json:"answer,omitempty" xml:"answer"`
}
