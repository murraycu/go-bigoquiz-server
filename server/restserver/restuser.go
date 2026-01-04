package restserver

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	domainuser "github.com/murraycu/go-bigoquiz-server/domain/user"
	restuser "github.com/murraycu/go-bigoquiz-server/server/restserver/user"
)

func (s *RestServer) HandleUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	result, err := s.getLoginInfoFromSessionAndDb(r)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "getLoginInfoFromSessionAndDb() failed: %v", err)
		return
	}

	marshalAndWriteOrHttpError(w, result.LoginInfo)
}

type GetProfileResult struct {
	Profile      *domainuser.Profile
	UserId       string
	InvalidToken bool
}

// Returns the LoginInfo and the userID.
func (s *RestServer) getProfileFromSessionAndDb(r *http.Request) (*GetProfileResult, error) {
	userIdResult, err := s.getUserIdAndTokenValidityFromSessionAndDb(r)
	if err != nil {
		return nil, fmt.Errorf("getUserIdFromSessionAndDb() failed: %v", err)
	}

	userId := userIdResult.UserId

	if len(userId) == 0 {
		// Not an error.
		// It's just not in the session cookie.
		return nil, nil
	}

	c := r.Context()
	profile, err := s.userDataClient.GetUserProfileById(c, userId)
	if err != nil {
		return nil, fmt.Errorf("GetUserProfileById() failed: %v", err)
	}

	return &GetProfileResult{
		Profile:      profile,
		UserId:       userId,
		InvalidToken: userIdResult.InvalidToken,
	}, nil
}

type GetLoginInfoResult struct {
	LoginInfo *restuser.LoginInfo
	UserId    string
}

// Returns the LoginInfo and the userID.
func (s *RestServer) getLoginInfoFromSessionAndDb(r *http.Request) (*GetLoginInfoResult, error) {
	var loginInfo restuser.LoginInfo

	getProfileResult, err := s.getProfileFromSessionAndDb(r)
	if err != nil {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = fmt.Sprintf("not logged in (%v)", err)
	}

	s.updateLoginInfoFromProfile(&loginInfo, getProfileResult.Profile, getProfileResult.InvalidToken)

	return &GetLoginInfoResult{
		LoginInfo: &loginInfo,
		UserId:    getProfileResult.UserId,
	}, err
}

func (s *RestServer) updateLoginInfoFromProfile(loginInfo *restuser.LoginInfo, profile *domainuser.Profile, invalidToken bool) {
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

	if invalidToken {
		if loginInfo.GoogleLinked {
			loginInfo.GoogleTokenExpired = true
		}

		if loginInfo.GitHubLinked {
			loginInfo.GitHubTokenExpired = true
		}

		if loginInfo.FacebookLinked {
			loginInfo.FacebookTokenExpired = true
		}
	}
}

/** getUserIdFromSessionAndDb() get the user ID (without checking if the OAuth token is valid).
 * Returns an empty user ID, and a nil error, if the user is not logged in.
 */
func (s *RestServer) getUserIdFromSessionAndDb(r *http.Request) (string, error) {
	result, err := s.getUserIdAndTokenValidityFromSessionAndDb(r)
	if err != nil {
		return "", fmt.Errorf("getUserIdAndTokenValidityFromSessionAndDb() failed: %v", err)
	}

	if result.InvalidToken {
		return "", nil
	}

	return result.UserId, nil
}

type UserIdResult struct {
	UserId       string
	InvalidToken bool
}

/** getUserIdAndTokenValidityFromSessionAndDb gets the User ID and gets whether the OAuth token is valid.
 * Returns an empty user ID, and a nil error, if the user is not logged in.
 */
func (s *RestServer) getUserIdAndTokenValidityFromSessionAndDb(r *http.Request) (*UserIdResult, error) {
	userId, token, err := s.userSessionStore.GetProfileFromSession(r)
	if err != nil {
		return nil, fmt.Errorf("GetProfileFromSession() failed: %v", err)
	}

	if !token.Valid() {
		// This is not an error
		// (it is normal for a token to expire.)
		// (the user is now not logged in.)
		return &UserIdResult{
			UserId:       userId,
			InvalidToken: true,
		}, nil
	}

	return &UserIdResult{
		UserId: userId,
	}, nil
}
