package user

import "golang.org/x/oauth2"

type Profile struct {
	Name  string `datastore:"name"`
	Email string `datastore:"email"`

	// Google's "sub" ID. See https://developers.google.com/identity/protocols/OpenIDConnect#obtainuserinfo
	GoogleId string `datastore:"googleId"`

	// GoogleAccessToken is actually an oauth2.Token, not an access token, but contains an access token (and a refresh token).
	GoogleAccessToken oauth2.Token `datastore:"googleAccessToken"`
	GoogleProfileUrl  string       `datastore:"googleProfileUrl"`

	// GitHub's ID. See https://developer.github.com/v3/users/#get-a-single-user
	GitHubId int `datastore:"gitHubId"`

	// GitHubAccessToken is actually an oauth2.Token, not an access token, but contains an access token (and a refresh token).
	GitHubAccessToken oauth2.Token `datastore:"gitHubAccessToken"`
	GitHubProfileUrl  string       `datastore:"gitHubProfileUrl"`

	// Facebook's ID. See TODO
	FacebookId string `datastore:"facebookId"`

	// FacebookAccessToken is actually an oauth2.Token, not an access token, but contains an access token (and a refresh token).
	FacebookAccessToken oauth2.Token `datastore:"facebookAccessToken"`
	FacebookProfileUrl  string       `datastore:"facebookProfileUrl"`
}
