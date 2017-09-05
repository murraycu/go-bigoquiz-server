package bigoquiz

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"io/ioutil"
	"net/http"
)

/** Get an oauth2 URL based on the secret .json file.
 * See credentialsFilename.
 */
func generateGoogleOAuthUrl(r *http.Request) string {
	c := appengine.NewContext(r)

	config := generateGoogleOAuthConfig(r)
	if config == nil {
		log.Errorf(c, "Unable to generate config.")
		return ""
	}

	return config.AuthCodeURL(oauthStateString)
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Redirect the user to the Google login page:
	url := generateGoogleOAuthUrl(r)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
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

	config := generateGoogleOAuthConfig(r)
	if config == nil {
		log.Errorf(c, "Unable to generate config.")
		return
	}

	token, err := config.Exchange(c, code)
	if err != nil {
		loginFailed(c, "config.Exchange() failed", err, w, r)
		return
	}

	if !token.Valid() {
		loginFailed(c, "loginFailedUrl.Exchange() returned an invalid token", err, w, r)
		return
	}

	client := config.Client(c, token)
	infoResponse, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		loginFailed(c, "client.Get() failed", err, w, r)
		return
	}

	session, err := store.New(r, defaultSessionID)
	if err != nil {
		loginFailed(c, "Could not create new session", err, w, r)
	}

	// Store the token in the cookie
	// so we can retrieve it from subsequent requests from the browser.
	session.Values[oauthTokenSessionKey] = token

	// Just to show that it worked:
	defer infoResponse.Body.Close()
	body, err := ioutil.ReadAll(infoResponse.Body)
	if err != nil {
		loginFailed(c, "ReadAll(body) failed", err, w, r)
		return
	}

	var userinfo GoogleUserInfo
	json.Unmarshal(body, &userinfo)

	/*
	session, err := store.Get(r, "sess")
	if err != nil {
		log.Errorf(c, "Failed to get session from store.", err)
		http.Redirect(w, r, baseUrl, http.StatusTemporaryRedirect)
		return
	}
	*/

	session.Values["name"] = userinfo.Name
	session.Values["accessToken"] = token.AccessToken
	session.Save(r, w)

	var userProfileUrl = baseUrl + "/user"
	http.Redirect(w, r, userProfileUrl, http.StatusTemporaryRedirect)
}

func loginFailed(c context.Context, message string, err error, w http.ResponseWriter, r *http.Request) {
	var loginFailedUrl = baseUrl + "/login?failed=true"

	log.Errorf(c, message + ":'%v'\n", err)
	http.Redirect(w, r, loginFailedUrl, http.StatusTemporaryRedirect)
}

var (
	// Some random string, random for each request
	// TODO: Actually be random, and somehow check it in the callback.
	oauthStateString = "random"

	// This file must be downloaded
	// (via the "DOWNLOAD JSON" link at https://console.developers.google.com/apis/credentials/oauthclient )
	// and added with this exact filename, next to this .go source file.
	credentialsFilename = "google_oauth2_credentials_secret.json"

	// See https://developers.google.com/+/web/api/rest/oauth#profile
	credentialsScopeProfile = "profile"

	// See https://developers.google.com/identity/protocols/googlescopes
	credentialsScopeEmail   = "https://www.googleapis.com/auth/userinfo.email"

	// We store the token in a session cookie.
	store *sessions.CookieStore
	defaultSessionID = "default"
	oauthTokenSessionKey = "oauth_token"

    baseUrl = "http://beta.bigoquiz.com"
)

/** Get an oauth2 Config object based on the secret .json file.
 * See credentialsFilename.
 */
func generateGoogleOAuthConfig(r *http.Request) *oauth2.Config {
	c := appengine.NewContext(r)

	b, err := ioutil.ReadFile(credentialsFilename)
	if err != nil {
		log.Errorf(c, "Unable to read client secret file (%s): %v", credentialsFilename, err)
	}

	config, err := google.ConfigFromJSON(b, credentialsScopeProfile, credentialsScopeEmail)
	if err != nil {
		log.Errorf(c, "Unable to parse client secret file (%) to config: %v", credentialsFilename, err)
	}

	return config
}
