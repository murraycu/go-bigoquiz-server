package quiz

type HasIdAndTitle struct {
	Id    string `xml:"id,attr" json:"id,omitempty"`
	Title string `xml:"title" json:"title,omitempty"`

	// A URL.
	Link string `xml:"link" json:"link,omitempty"`
}
