package loginserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/murraycu/go-bigoquiz-server/config"
	"github.com/murraycu/go-bigoquiz-server/repositories/db"
	oauthparsers2 "github.com/murraycu/go-bigoquiz-server/server/loginserver/oauthparsers"
	"github.com/murraycu/go-bigoquiz-server/server/usersessionstore"
	"golang.org/x/oauth2"
)

type LoginServer struct {
	userDataClient   db.UserDataRepository
	oAuthStateClient *db.OAuthStateDataRepository

	// Session cookie store.
	userSessionStore usersessionstore.UserSessionStore

	confOAuthGoogle   *oauth2.Config
	confOAuthGitHub   *oauth2.Config
	confOAuthFacebook *oauth2.Config

	config *config.Config
}

func NewLoginServer(userSessionStore usersessionstore.UserSessionStore, conf *config.Config) (*LoginServer, error) {
	result := &LoginServer{}
	result.config = conf

	var err error
	result.userDataClient, err = db.NewUserDataRepository()
	if err != nil {
		return nil, fmt.Errorf("NewUserDataRepository() failed: %v", err)
	}

	result.oAuthStateClient, err = db.NewOAuthStateDataRepository()
	if err != nil {
		return nil, fmt.Errorf("NewOAuthStateDataRepository() failed: %v", err)
	}

	result.userSessionStore = userSessionStore

	result.confOAuthGoogle, err = config.GenerateGoogleOAuthConfig(conf)
	if err != nil {
		log.Fatalf("Unable to generate Google OAuth config: %v", err)
	}

	result.confOAuthGitHub, err = config.GenerateGitHubOAuthConfig(conf)
	if err != nil {
		log.Fatalf("Unable to generate GitHub OAuth config: %v", err)
	}

	result.confOAuthFacebook, err = config.GenerateFacebookOAuthConfig(conf)
	if err != nil {
		log.Fatalf("Unable to generate Facebook OAuth config: %v", err)
	}

	return result, nil
}

/** Get an oauth2 URL based on the oauth config.
 */
func (s *LoginServer) generateOAuthUrl(r *http.Request, oauthConfig *oauth2.Config) (string, error) {
	c := r.Context()

	state, err := s.generateState(c)
	if err != nil {
		return "", fmt.Errorf("unable to generate state: %v", err)
	}

	return oauthConfig.AuthCodeURL(state), nil
}

func (s *LoginServer) HandleGoogleLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.redirectToGoogleLogin(w, r)
}

func (s *LoginServer) redirectToGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect the user to the Google login page:
	url, err := s.generateOAuthUrl(r, s.confOAuthGoogle)
	if err != nil {
		loginStartFailedErr("generateOAuthUrl() failed", err, w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (s *LoginServer) generateState(c context.Context) (string, error) {
	state := rand.Int63()
	err := s.oAuthStateClient.StoreOAuthState(c, state)
	if err != nil {
		return "", fmt.Errorf("StoreOAuthState() failed: %v", err)
	}

	return strconv.FormatInt(state, 10), nil
}

func (s *LoginServer) checkState(c context.Context, state string) error {
	stateNum, err := strconv.ParseInt(state, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt() failed: %v", err)
	}

	err = s.oAuthStateClient.CheckOAuthState(c, stateNum)
	if err != nil {
		return fmt.Errorf("db.CheckOAuthState() failed: %v", err)
	}

	return nil
}

func (s *LoginServer) removeState(c context.Context, state string) error {
	stateNum, err := strconv.ParseInt(state, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt() failed: %v", err)
	}

	return s.oAuthStateClient.RemoveOAuthState(c, stateNum)
}

func (s *LoginServer) checkStateAndGetCode(c context.Context, r *http.Request) (string, error) {
	state := r.FormValue("state")
	err := s.checkState(c, state)
	if err != nil {
		return "", fmt.Errorf("invalid oauth state ('%s): %v", state, err)
	}

	// The state will not be used again,
	// so remove it from the datastore.
	err = s.removeState(c, state)
	if err != nil {
		return "", fmt.Errorf("removeState() failed: %v", err)
	}

	return r.FormValue("code"), nil
}

func (s *LoginServer) HandleGoogleCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := r.Context()

	token, body, ok := s.checkStateAndGetBody(w, r, s.confOAuthGoogle, "https://www.googleapis.com/oauth2/v3/userinfo", c)
	if !ok {
		// checkStateAndGetBody() already called loginFailed().
		return
	}

	var userinfo oauthparsers2.GoogleUserInfo
	err := json.Unmarshal(body, &userinfo)
	if err != nil {
		loginCallbackFailedErr("Unmarshaling of JSON from oauth2 callback failed", err, w, r)
		return
	}

	// Get the existing logged-in user's userId, if any, from the cookie, if any:
	userId, _, err := s.userSessionStore.GetProfileFromSession(r)
	if err != nil {
		loginCallbackFailedErr("getProfileFromSession() failed", err, w, r)
		return
	}
	// Store in the database,
	// either creating a new user or updating an existing user.
	userId, err = s.userDataClient.StoreGoogleLoginInUserProfile(c, userinfo, userId, token)
	if err != nil {
		loginCallbackFailedErr("StoreGoogleLoginInUserProfile() failed", err, w, r)
		return
	}

	s.storeCookieAndRedirect(r, w, c, userId, token)
}

func (s *LoginServer) checkStateAndGetBody(w http.ResponseWriter, r *http.Request, conf *oauth2.Config, url string, c context.Context) (*oauth2.Token, []byte, bool) {
	code, err := s.checkStateAndGetCode(c, r)
	if err != nil {
		log.Printf("checkStateAndGetCode() failed: %v", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil, nil, false
	}

	return s.exchangeAndGetUserBody(w, r, conf, code, url, c)
}

func (s *LoginServer) exchangeAndGetUserBody(w http.ResponseWriter, r *http.Request, conf *oauth2.Config, code string, url string, c context.Context) (*oauth2.Token, []byte, bool) {
	token, err := conf.Exchange(c, code)
	if err != nil {
		s.loginFailed("config.Exchange() failed", err, w, r)
		return nil, nil, false
	}

	if !token.Valid() {
		s.loginFailed("loginFailedUrl.Exchange() returned an invalid token", err, w, r)
		return nil, nil, false
	}

	client := conf.Client(c, token)
	infoResponse, err := client.Get(url)
	if err != nil {
		s.loginFailed("client.Get() failed", err, w, r)
		return nil, nil, false
	}

	defer func() {
		err := infoResponse.Body.Close()
		if err != nil {
			s.loginFailed("Body.Close() failed", err, w, r)
		}
	}()

	body, err := ioutil.ReadAll(infoResponse.Body)
	if err != nil {
		s.loginFailed("ReadAll(body) failed", err, w, r)
		return nil, nil, false
	}

	return token, body, true
}

// Called after user info has been successfully stored in the database.
func (s *LoginServer) storeCookieAndRedirect(r *http.Request, w http.ResponseWriter, c context.Context, strUserId string, token *oauth2.Token) {
	// Store the token in the cookie
	// so we can retrieve it from subsequent requests from the browser.
	session, err := s.userSessionStore.GetSession(r)
	if err != nil {
		s.loginFailed("Could not create new session", err, w, r)
		return
	}

	session.Values[usersessionstore.OAuthTokenSessionKey] = token
	session.Values[usersessionstore.UserIdSessionKey] = strUserId

	if err := session.Save(r, w); err != nil {
		s.loginFailed("Could not save session", err, w, r)
		return
	}

	// Redirect the user back to a page to show they are logged in:
	var userProfileUrl = s.config.BaseUrl + "/user"
	http.Redirect(w, r, userProfileUrl, http.StatusFound)
}

func (s *LoginServer) HandleLogout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Wipe the cookie:
	session, err := s.userSessionStore.GetSession(r)
	if err != nil {
		logoutError("could not get default session", err, w)
		return
	}

	session.Options.MaxAge = -1 // Clear session.

	if err := session.Save(r, w); err != nil {
		logoutError("Could not save session", err, w)
		return
	}

	redirectURL := r.FormValue("redirect")
	if redirectURL == "" {
		redirectURL = "/"
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

/** Get an oauth2 URL based on the secret .json file.
 * See githubConfigCredentialsFilename.
 */
func (s *LoginServer) generateGitHubOAuthUrl(r *http.Request) string {
	c := r.Context()

	state, err := s.generateState(c)
	if err != nil {
		log.Printf("Unable to generate state.")
		return ""
	}

	return s.confOAuthGitHub.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (s *LoginServer) HandleGitHubLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.redirectToGitHubLogin(w, r)
}

func (s *LoginServer) redirectToGitHubLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect the user to the GitHub login page:
	url, err := s.generateOAuthUrl(r, s.confOAuthGitHub)
	if err != nil {
		loginStartFailedErr("generateOAuthUrl() failed", err, w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (s *LoginServer) HandleGitHubCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := r.Context()

	token, body, ok := s.checkStateAndGetBody(w, r, s.confOAuthGitHub, "https://api.github.com/user", c)
	if !ok {
		// checkStateAndGetBody() already called loginFailed().
		return
	}

	var userinfo oauthparsers2.GitHubUserInfo
	err := json.Unmarshal(body, &userinfo)
	if err != nil {
		loginCallbackFailedErr("Unmarshaling of JSON from oauth2 callback failed", err, w, r)
		return
	}

	// Get the existing logged-in user's userId, if any, from the cookie, if any:
	userId, _, err := s.userSessionStore.GetProfileFromSession(r)
	if err != nil {
		loginCallbackFailedErr("getProfileFromSession() failed", err, w, r)
		return
	}

	userId, err = s.userDataClient.StoreGitHubLoginInUserProfile(c, userinfo, userId, token)
	if err != nil {
		loginCallbackFailedErr("StoreGitHubLoginInUserProfile() failed", err, w, r)
		return
	}

	s.storeCookieAndRedirect(r, w, c, userId, token)
}

func (s *LoginServer) HandleFacebookLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.redirectToFacebookLogin(w, r)
}

func (s *LoginServer) redirectToFacebookLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect the user to the Facebook login page:
	url, err := s.generateOAuthUrl(r, s.confOAuthFacebook)
	if err != nil {
		loginStartFailedErr("generateOAuthUrl() failed", err, w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (s *LoginServer) HandleFacebookCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := r.Context()

	token, body, ok := s.checkStateAndGetBody(w, r, s.confOAuthFacebook, "https://graph.facebook.com/me?fields=link,name,email", c)
	if !ok {
		// checkStateAndGetBody() already called loginFailed().
		return
	}

	var userinfo oauthparsers2.FacebookUserInfo
	err := json.Unmarshal(body, &userinfo)
	if err != nil {
		loginCallbackFailedErr("Unmarshaling of JSON from oauth2 callback failed", err, w, r)
		return
	}

	// Get the existing logged-in user's userId, if any, from the cookie, if any:
	// (This lets us associate multiple oauth2 logins with a single user ID.)
	userId, _, err := s.userSessionStore.GetProfileFromSession(r)
	if err != nil {
		loginCallbackFailedErr("getProfileFromSession() failed.", err, w, r)
		return
	}

	// Store in the database:
	userId, err = s.userDataClient.StoreFacebookLoginInUserProfile(c, userinfo, userId, token)
	if err != nil {
		loginCallbackFailedErr("StoreFacebookLoginInUserProfile() failed.", err, w, r)
		return
	}

	s.storeCookieAndRedirect(r, w, c, userId, token)
}

func (s *LoginServer) loginFailed(message string, err error, w http.ResponseWriter, r *http.Request) {
	var loginFailedUrl = s.config.BaseUrl + "/login?failed=true"

	log.Printf(message+":'%v'\n", err)
	http.Redirect(w, r, loginFailedUrl, http.StatusTemporaryRedirect)
}

func loginStartFailedErr(message string, err error, w http.ResponseWriter, r *http.Request) {
	log.Printf(message+":'%v'\n", err)
	http.Redirect(w, r, "/", http.StatusInternalServerError)
}

func loginCallbackFailedErr(message string, err error, w http.ResponseWriter, r *http.Request) {
	log.Printf(message+":'%v'\n", err)
	http.Redirect(w, r, "/", http.StatusInternalServerError)
}

func logoutError(message string, err error, w http.ResponseWriter) {
	handleErrorAsHttpError(w, http.StatusInternalServerError, "messsage: %v", err)
}

func handleErrorAsHttpError(w http.ResponseWriter, code int, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	log.Print(msg)

	http.Error(w, msg, code)
}
