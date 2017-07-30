package quiz

type Section struct {
	HasIdAndTitle
	Questions []Question `json:"questions,omitempty" xml:"question"`
}
