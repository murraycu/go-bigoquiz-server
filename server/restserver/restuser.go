package restserver

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	domainuser "github.com/murraycu/go-bigoquiz-server/domain/user"
	restuser "github.com/murraycu/go-bigoquiz-server/server/restserver/user"
)

func (s *RestServer) HandleUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	loginInfo, _, err := s.getLoginInfoFromSessionAndDb(w, r)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "getLoginInfoFromSessionAndDb() failed: %v", err)
		return
	}

	marshalAndWriteOrHttpError(w, &loginInfo)
}

// Returns the LoginInfo and the userID.
func (s *RestServer) getProfileFromSessionAndDb(w http.ResponseWriter, r *http.Request) (*domainuser.Profile, string, error) {
	userId, err := s.getUserIdFromSessionAndDb(w, r)
	if err != nil {
		return nil, "", fmt.Errorf("getUserIdFromSessionAndDb() failed: %v", err)
	}

	if len(userId) == 0 {
		// Not an error.
		// It's just not in the session cookie.
		return nil, "", nil
	}

	c := r.Context()
	profile, err := s.userDataClient.GetUserProfileById(c, userId)
	if err != nil {
		return nil, "", fmt.Errorf("GetUserProfileById() failed: %v", err)
	}

	return profile, userId, nil
}

// Returns the LoginInfo and the userID.
func (s *RestServer) getLoginInfoFromSessionAndDb(w http.ResponseWriter, r *http.Request) (*restuser.LoginInfo, string, error) {
	var loginInfo restuser.LoginInfo

	profile, userId, err := s.getProfileFromSessionAndDb(w, r)
	if err != nil {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = fmt.Sprintf("not logged in (%v)", err)
	}

	s.updateLoginInfoFromProfile(&loginInfo, profile)

	return &loginInfo, userId, err
}

func (s *RestServer) updateLoginInfoFromProfile(loginInfo *restuser.LoginInfo, profile *domainuser.Profile) {
	if profile == nil {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (no profile found)"
	} else {
		loginInfo.LoggedIn = true
		loginInfo.Nickname = profile.Name

		loginInfo.GoogleLinked = len(profile.GoogleProfileUrl) != 0
		loginInfo.GoogleProfileUrl = profile.GoogleProfileUrl
		loginInfo.GitHubLinked = len(profile.GitHubProfileUrl) != 0
		loginInfo.GitHubProfileUrl = profile.GitHubProfileUrl
		loginInfo.FacebookLinked = len(profile.FacebookProfileUrl) != 0
		loginInfo.FacebookProfileUrl = profile.FacebookProfileUrl
	}
}

/** Get the user ID.
 * Returns an empty user ID, and a nil error, if the user is not logged in.
 */
func (s *RestServer) getUserIdFromSessionAndDb(w http.ResponseWriter, r *http.Request) (string, error) {
	userId, token, err := s.userSessionStore.GetProfileFromSession(r)
	if err != nil {
		return "", fmt.Errorf("GetProfileFromSession() failed: %v", err)
	}

	if !token.Valid() {
		// TODO: Revalidate it.

		// This is not an error
		// (it is normal for a token to expire.)
		// (the user is now not logged in.)
		return "", nil
	}

	return userId, nil
}
