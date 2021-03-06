package quiz

type HasIdAndTitle struct {
	Id    string
	Title string

	// A URL.
	Link string
}

func (self *HasIdAndTitle) CopyHasIdAndTitle(dest *HasIdAndTitle) {
	dest.Id = self.Id
	dest.Title = self.Title
	dest.Link = self.Link
}
