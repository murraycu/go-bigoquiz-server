package quiz

type Quiz struct {
	HasIdAndTitle
	IsPrivate bool `json:"isPrivate"`

	Sections []*Section `json:"sections,omitempty"`
	// (Unlike the DTO struct, this only has questions inside a Section or SubSection.)

	UsesMathML bool `json:"usesMathML"`
}
