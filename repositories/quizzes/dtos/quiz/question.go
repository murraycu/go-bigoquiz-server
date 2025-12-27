package quiz

type Question struct {
	Id string `xml:"id,attr" json:"id,omitempty"`

	// A URL.
	Link string `xml:"link" json:"link,omitempty"`

	Text Text `xml:"text" json:"text,omitempty"`
}
