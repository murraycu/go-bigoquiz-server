package quiz

type Question struct {
	HasIdAndTitle
	Text Text `json:"text,omitempty" xml:"text"`
}
