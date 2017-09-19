package user

import "golang.org/x/oauth2"

type Profile struct {
	Name  string `datastore:"name"`
	Email string `datastore:"email"`

	// Google's "sub" ID. See https://developers.google.com/identity/protocols/OpenIDConnect#obtainuserinfo
	GoogleId string `datastore:"googleId"`

	GoogleAccessToken oauth2.Token `datastore:"googleAccessToken"`
}
