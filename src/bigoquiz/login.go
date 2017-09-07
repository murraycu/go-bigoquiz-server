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
	"google.golang.org/appengine/datastore"
	"user"
	"golang.org/x/oauth2"
	"fmt"
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

	// Store in the database:
	userId, err := storeGoogleLoginInUserProfile(c, userinfo, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return;
	}

	// Store the token in the cookie
	// so we can retrieve it from subsequent requests from the browser.
	session, err := store.New(r, defaultSessionID)
	if err != nil {
		loginFailed(c, "Could not create new session", err, w, r)
	}

	session.Values[oauthTokenSessionKey] = token
	session.Values[userIdSessionKey] = userId

	if err := session.Save(r, w); err != nil {
		loginFailed(c, "Could not save session", err, w, r)
	}

	// Redirect the user back to a page to show they are logged in:
	var userProfileUrl = config.BaseUrl + "/user"
	http.Redirect(w, r, userProfileUrl, http.StatusFound)
}

// Get the UserProfile via the GoogleID, adding it if necessary.
func storeGoogleLoginInUserProfile(c context.Context, userInfo GoogleUserInfo, token *oauth2.Token) (*datastore.Key, error) {
	q := datastore.NewQuery("user.Profile").
		Filter("GoogleId =", userInfo.Sub).
		Limit(1)
	iter := q.Run(c)
	if iter == nil {
		return nil, fmt.Errorf("datastore query for GoogleId failed.")
	}

	var profile user.Profile
	var key *datastore.Key
	var err error
	key, err = iter.Next(&profile)
	if err == datastore.Done {
		// It is not in the datastore yet, so we add it.
		updateProfileFromGoogleUserInfo(&profile, &userInfo)
		profile.GoogleAccessToken = *token

		key = datastore.NewIncompleteKey(c, "stringId", nil)
		if key, err = datastore.Put(c, key, &profile); err != nil {
			return nil, fmt.Errorf("datastore.Put(with incomplete key %v) failed: %v", key, err)
		}
	} else if err != nil {
		// An unexpected error.
		return nil, err
	} else {
		// Update the Profile:
		updateProfileFromGoogleUserInfo(&profile, &userInfo)
		profile.GoogleAccessToken = *token

		if key, err = datastore.Put(c, key, &profile); err != nil {
			return nil, fmt.Errorf("datastore.Put(with key %v) failed: %v", key, err)
		}
	}

	return key, nil
}

func updateProfileFromGoogleUserInfo(profile *user.Profile, userInfo *GoogleUserInfo) {
	profile.GoogleId = userInfo.Sub
	profile.Name = userInfo.Name

	if userInfo.EmailVerified {
		profile.Email = userInfo.Email
	}
}

func getUserProfileById(c context.Context, userId *datastore.Key) (*user.Profile, error) {
	var profile user.Profile
	if err := datastore.Get(c, userId, &profile); err != nil {
		return nil, fmt.Errorf("datastore.Get() failed with key: %v: %v", userId, err)
	}

	return &profile, nil
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
	var loginFailedUrl = config.BaseUrl + "/login?failed=true"

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

	defaultSessionID     = "default"
	oauthTokenSessionKey = "oauth_token"
	userIdSessionKey     = "id" // A generic user ID, not a google user ID.
)

// We store the token in a session cookie.
var store *sessions.CookieStore