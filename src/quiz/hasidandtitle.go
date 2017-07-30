package quiz

type HasIdAndTitle struct {
	Id    string `json:"id" xml:"id,attr"`
	Title string `json:"title,omitempty" xml:"title"`

	// A URL.
	Link string `json:"link,omitempty" xml:"link"`
}
