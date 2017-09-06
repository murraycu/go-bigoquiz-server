package bigoquiz

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"io/ioutil"
	"net/http"
	"config"
)

/** Get an oauth2 URL based on the secret .json file.
 * See configCredentialsFilename.
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

func handleGoogleCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := appengine.NewContext(r)

	state := r.FormValue("state")
	if state != oauthStateString {
		log.Errorf(c, "invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	conf := config.GenerateGoogleOAuthConfig(r)
	if conf == nil {
		log.Errorf(c, "Unable to generate config.")
		return
	}

	token, err := conf.Exchange(c, code)
	if err != nil {
		loginFailed(c, "config.Exchange() failed", err, w, r)
		return
	}

	if !token.Valid() {
		loginFailed(c, "loginFailedUrl.Exchange() returned an invalid token", err, w, r)
		return
	}

	client := conf.Client(c, token)
	infoResponse, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		loginFailed(c, "client.Get() failed", err, w, r)
		return
	}

	// Just to show that it worked:
	defer infoResponse.Body.Close()
	body, err := ioutil.ReadAll(infoResponse.Body)
	if err != nil {
		loginFailed(c, "ReadAll(body) failed", err, w, r)
		return
	}

	var userinfo GoogleUserInfo
	json.Unmarshal(body, &userinfo)


	session, err := store.New(r, defaultSessionID)
	if err != nil {
		loginFailed(c, "Could not create new session", err, w, r)
	}

	// Store the token in the cookie
	// so we can retrieve it from subsequent requests from the browser.
	session.Values[oauthTokenSessionKey] = token
	session.Values[nameSessionKey] = &(userinfo.Name)
	if err := session.Save(r, w); err != nil {
		loginFailed(c, "Could not save session", err, w, r)
	}

	var userProfileUrl = baseUrl + "/user"
	http.Redirect(w, r, userProfileUrl, http.StatusFound)
}

func handleGoogleLogout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

func loginFailed(c context.Context, message string, err error, w http.ResponseWriter, r *http.Request) {
	var loginFailedUrl = baseUrl + "/login?failed=true"

	log.Errorf(c, message + ":'%v'\n", err)
	http.Redirect(w, r, loginFailedUrl, http.StatusTemporaryRedirect)
}

func logoutError(message string, err error, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	log.Errorf(c, message + ":'%v'\n", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

const (
	// Some random string, random for each request
	// TODO: Actually be random, and somehow check it in the callback.
	oauthStateString = "random"

	defaultSessionID = "default"
	oauthTokenSessionKey = "oauth_token"
	nameSessionKey = "name"

	baseUrl = "http://beta.bigoquiz.com"
)

// We store the token in a session cookie.
var store *sessions.CookieStore