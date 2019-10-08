package oauthparsers

/** A representation of some of the JSON returned by
 * Google.
 */
type GoogleUserInfo struct {
	// The unique ID of the Google user.
	Sub string `json:"sub"` // See https://developers.google.com/identity/protocols/OpenIDConnect#obtaininguserprofileinformation

	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`

	// For instance https://plus.google.com/+MurrayCumming
	ProfileUrl string `json:"profile"`
}
