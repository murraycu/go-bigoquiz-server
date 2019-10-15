package quiz

type HasIdAndTitle struct {
	Id    string `json:"id"`
	Title string `json:"title,omitempty"`

	// A URL.
	Link string `json:"link,omitempty"`
}

func (self *HasIdAndTitle) CopyHasIdAndTitle(dest *HasIdAndTitle) {
	dest.Id = self.Id
	dest.Title = self.Title
	dest.Link = self.Link
}
