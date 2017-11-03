package quiz

type Text struct {
	// Note: using ',innerxml' instead of chardata would not unescape the text/xml.
	Text string `json:"text,omitempty" xml:",chardata"`

	IsHtml bool `json:"isHtml,omitempty" xml:"is_html,attr"`
}
