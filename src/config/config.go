package config

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"io/ioutil"
	"net/http"
)

const (
	BaseUrl = "https://beta.bigoquiz.com"
	// When running angular-bigoquiz-client with ng serve: BaseUrl = "http://localhost:4200"

	// This file must be downloaded
	// (via the "DOWNLOAD JSON" link at https://console.developers.google.com/apis/credentials/oauthclient )
	// and added with this exact filename, next to the bigoquiz.go source file.
	configCredentialsFilename = "config_google_oauth2_credentials_secret.json"

	// This file contains other secrets, such as the keys for the encrypted cookie store.
	// The file format is like so:
	//
	// {
	//   "cookie-store-key": "something-secret"
	// }
	configFilename = "config.json"

	// See https://developers.google.com/+/web/api/rest/oauth#profile
	credentialsScopeProfile = "profile"

	// See https://developers.google.com/identity/protocols/googlescopes
	credentialsScopeEmail = "https://www.googleapis.com/auth/userinfo.email"
)

/** Get general configuration.
 * See configFilename.
 */
type Config struct {
	CookieKey string
}

func GenerateConfig() (*Config, error) {
	b, err := ioutil.ReadFile(configFilename)
	if err != nil {
		// log.Errorf("Unable to read config file (%s): %v", configFilename, err)
		return nil, err
	}

	var result Config
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal failed", err)
	}

	return &result, nil
}

/** Get an oauth2 Config object based on the secret .json file.
 * See configCredentialsFilename.
 */
func GenerateGoogleOAuthConfig(r *http.Request) *oauth2.Config {
	c := appengine.NewContext(r)

	b, err := ioutil.ReadFile(configCredentialsFilename)
	if err != nil {
		log.Errorf(c, "Unable to read client secret file (%s): %v", configCredentialsFilename, err)
		return nil
	}

	config, err := google.ConfigFromJSON(b, credentialsScopeProfile, credentialsScopeEmail)
	if err != nil {
		log.Errorf(c, "Unable to parse client secret file (%) to config: %v", configCredentialsFilename, err)
		return nil
	}

	return config
}
