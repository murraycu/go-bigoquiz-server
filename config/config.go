package config

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"path/filepath"
)

const (
	BaseUrl = "https://bigoquiz.com"
	// When running angular-bigoquiz-client with ng serve: BaseUrl = "http://localhost:4200"

	BaseApiUrl = "https://api.bigoquiz.com"
	// When running locally: BaseApuUrl = "http://localhost:8080

	PART_URL_LOGIN_CALLBACK_GOOGLE   = "callback-google"
	PART_URL_LOGIN_CALLBACK_GITHUB   = "callback-github"
	PART_URL_LOGIN_CALLBACK_FACEBOOK = "callback-facebook"

	// This file contains other secrets, such as the keys for the encrypted cookie store.
	// The file format is like so:
	//
	// {
	//   "cookie-store-key": "something-secret"
	// }
	configFilename = "config.json"

	// See https://developers.google.com/+/web/api/rest/oauth#profile
	googleCredentialsScopeProfile = "profile"

	// See https://developers.google.com/identity/protocols/googlescopes
	googleCredentialsScopeEmail = "https://www.googleapis.com/auth/userinfo.email"

	// This has the same format, and location, as googleConfigCredentialsFilename,
	// but is maintained manually instead of being downloaded.

	// for clues

	// See https://developer.github.com/apps/building-integrations/setting-up-and-registering-oauth-apps/about-scopes-for-oauth-apps/
	githubCredentialsScopeUser  = "read:user"
	githubCredentialsScopeEmail = "user:email"

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
		return nil, fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	return &result, nil
}

/** Get an oauth2 Config object based on the secret .json file,
 * These files contains the client_id and client_secret for the OAuth2 authentication.
 * See github_credentials_secret.json.example, for instance.
 *
 * Google:
 * (See the "DOWNLOAD JSON" link at https://console.developers.google.com/apis/credentials/oauthclient )
 * and added with this exact filename, next to the main.go source file.
 *
 * Github:
 * See https://github.com/settings/applications/SOME_ID_OF_YOUR_OWN_APP
 * and https://github.com/golang/oauth2
 *
 * Facebook:
 * See https://developers.facebook.com/apps/YOUR_OWN_APPS_CLIENT_ID/settings/basic/
 * and https://github.com/golang/oauth2/blob/master/facebook/facebook.go
 * for clues.
 */
func addSecretsToOAuthConfig(credentialsFilenamePrefix string, config oauth2.Config) (*oauth2.Config, error) {
	credentialsFilename := fmt.Sprintf("%s_credentials_secret.json", credentialsFilenamePrefix)
	path := filepath.Join("config_oauth2", credentialsFilename)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file (%s): %v", credentialsFilename, err)
	}

	// Load the secrets:
	type cred struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}

	var secretConfig cred
	if err := json.Unmarshal(b, &secretConfig); err != nil {
		return nil, fmt.Errorf("unable to parse client secret file (%s) to config (json.Unmarshal() failed): %v", credentialsFilename, err)
	}

	config.ClientID = secretConfig.ClientID
	config.ClientSecret = secretConfig.ClientSecret

	return &config, nil
}

/** Get an oauth2 Config object based on the secret .json file.
 * See googleConfigCredentialsFilename.
 */
func GenerateGoogleOAuthConfig() (*oauth2.Config, error) {
	config := oauth2.Config{
		ClientID:     "", // Filled in from secrets
		ClientSecret: "", // Filled in from secrets
		Endpoint:     google.Endpoint,
		RedirectURL:  callbackUrl(PART_URL_LOGIN_CALLBACK_GOOGLE),
		Scopes:       []string{googleCredentialsScopeProfile, googleCredentialsScopeEmail},
	}

	result, err := addSecretsToOAuthConfig("google", config)
	if err != nil {
		return nil, fmt.Errorf("addSecretsToOAuthConfig() failed: %v", err)
	}

	return result, nil
}

func callbackUrl(suffix string) string {
	return BaseApiUrl + "/login/" + suffix
}

/** Get an oauth2 Config object based on the secret .json file.
 * See githubConfigCredentialsFilename.
 */
func GenerateGitHubOAuthConfig() (*oauth2.Config, error) {
	config := oauth2.Config{
		ClientID:     "", // Filled in from secrets
		ClientSecret: "", // Filled in from secrets
		Endpoint:     github.Endpoint,
		RedirectURL:  callbackUrl(PART_URL_LOGIN_CALLBACK_GITHUB),
		Scopes:       []string{githubCredentialsScopeUser, githubCredentialsScopeEmail},
	}

	return addSecretsToOAuthConfig("github", config)
}

/** Get an oauth2 Config object based on the secret .json file.
 * See githubConfigCredentialsFilename.
 */
func GenerateFacebookOAuthConfig() (*oauth2.Config, error) {
	config := oauth2.Config{
		ClientID:     "", // Filled in from secrets
		ClientSecret: "", // Filled in from secrets
		Endpoint:     facebook.Endpoint,
		RedirectURL:  callbackUrl(PART_URL_LOGIN_CALLBACK_FACEBOOK),
		Scopes:       []string{facebookCredentialsScopePublicProfile, facebookCredentialsScopeEmail},
	}

	return addSecretsToOAuthConfig("facebook", config)
}
