package bigoquiz

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"io"
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

	const baseUrl = "http://beta.bigoquiz.com/"

	token, err := config.Exchange(c, code)
	if err != nil {
		log.Errorf(c, "config.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, baseUrl, http.StatusTemporaryRedirect)
		return
	}

	client := config.Client(c, token)
	infoResponse, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Errorf(c, "client.Get() failed with '%s'\n", err)
		http.Redirect(w, r, baseUrl, http.StatusTemporaryRedirect)
		return
	}

	// Just to show that it worked:
	defer infoResponse.Body.Close()
	body, err := ioutil.ReadAll(infoResponse.Body)
	if err != nil {
		log.Errorf(c, "ReadAll(body) failed with '%s'\n", err)
		http.Redirect(w, r, baseUrl, http.StatusTemporaryRedirect)
		return
	}

	var userinfo GoogleUserInfo
	json.Unmarshal(body, &userinfo)

	io.WriteString(w, userinfo.Name)
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
