package oauthparsers

/** A representation of some of the JSON returned by
 * Github.
 * See https://developer.github.com/v3/users/#get-a-single-user
 */
type GitHubUserInfo struct {
	// The unique ID of the Github user.
	Id int `json:"id"`

	Name       string `json:"name"`
	Email      string `json:"email"`
	ProfileUrl string `json:"html_url"`
}
