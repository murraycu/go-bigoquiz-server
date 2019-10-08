package loginserver

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/murraycu/go-bigoquiz-server/config"
	"github.com/murraycu/go-bigoquiz-server/repositories/db"
	"github.com/murraycu/go-bigoquiz-server/usersessionstore"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type LoginServer struct {
	UserDataClient *db.UserDataRepository

	// Session cookie store.
	UserSessionStore *usersessionstore.UserSessionStore

	ConfOAuthGoogle   *oauth2.Config
	ConfOAuthGitHub   *oauth2.Config
	ConfOAuthFacebook *oauth2.Config
}

func NewLoginServer(userSessionStore *usersessionstore.UserSessionStore) (*LoginServer, error) {
	result := &LoginServer{}

	var err error
	result.UserDataClient, err = db.NewUserDataRepository()
	if err != nil {
		return nil, fmt.Errorf("NewUserDataRepository() failed: %v", err)
	}

	result.UserSessionStore = userSessionStore

	result.ConfOAuthGoogle, err = config.GenerateGoogleOAuthConfig()
	if err != nil {
		log.Fatalf("Unable to generate Google OAuth config: %v", err)
	}

	result.ConfOAuthGitHub, err = config.GenerateGitHubOAuthConfig()
	if err != nil {
		log.Fatalf("Unable to generate GitHub OAuth config: %v", err)
	}

	result.ConfOAuthFacebook, err = config.GenerateFacebookOAuthConfig()
	if err != nil {
		log.Fatalf("Unable to generate Facebook OAuth config: %v", err)
	}

	return result, nil
}

/** Get an oauth2 URL based on the oauth config.
 */
func (s *LoginServer) generateOAuthUrl(r *http.Request, oauthConfig *oauth2.Config) string {
	c := r.Context()

	state, err := generateState(c)
	if err != nil {
		log.Printf("Unable to generate state: %v", err)
		return ""
	}

	return oauthConfig.AuthCodeURL(state)
}

func (s *LoginServer) HandleGoogleLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Redirect the user to the Google login page:
	url := s.generateOAuthUrl(r, s.ConfOAuthGoogle)
	http.Redirect(w, r, url, http.StatusFound)
}

func generateState(c context.Context) (string, error) {
	dbClient, err := db.NewOAuthStateDataRepository()
	if err != nil {
		return "", fmt.Errorf("NewOAuthStateDataRepository() failed: %v", err)
	}

	state := rand.Int63()
	err = dbClient.StoreOAuthState(c, state)
	if err != nil {
		return "", fmt.Errorf("StoreOAuthState() failed: %v", err)
	}

	return strconv.FormatInt(state, 10), nil
}

func checkState(c context.Context, state string) error {
	stateNum, err := strconv.ParseInt(state, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt() failed: %v", err)
	}

	dbClient, err := db.NewOAuthStateDataRepository()
	if err != nil {
		return fmt.Errorf("NewUserDataRepository() failed: %v", err)
	}

	err = dbClient.CheckOAuthState(c, stateNum)
	if err != nil {
		return fmt.Errorf("db.CheckOAuthState() failed: %v", err)
	}

	return nil
}

func removeState(c context.Context, state string) error {
	stateNum, err := strconv.ParseInt(state, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt() failed: %v", err)
	}

	dbClient, err := db.NewOAuthStateDataRepository()
	if err != nil {
		return fmt.Errorf("NewUserDataRepository() failed: %v", err)
	}

	return dbClient.RemoveOAuthState(c, stateNum)
}

func checkStateAndGetCode(c context.Context, r *http.Request) (string, error) {
	state := r.FormValue("state")
	err := checkState(c, state)
	if err != nil {
		return "", fmt.Errorf("invalid oauth state ('%s): %v", state, err)
	}

	// The state will not be used again,
	// so remove it from the datastore.
	err = removeState(c, state)
	if err != nil {
		return "", fmt.Errorf("removeState() failed: %v", err)
	}

	return r.FormValue("code"), nil
}

func (s *LoginServer) HandleGoogleCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := r.Context()

	token, body, ok := checkStateAndGetBody(w, r, s.ConfOAuthGoogle, "https://www.googleapis.com/oauth2/v3/userinfo", c)
	if !ok {
		// checkStateAndGetBody() already called loginFailed().
		return
	}

	var userinfo db.GoogleUserInfo
	err := json.Unmarshal(body, &userinfo)
	if err != nil {
		loginCallbackFailedErr("Unmarshaling of JSON from oauth2 callback failed", err, w, r)
		return
	}

	// Get the existing logged-in user's userId, if any, from the cookie, if any:
	userId, _, err := s.UserSessionStore.GetProfileFromSession(r)
	if err != nil {
		loginCallbackFailedErr("getProfileFromSession() failed", err, w, r)
		return
	}

	// Store in the database:
	dbClient, err := db.NewUserDataRepository()
	if err != nil {
		loginCallbackFailedErr("dbClient.NewUserDataRepository() failed", err, w, r)
		return
	}

	// Store in the database,
	// either creating a new user or updating an existing user.
	userId, err = dbClient.StoreGoogleLoginInUserProfile(c, userinfo, userId, token)
	if err != nil {
		loginCallbackFailedErr("StoreGoogleLoginInUserProfile() failed", err, w, r)
		return
	}

	s.storeCookieAndRedirect(r, w, c, userId, token)
}

func checkStateAndGetBody(w http.ResponseWriter, r *http.Request, conf *oauth2.Config, url string, c context.Context) (*oauth2.Token, []byte, bool) {
	code, err := checkStateAndGetCode(c, r)
	if err != nil {
		log.Printf("checkStateAndGetCode() failed: %v", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil, nil, false
	}

	return exchangeAndGetUserBody(w, r, conf, code, url, c)
}

func exchangeAndGetUserBody(w http.ResponseWriter, r *http.Request, conf *oauth2.Config, code string, url string, c context.Context) (*oauth2.Token, []byte, bool) {
	token, err := conf.Exchange(c, code)
	if err != nil {
		loginFailed("config.Exchange() failed", err, w, r)
		return nil, nil, false
	}

	if !token.Valid() {
		loginFailed("loginFailedUrl.Exchange() returned an invalid token", err, w, r)
		return nil, nil, false
	}

	client := conf.Client(c, token)
	infoResponse, err := client.Get(url)
	if err != nil {
		loginFailed("client.Get() failed", err, w, r)
		return nil, nil, false
	}

	defer infoResponse.Body.Close()
	body, err := ioutil.ReadAll(infoResponse.Body)
	if err != nil {
		loginFailed("ReadAll(body) failed", err, w, r)
		return nil, nil, false
	}

	return token, body, true
}

// Called after user info has been successful stored in the database.
func (s *LoginServer) storeCookieAndRedirect(r *http.Request, w http.ResponseWriter, c context.Context, strUserId string, token *oauth2.Token) {
	// Store the token in the cookie
	// so we can retrieve it from subsequent requests from the browser.
	session, err := s.UserSessionStore.Store.New(r, usersessionstore.DefaultSessionID)
	if err != nil {
		loginFailed("Could not create new session", err, w, r)
		return
	}

	session.Values[usersessionstore.OAuthTokenSessionKey] = token
	session.Values[usersessionstore.UserIdSessionKey] = strUserId

	if err := session.Save(r, w); err != nil {
		loginFailed("Could not save session", err, w, r)
		return
	}

	// Redirect the user back to a page to show they are logged in:
	var userProfileUrl = config.BaseUrl + "/user"
	http.Redirect(w, r, userProfileUrl, http.StatusFound)
}

func (s *LoginServer) HandleLogout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Wipe the cookie:
	session, err := s.UserSessionStore.Store.New(r, usersessionstore.DefaultSessionID)
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

	state, err := generateState(c)
	if err != nil {
		log.Printf("Unable to generate state.")
		return ""
	}

	return s.ConfOAuthGitHub.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (s *LoginServer) HandleGitHubLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Redirect the user to the GitHub login page:
	url := s.generateOAuthUrl(r, s.ConfOAuthGitHub)
	http.Redirect(w, r, url, http.StatusFound)
}

func (s *LoginServer) HandleGitHubCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := r.Context()

	token, body, ok := checkStateAndGetBody(w, r, s.ConfOAuthGitHub, "https://api.github.com/user", c)
	if !ok {
		// checkStateAndGetBody() already called loginFailed().
		return
	}

	var userinfo db.GitHubUserInfo
	err := json.Unmarshal(body, &userinfo)
	if err != nil {
		loginCallbackFailedErr("Unmarshaling of JSON from oauth2 callback failed", err, w, r)
		return
	}

	// Get the existing logged-in user's userId, if any, from the cookie, if any:
	userId, _, err := s.UserSessionStore.GetProfileFromSession(r)
	if err != nil {
		loginCallbackFailedErr("getProfileFromSession() failed", err, w, r)
		return
	}

	// Store in the database:
	dbClient, err := db.NewUserDataRepository()
	if err != nil {
		loginCallbackFailedErr("dbClient.NewUserDataRepository() failed", err, w, r)
		return
	}

	userId, err = dbClient.StoreGitHubLoginInUserProfile(c, userinfo, userId, token)
	if err != nil {
		loginCallbackFailedErr("StoreGitHubLoginInUserProfile() failed", err, w, r)
		return
	}

	s.storeCookieAndRedirect(r, w, c, userId, token)
}

func (s *LoginServer) HandleFacebookLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Redirect the user to the Facebook login page:
	url := s.generateOAuthUrl(r, s.ConfOAuthFacebook)
	http.Redirect(w, r, url, http.StatusFound)
}

func (s *LoginServer) HandleFacebookCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := r.Context()

	token, body, ok := checkStateAndGetBody(w, r, s.ConfOAuthFacebook, "https://graph.facebook.com/me?fields=link,name,email", c)
	if !ok {
		// checkStateAndGetBody() already called loginFailed().
		return
	}

	var userinfo db.FacebookUserInfo
	err := json.Unmarshal(body, &userinfo)
	if err != nil {
		loginCallbackFailedErr("Unmarshaling of JSON from oauth2 callback failed", err, w, r)
		return
	}

	// Get the existing logged-in user's userId, if any, from the cookie, if any:
	userId, _, err := s.UserSessionStore.GetProfileFromSession(r)
	if err != nil {
		loginCallbackFailedErr("getProfileFromSession() failed.", err, w, r)
		return
	}

	// Store in the database:
	dbClient, err := db.NewUserDataRepository()
	if err != nil {
		loginCallbackFailedErr("dbClient.NewUserDataRepository() failed", err, w, r)
		return
	}

	// Store in the database:
	userId, err = dbClient.StoreFacebookLoginInUserProfile(c, userinfo, userId, token)
	if err != nil {
		loginCallbackFailedErr("StoreFacebookLoginInUserProfile() failed.", err, w, r)
		return
	}

	s.storeCookieAndRedirect(r, w, c, userId, token)
}

func loginFailed(message string, err error, w http.ResponseWriter, r *http.Request) {
	var loginFailedUrl = config.BaseUrl + "/login?failed=true"

	log.Printf(message+":'%v'\n", err)
	http.Redirect(w, r, loginFailedUrl, http.StatusTemporaryRedirect)
}

func loginCallbackFailedErr(message string, err error, w http.ResponseWriter, r *http.Request) {
	log.Printf(message+":'%v'\n", err)
	http.Redirect(w, r, "/", http.StatusInternalServerError)
}

func logoutError(message string, err error, w http.ResponseWriter) {
	log.Printf(message+":'%v'\n", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
