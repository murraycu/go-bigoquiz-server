package db

/** A representation of some of the JSON returned by
 * Github.
 * See https://developer.github.com/v3/users/#get-a-single-user
 */
type FacebookUserInfo struct {
	// The unique ID of the Facebook user.
	Id string `json:"id"`

	Name       string `json:"name"`
	Email      string `json:"email"`
	ProfileUrl string `json:"link"`
}
