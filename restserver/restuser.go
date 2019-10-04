package restserver

import (
	"cloud.google.com/go/datastore"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/murraycu/go-bigoquiz-server/db"
	"github.com/murraycu/go-bigoquiz-server/user"
	"golang.org/x/oauth2"
	"net/http"
)

func (s *RestServer) HandleUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	loginInfo, err := s.getLoginInfoFromSessionAndDb(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.Marshal(loginInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *RestServer) getProfileFromSessionAndDb(r *http.Request) (*user.Profile, *datastore.Key, *oauth2.Token, error) {
	userId, token, err := s.UserSessionStore.GetProfileFromSession(r)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("GetProfileFromSession() failed: %v", err)
	}

	if userId == nil {
		// Not an error.
		// It's just not in the session cookie.
		return nil, nil, nil, nil
	}

	dbClient, err := db.NewUserDataRepository()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("NewUserDataRepository() failed: %v", err)
	}

	c := r.Context()
	profile, err := dbClient.GetUserProfileById(c, userId)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("GetUserProfileById() failed: %v", err)
	}

	return profile, userId, token, nil
}

func (s *RestServer) getLoginInfoFromSessionAndDb(r *http.Request) (*user.LoginInfo, error) {
	var loginInfo user.LoginInfo

	profile, userId, token, err := s.getProfileFromSessionAndDb(r)
	if err != nil {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = fmt.Sprintf("not logged in (%v)", err)
	}

	s.updateLoginInfoFromProfile(&loginInfo, profile, token, userId)

	return &loginInfo, err
}

func (s *RestServer) updateLoginInfoFromProfile(loginInfo *user.LoginInfo, profile *user.Profile, token *oauth2.Token, userId *datastore.Key) {
	if profile == nil {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (no profile found)"
	} else if !token.Valid() {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (invalid token)"
	} else {
		loginInfo.LoggedIn = true
		loginInfo.Nickname = profile.Name
		loginInfo.UserId = userId // Not for the JSON, but useful to callers.

		loginInfo.GoogleLinked = profile.GoogleId != ""
		loginInfo.GoogleProfileUrl = profile.GoogleProfileUrl
		loginInfo.GitHubLinked = profile.GitHubId != 0
		loginInfo.GitHubProfileUrl = profile.GitHubProfileUrl
		loginInfo.FacebookLinked = profile.FacebookId != ""
		loginInfo.FacebookProfileUrl = profile.FacebookProfileUrl
	}
}

/** Get the user ID.
 * Returns a nil Key, and a nil error, if the user is not logged in.
 */
func (s *RestServer) getUserIdFromSessionAndDb(r *http.Request, w http.ResponseWriter) (*datastore.Key, error) {
	loginInfo, err := s.getLoginInfoFromSessionAndDb(r)
	if err != nil {
		return nil, fmt.Errorf("getLoginInfoFromSessionAndDb() failed: %v", err)
	}

	return loginInfo.UserId, nil
}
