package quiz

type Text struct {
	// Note: using ',innerxml' instead of chardata would not unescape the text/xml.
	Text string

	IsHtml bool
}
