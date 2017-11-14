package bigoquiz

import (
	"config"
	"db"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"io/ioutil"
	"net/http"
)

/** Get an oauth2 URL based on the secret .json file.
 * See googleConfigCredentialsFilename.
 */
func generateGoogleOAuthUrl(r *http.Request) string {
	c := appengine.NewContext(r)

	conf := config.GenerateGoogleOAuthConfig(r)
	if conf == nil {
		log.Errorf(c, "Unable to generate config.")
		return ""
	}

	return conf.AuthCodeURL(oauthStateString)
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Redirect the user to the Google login page:
	url := generateGoogleOAuthUrl(r)
	http.Redirect(w, r, url, http.StatusFound)
}

func checkStateAndGetCode(r *http.Request) (string, error) {
	state := r.FormValue("state")
	if state != oauthStateString {
		return "", fmt.Errorf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
	}

	return r.FormValue("code"), nil
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := appengine.NewContext(r)

	code, err := checkStateAndGetCode(r)
	if err != nil {
		log.Errorf(c, "checkStateAndGetCode() failed", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	conf := config.GenerateGoogleOAuthConfig(r)
	if conf == nil {
		log.Errorf(c, "Unable to generate config.")
		return
	}

	token, body, ok := exchangeAndGetUserBody(w, r, conf, code, "https://www.googleapis.com/oauth2/v3/userinfo", c)
	if !ok {
		// exchangeAndGetUserBody() already called loginFailed().
		return
	}

	var userinfo db.GoogleUserInfo
	err = json.Unmarshal(body, &userinfo)
	if err != nil {
		log.Errorf(c, "Unmarshaling of JSON from oauth2 callback failed:'%v'\n", err)
		http.Error(w, "Unmarshaling of JSON from oauth2 callback failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the existing logged-in user's userId, if any, from the cookie, if any:
	userId, _, err := getProfileFromSession(r)
	if err != nil {
		log.Errorf(c, "getProfileFromSession() failed:'%v'\n", err)
		http.Error(w, "getProfileFromSession() failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Store in the database,
	// either creating a new user or updating an existing user.
	userId, err = db.StoreGoogleLoginInUserProfile(c, userinfo, userId, token)
	if err != nil {
		log.Errorf(c, "StoreGoogleLoginInUserProfile() failed:'%v'\n", err)
		http.Error(w, "StoreGoogleLoginInUserProfile() failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	storeCookieAndRedirect(r, w, c, userId, token)
}

func exchangeAndGetUserBody(w http.ResponseWriter, r *http.Request, conf *oauth2.Config, code string, url string, c context.Context) (*oauth2.Token, []byte, bool) {
	token, err := conf.Exchange(c, code)
	if err != nil {
		loginFailed(c, "config.Exchange() failed", err, w, r)
		return nil, nil, false
	}

	if !token.Valid() {
		loginFailed(c, "loginFailedUrl.Exchange() returned an invalid token", err, w, r)
		return nil, nil, false
	}

	client := conf.Client(c, token)
	infoResponse, err := client.Get(url)
	if err != nil {
		loginFailed(c, "client.Get() failed", err, w, r)
		return nil, nil, false
	}

	defer infoResponse.Body.Close()
	body, err := ioutil.ReadAll(infoResponse.Body)
	if err != nil {
		loginFailed(c, "ReadAll(body) failed", err, w, r)
		return nil, nil, false
	}

	return token, body, true
}

// Called after user info has been successful stored in the database.
func storeCookieAndRedirect(r *http.Request, w http.ResponseWriter, c context.Context, userId *datastore.Key, token *oauth2.Token) {
	// Store the token in the cookie
	// so we can retrieve it from subsequent requests from the browser.
	session, err := store.New(r, defaultSessionID)
	if err != nil {
		loginFailed(c, "Could not create new session", err, w, r)
		return
	}

	session.Values[oauthTokenSessionKey] = token
	session.Values[userIdSessionKey] = userId

	if err := session.Save(r, w); err != nil {
		loginFailed(c, "Could not save session", err, w, r)
		return
	}

	// Redirect the user back to a page to show they are logged in:
	var userProfileUrl = config.BaseUrl + "/user"
	http.Redirect(w, r, userProfileUrl, http.StatusFound)
}

func handleLogout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Wipe the cookie:
	session, err := store.New(r, defaultSessionID)
	if err != nil {
		logoutError("could not get default session", err, w, r)
		return
	}

	session.Options.MaxAge = -1 // Clear session.

	if err := session.Save(r, w); err != nil {
		logoutError("Could not save session", err, w, r)
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
func generateGitHubOAuthUrl(r *http.Request) string {
	c := appengine.NewContext(r)

	conf := config.GenerateGitHubOAuthConfig(r)
	if conf == nil {
		log.Errorf(c, "Unable to generate config.")
		return ""
	}

	return conf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
}

func handleGitHubLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Redirect the user to the GitHub login page:
	url := generateGitHubOAuthUrl(r)
	http.Redirect(w, r, url, http.StatusFound)
}

func handleGitHubCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := appengine.NewContext(r)

	code, err := checkStateAndGetCode(r)
	if err != nil {
		log.Errorf(c, "checkStateAndGetCode() failed", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	conf := config.GenerateGitHubOAuthConfig(r)
	if conf == nil {
		log.Errorf(c, "Unable to generate config.")
		return
	}

	token, body, ok := exchangeAndGetUserBody(w, r, conf, code, "https://api.github.com/user", c)
	if !ok {
		// exchangeAndGetUserBody() already called loginFailed().
		return
	}

	var userinfo db.GitHubUserInfo
	err = json.Unmarshal(body, &userinfo)
	if err != nil {
		log.Errorf(c, "Unmarshaling of JSON from oauth2 callback failed:'%v'\n", err)

		msg := fmt.Sprintf("Unmarshaling of JSON from oauth2 callback failed: %s\nfor:\n%v", err.Error(), string(body))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Get the existing logged-in user's userId, if any, from the cookie, if any:
	userId, _, err := getProfileFromSession(r)
	if err != nil {
		log.Errorf(c, "getProfileFromSession() failed:'%v'\n", err)
		http.Error(w, "getProfileFromSession() failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Store in the database:
	userId, err = db.StoreGitHubLoginInUserProfile(c, userinfo, userId, token)
	if err != nil {
		log.Errorf(c, "StoreGitHubLoginInUserProfile() failed:'%v'\n", err)
		http.Error(w, "StoreGitHubLoginInUserProfile() failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	storeCookieAndRedirect(r, w, c, userId, token)
}

/** Get an oauth2 URL based on the secret .json file.
 * See githubConfigCredentialsFilename.
 */
func generateFacebookOAuthUrl(r *http.Request) string {
	c := appengine.NewContext(r)

	conf := config.GenerateFacebookOAuthConfig(r)
	if conf == nil {
		log.Errorf(c, "Unable to generate config.")
		return ""
	}

	return conf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
}

func handleFacebookLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Redirect the user to the GitHub login page:
	url := generateFacebookOAuthUrl(r)
	http.Redirect(w, r, url, http.StatusFound)
}

func handleFacebookCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := appengine.NewContext(r)

	code, err := checkStateAndGetCode(r)
	if err != nil {
		log.Errorf(c, "checkStateAndGetCode() failed", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	conf := config.GenerateFacebookOAuthConfig(r)
	if conf == nil {
		log.Errorf(c, "Unable to generate config.")
		return
	}

	token, body, ok := exchangeAndGetUserBody(w, r, conf, code, "https://graph.facebook.com/me?fields=link,name,email", c)
	if !ok {
		// exchangeAndGetUserBody() already called loginFailed().
		return
	}

	var userinfo db.FacebookUserInfo
	err = json.Unmarshal(body, &userinfo)
	if err != nil {
		log.Errorf(c, "Unmarshaling of JSON from oauth2 callback failed:'%v'\n", err)

		msg := fmt.Sprintf("Unmarshaling of JSON from oauth2 callback failed: %s\nfor:\n%v", err.Error(), string(body))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Get the existing logged-in user's userId, if any, from the cookie, if any:
	userId, _, err := getProfileFromSession(r)
	if err != nil {
		log.Errorf(c, "getProfileFromSession() failed:'%v'\n", err)
		http.Error(w, "getProfileFromSession() failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Store in the database:
	userId, err = db.StoreFacebookLoginInUserProfile(c, userinfo, userId, token)
	if err != nil {
		log.Errorf(c, "StoreFacebookLoginInUserProfile() failed:'%v'\n", err)
		http.Error(w, "StoreFacebookLoginInUserProfile() failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	storeCookieAndRedirect(r, w, c, userId, token)
}

func loginFailed(c context.Context, message string, err error, w http.ResponseWriter, r *http.Request) {
	var loginFailedUrl = config.BaseUrl + "/login?failed=true"

	log.Errorf(c, message+":'%v'\n", err)
	http.Redirect(w, r, loginFailedUrl, http.StatusTemporaryRedirect)
}

func logoutError(message string, err error, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	log.Errorf(c, message+":'%v'\n", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

const (
	// Some random string, random for each request
	// TODO: Actually be random, and somehow check it in the callback.
	oauthStateString = "random"

	defaultSessionID     = "default"
	oauthTokenSessionKey = "oauth_token"
	userIdSessionKey     = "id" // A generic user ID, not a google user ID.
)

// We store the token in a session cookie.
var store *sessions.CookieStore
