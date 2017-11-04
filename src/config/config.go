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

	// This file contains other secrets, such as the keys for the encrypted cookie store.
	// The file format is like so:
	//
	// {
	//   "cookie-store-key": "something-secret"
	// }
	configFilename = "config.json"

	// This file must be downloaded
	// (via the "DOWNLOAD JSON" link at https://console.developers.google.com/apis/credentials/oauthclient )
	// and added with this exact filename, next to the bigoquiz.go source file.
	googleConfigCredentialsFilename = "config_google_oauth2_credentials_secret.json"

	// See https://developers.google.com/+/web/api/rest/oauth#profile
	googleCredentialsScopeProfile = "profile"

	// See https://developers.google.com/identity/protocols/googlescopes
	googleCredentialsScopeEmail = "https://www.googleapis.com/auth/userinfo.email"

	// This has the same format, and location, as googleConfigCredentialsFilename,
	// but is maintained manually instead of being downloaded.
	githubConfigCredentialsFilename = "config_github_oauth2_credentials_secret.json"

	// See https://developer.github.com/apps/building-integrations/setting-up-and-registering-oauth-apps/about-scopes-for-oauth-apps/
	githubCredentialsScopeUser  = "read:user"
	githubCredentialsScopeEmail = "user:email"

	// This has the same format, and location, as googleConfigCredentialsFilename,
	// but is maintained manually instead of being downloaded.
	facebookConfigCredentialsFilename = "config_facebook_oauth2_credentials_secret.json"

	// See https://developers.facebook.com/docs/facebook-login/permissions
	facebookCredentialsScopePublicProfile = "public_profile"
	facebookCredentialsScopeEmail         = "email"
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

/** Get an oauth2 Config object based on the secret .json file,
 * in the Google format expected by google.ConfigFromJSON(),
 * though this is also for non-Google credentials.
 */
func generateOAuthConfig(r *http.Request, credentialsFilename string, scope ...string) *oauth2.Config {
	c := appengine.NewContext(r)

	b, err := ioutil.ReadFile(credentialsFilename)
	if err != nil {
		log.Errorf(c, "Unable to read client secret file (%s): %v", credentialsFilename, err)
		return nil
	}

	config, err := google.ConfigFromJSON(b, scope...)
	if err != nil {
		log.Errorf(c, "Unable to parse client secret file (%) to config: %v", credentialsFilename, err)
		return nil
	}

	return config
}

/** Get an oauth2 Config object based on the secret .json file.
 * See googleConfigCredentialsFilename.
 */
func GenerateGoogleOAuthConfig(r *http.Request) *oauth2.Config {
	return generateOAuthConfig(r, googleConfigCredentialsFilename, googleCredentialsScopeProfile, googleCredentialsScopeEmail)
}

/** Get an oauth2 Config object based on the secret .json file.
 * See githubConfigCredentialsFilename.
 */
func GenerateGitHubOAuthConfig(r *http.Request) *oauth2.Config {
	return generateOAuthConfig(r, githubConfigCredentialsFilename, githubCredentialsScopeUser, githubCredentialsScopeEmail)
}

/** Get an oauth2 Config object based on the secret .json file.
 * See githubConfigCredentialsFilename.
 */
func GenerateFacebookOAuthConfig(r *http.Request) *oauth2.Config {
	return generateOAuthConfig(r, facebookConfigCredentialsFilename, facebookCredentialsScopePublicProfile, facebookCredentialsScopeEmail)
}
