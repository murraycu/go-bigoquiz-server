package bigoquiz

/** A representation of some of the JSON returned by
 * Google.
 */
type GoogleUserInfo struct {
  Name string `json:"name"`
  Picture string `json:"picture"`
  Email string `json:"email"`
  EmailVerified bool `json:"email_verified"`
}
