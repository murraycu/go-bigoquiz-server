package user

import "golang.org/x/oauth2"

type Profile struct {
	Id string
	Name string
	Email string

	// Google's "sub" ID. See https://developers.google.com/identity/protocols/OpenIDConnect#obtainuserinfo
	GoogleId string

	GoogleAccessToken oauth2.Token
}
