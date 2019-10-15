package quiz

type Question struct {
	Id string `xml:"id,attr"`

	// A URL.
	Link string `xml:"link"`

	Text Text `xml:"text"`
}
