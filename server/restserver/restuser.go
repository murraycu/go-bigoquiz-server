package restserver

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	domainuser "github.com/murraycu/go-bigoquiz-server/domain/user"
	restuser "github.com/murraycu/go-bigoquiz-server/server/restserver/user"
)

func (s *RestServer) HandleUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	result, err := s.getLoginInfoFromSessionAndDb(w, r)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "getLoginInfoFromSessionAndDb() failed: %v", err)
		return
	}

	marshalAndWriteOrHttpError(w, result.LoginInfo)
}

type getProfileResult struct {
	Profile *domainuser.Profile
	UserId  string
}

// Returns the LoginInfo and the userID.
func (s *RestServer) getProfileFromSessionAndDb(w http.ResponseWriter, r *http.Request) (*getProfileResult, error) {
	userId, err := s.getUserIdFromSessionAndDb(w, r)
	if err != nil {
		return nil, fmt.Errorf("getUserIdFromSessionAndDb() failed: %v", err)
	}

	if len(userId) == 0 {
		// Not an error.
		// It's just not in the session cookie.
		return &getProfileResult{}, nil
	}

	c := r.Context()
	profile, err := s.userDataClient.GetUserProfileById(c, userId)
	if err != nil {
		return nil, fmt.Errorf("GetUserProfileById() failed: %v", err)
	}

	return &getProfileResult{
		Profile: profile,
		UserId:  userId,
	}, nil
}

type getLoginInfoResult struct {
	LoginInfo *restuser.LoginInfo
	UserId    string
}

// Returns the LoginInfo and the userID.
func (s *RestServer) getLoginInfoFromSessionAndDb(w http.ResponseWriter, r *http.Request) (*getLoginInfoResult, error) {
	var loginInfo restuser.LoginInfo

	getProfileResult, err := s.getProfileFromSessionAndDb(w, r)
	if err != nil {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = fmt.Sprintf("not logged in (%v)", err)
		return &getLoginInfoResult{
			LoginInfo: &loginInfo,
		}, nil
	}

	if getProfileResult == nil {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in (nil getProfileResult)"
		return &getLoginInfoResult{
			LoginInfo: &loginInfo,
		}, nil
	}

	s.updateLoginInfoFromProfile(&loginInfo, getProfileResult.Profile)
	return &getLoginInfoResult{
		LoginInfo: &loginInfo,
		UserId:    getProfileResult.UserId,
	}, err
}

func (s *RestServer) updateLoginInfoFromProfile(loginInfo *restuser.LoginInfo, profile *domainuser.Profile) {
	if profile == nil {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (no profile found)"
	} else {
		loginInfo.LoggedIn = true
		loginInfo.Nickname = profile.Name

		loginInfo.GoogleLinked = profile.GoogleLinked
		loginInfo.GoogleProfileUrl = profile.GoogleProfileUrl
		loginInfo.GitHubLinked = profile.GitHubLinked
		loginInfo.GitHubProfileUrl = profile.GitHubProfileUrl
		loginInfo.FacebookLinked = profile.FacebookLinked
		loginInfo.FacebookProfileUrl = profile.FacebookProfileUrl
	}
}

/** Get the user ID.
 * Returns an empty user ID, and a nil error, if the user is not logged in.
 */
func (s *RestServer) getUserIdFromSessionAndDb(w http.ResponseWriter, r *http.Request) (string, error) {
	userIdAndToken, err := s.userSessionStore.GetUserIdAndOAuthTokenFromSession(r)
	if err != nil {
		return "", fmt.Errorf("GetUserIdAndOAuthTokenFromSession() failed: %v", err)
	}

	if userIdAndToken == nil {
		return "", fmt.Errorf("GetUserIdAndOAuthTokenFromSession() returned nil")
	}

	token := userIdAndToken.Token
	if token == nil {
		return "", fmt.Errorf("GetUserIdAndOAuthTokenFromSession() returned a nil token")
	}

	if !token.Valid() {
		// TODO: Revalidate it.

		// This is not an error
		// (it is normal for a token to expire.)
		// (the user is now not logged in.)
		return "", nil
	}

	return userIdAndToken.UserId, nil
}
