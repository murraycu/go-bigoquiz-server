package quiz

type HasIdAndTitle struct {
	Id    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`

	// A URL.
	Link string `json:"link,omitempty"`
}
