package quiz

type HasIdAndTitle struct {
	Id    string `xml:"id,attr"`
	Title string `xml:"title"`

	// A URL.
	Link string `xml:"link"`
}
