package user

type LoginInfo struct {
	LoggedIn bool `json:"loggedIn"`

	LoginUrl  string `json:"loginUrl,omitempty"`
	LogoutUrl string `json:"logoutUrl,omitempty"`

	UserId   string `json:"userId,omitempty"`
	Nickname string `json:"nickname,omitempty"`
}
