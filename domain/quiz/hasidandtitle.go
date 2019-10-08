package quiz

type HasIdAndTitle struct {
	Id    string `json:"id" xml:"id,attr"`
	Title string `json:"title,omitempty" xml:"title"`

	// A URL.
	Link string `json:"link,omitempty" xml:"link"`
}

func (self *HasIdAndTitle) CopyHasIdAndTitle(dest *HasIdAndTitle) {
	dest.Id = self.Id
	dest.Title = self.Title
	dest.Link = self.Link
}
