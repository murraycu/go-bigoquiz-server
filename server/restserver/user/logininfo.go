package user

type LoginInfo struct {
	LoggedIn bool `json:"loggedIn"`

	LoginUrl  string `json:"loginUrl,omitempty"`
	LogoutUrl string `json:"logoutUrl,omitempty"`

	Nickname string `json:"nickname,omitempty"`

	// If the user account is linked to these oauth2 accounts:
	GoogleLinked     bool   `json:"googleLinked"`
	GoogleProfileUrl string `json:"googleProfileUrl"`
	// GoogleTokenExpired will be true if we know the login details but the OAuth token has expired (and could be re-obtained)
	GoogleTokenExpired bool `json:"googleTokenExpired"`

	GitHubLinked     bool   `json:"gitHubLinked"`
	GitHubProfileUrl string `json:"gitHubProfileUrl"`
	// GitHubTokenExpired will be true if we know the login details but the OAuth token has expired (and could be re-obtained)
	GitHubTokenExpired bool `json:"gitHubProfileExpired"`

	FacebookLinked     bool   `json:"facebookLinked"`
	FacebookProfileUrl string `json:"facebookProfileUrl"`
	// FacebookTokenExpired will be true if we know the login details but the OAuth token has expired (and could be re-obtained)
	FacebookTokenExpired bool `json:"facebookTokenExpired"`

	// This is just for debugging.
	ErrorMessage string `json:"errorMessage,omitempty"`
}
