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

	// This is just for debugging.
	ErrorMessage string `json:"errorMessage,omitempty"`
}
