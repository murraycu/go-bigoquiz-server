package restserver

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	domainuser "github.com/murraycu/go-bigoquiz-server/domain/user"
	"github.com/murraycu/go-bigoquiz-server/repositories/db"
	restuser "github.com/murraycu/go-bigoquiz-server/server/restserver/user"
	"golang.org/x/oauth2"
	"net/http"
)

func (s *RestServer) HandleUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	loginInfo, _, err := s.getLoginInfoFromSessionAndDb(r)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonStr, err := json.Marshal(loginInfo)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, err.Error())
	}
}

// Returns the LoginInfo and the userID.
func (s *RestServer) getProfileFromSessionAndDb(r *http.Request) (*domainuser.Profile, string, *oauth2.Token, error) {
	userId, token, err := s.userSessionStore.GetProfileFromSession(r)
	if err != nil {
		return nil, "", nil, fmt.Errorf("GetProfileFromSession() failed: %v", err)
	}

	if len(userId) == 0 {
		// Not an error.
		// It's just not in the session cookie.
		return nil, "", nil, nil
	}

	dbClient, err := db.NewUserDataRepository()
	if err != nil {
		return nil, "", nil, fmt.Errorf("NewUserDataRepository() failed: %v", err)
	}

	c := r.Context()
	profile, err := dbClient.GetUserProfileById(c, userId)
	if err != nil {
		return nil, "", nil, fmt.Errorf("GetUserProfileById() failed: %v", err)
	}

	return profile, userId, token, nil
}

// Returns the LoginInfo and the userID.
func (s *RestServer) getLoginInfoFromSessionAndDb(r *http.Request) (*restuser.LoginInfo, string, error) {
	var loginInfo restuser.LoginInfo

	profile, userId, token, err := s.getProfileFromSessionAndDb(r)
	if err != nil {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = fmt.Sprintf("not logged in (%v)", err)
	}

	s.updateLoginInfoFromProfile(&loginInfo, profile, token)

	return &loginInfo, userId, err
}

func (s *RestServer) updateLoginInfoFromProfile(loginInfo *restuser.LoginInfo, profile *domainuser.Profile, token *oauth2.Token) {
	if profile == nil {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (no profile found)"
	} else if !token.Valid() {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (invalid token)"
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
func (s *RestServer) getUserIdFromSessionAndDb(r *http.Request, w http.ResponseWriter) (string, error) {
	userId, _, err := s.userSessionStore.GetProfileFromSession(r)
	if err != nil {
		return "", fmt.Errorf("GetProfileFromSession() failed: %v", err)
	}

	return userId, nil
}
