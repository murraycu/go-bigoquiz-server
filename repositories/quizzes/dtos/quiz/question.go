package quiz

type Question struct {
	Id string `xml:"id,attr" json:"id,omitempty"`

	// A URL.
	Link string `xml:"link" json:"link,omitempty"`

	TextDetail Text `xml:"text" json:"textDetail,omitempty"`
}
