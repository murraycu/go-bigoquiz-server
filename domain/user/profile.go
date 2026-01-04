package user

type Profile struct {
	Name   string
	Email  string
	UserId string

	GoogleLinked     bool
	GoogleProfileUrl string

	GitHubLinked     bool
	GitHubProfileUrl string

	FacebookLinked     bool
	FacebookProfileUrl string
}
