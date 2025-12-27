package quiz

type Question struct {
	Id string `xml:"id,attr" json:"id,omitempty"`

	// A URL.
	Link string `xml:"link" json:"link,omitempty"`

	// TextDetail is an alternative to TextSimple.
	// Only one of these should be set.
	TextDetail Text `xml:"text" json:"textDetail,omitempty,omitzero"`

	// TextSimple is an alternative to TextDetail.
	// Only one of these should be set.
	TextSimple string `json:"text,omitempty"`
}
