package user

import "google.golang.org/appengine/datastore"

type LoginInfo struct {
	// Key *datastore.Key `json:"-"`

	LoggedIn bool `json:"loggedIn"`

	LoginUrl  string `json:"loginUrl,omitempty"`
	LogoutUrl string `json:"logoutUrl,omitempty"`

	// TODO: Show a string-based version of this,
	// for human-readability and convenience?
	UserId *datastore.Key `json:"-"`

	Nickname string `json:"nickname,omitempty"`

	// If the user account is linked to these oauth2 accounts:
	GoogleLinked       bool   `json:"googleLinked"`
	GoogleProfileUrl   string `json:"googleProfileUrl"`
	GitHubLinked       bool   `json:"gitHubLinked"`
	GitHubProfileUrl   string `json:"gitHubProfileUrl"`
	FacebookLinked     bool   `json:"facebookLinked"`
	FacebookProfileUrl string `json:"facebookProfileUrl"`

	// This is just for debugging.
	ErrorMessage string `json:"errorMessage,omitempty"`
}
