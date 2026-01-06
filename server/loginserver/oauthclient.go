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

	"github.com/murraycu/go-bigoquiz-server/config"
	"github.com/murraycu/go-bigoquiz-server/repositories/db"
	"github.com/murraycu/go-bigoquiz-server/server/usersessionstore"
	"golang.org/x/oauth2"
)

type OAuthClient struct {
	oAuthStateClient *db.OAuthStateDataRepository

	// Session cookie store.
	userSessionStore usersessionstore.UserSessionStore

	confOAuthGoogle   *oauth2.Config
	confOAuthGitHub   *oauth2.Config
	confOAuthFacebook *oauth2.Config

	config *config.Config

	userDataClient db.UserDataRepository
}

func NewOAuthClient(userSessionStore usersessionstore.UserSessionStore, userDataClient db.UserDataRepository, conf *config.Config) (*OAuthClient, error) {
	result := &OAuthClient{}
	result.config = conf

	var err error
	result.oAuthStateClient, err = db.NewOAuthStateDataRepository()
	if err != nil {
		return nil, fmt.Errorf("NewOAuthStateDataRepository() failed: %v", err)
	}

	result.userSessionStore = userSessionStore
	result.userDataClient = userDataClient

	result.confOAuthGoogle, err = config.GenerateGoogleOAuthConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("unable to generate Google OAuth config: %v", err)
	}

	result.confOAuthGitHub, err = config.GenerateGitHubOAuthConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("unable to generate GitHub OAuth config: %v", err)
	}

	result.confOAuthFacebook, err = config.GenerateFacebookOAuthConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("unable to generate Facebook OAuth config: %v", err)
	}

	return result, nil
}

/** Get an oauth2 URL based on the oauth config.
 */
func (o *OAuthClient) generateOAuthUrl(r *http.Request, oauthConfig *oauth2.Config) (string, error) {
	ctx := r.Context()

	state, err := o.generateOAuthState(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to generate state: %v", err)
	}

	// Use oauth2.AccessTypeOffline ("Offline Access"), instead of oauth2.AccessTyoeOnline, so we also receive an OAuth
	// refresh token (longer lived), not just an OAuth access token (short-lived - approximately 30 minutes). We will
	// use the refresh token to retrieve a new access token.
	// This is appropriate because this is a "Web Server Application" (also known as a "Server-side Web App"):
	// https://developers.google.com/identity/protocols/oauth2/web-server
	// so it's considered safe for us to store the (longer lived) refresh token on the server.
	//
	// In contrast, a "Client-side Web Applications" would not be a safe place to expose the refresh token. So it would
	// instead redirect the web page to the Google/GitHub/Facebook login page again, which would, most of the time,
	// just quickly redirect back again without any user interaction.
	//
	// Use oauth2.ApprovalForce to specify the "prompt" option.
	// (This seems to be necessary to get the RefreshToken too, though maybe only after a previous consent was already
	// granted without the "Offline Access".
	return oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

func (o *OAuthClient) generateOAuthState(ctx context.Context) (string, error) {
	state := rand.Int63()
	err := o.oAuthStateClient.StoreOAuthState(ctx, state)
	if err != nil {
		return "", fmt.Errorf("StoreOAuthState() failed: %v", err)
	}

	return strconv.FormatInt(state, 10), nil
}

func (o *OAuthClient) checkOAuthResponseState(ctx context.Context, state string) error {
	stateNum, err := strconv.ParseInt(state, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt() failed: %v", err)
	}

	err = o.oAuthStateClient.CheckOAuthState(ctx, stateNum)
	if err != nil {
		return fmt.Errorf("db.CheckOAuthState() failed: %v", err)
	}

	return nil
}

func (o *OAuthClient) removeOAuthState(ctx context.Context, state string) error {
	stateNum, err := strconv.ParseInt(state, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt() failed: %v", err)
	}

	return o.oAuthStateClient.RemoveOAuthState(ctx, stateNum)
}

func (o *OAuthClient) checkOAuthResponseStateAndGetCode(ctx context.Context, r *http.Request) (string, error) {
	state := r.FormValue("state")
	err := o.checkOAuthResponseState(ctx, state)
	if err != nil {
		return "", fmt.Errorf("invalid oauth state ('%s): %v", state, err)
	}

	// The state will not be used again,
	// so remove it from the datastore.
	err = o.removeOAuthState(ctx, state)
	if err != nil {
		return "", fmt.Errorf("removeOAuthState() failed: %v", err)
	}

	return r.FormValue("code"), nil
}

func (o *OAuthClient) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	handleOAuthCallback(o, w, r, "https://www.googleapis.com/oauth2/v3/userinfo", o.confOAuthGoogle, o.userDataClient.StoreGoogleLoginInUserProfile)
}

func (o *OAuthClient) HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	handleOAuthCallback(o, w, r, "https://api.github.com/user", o.confOAuthGitHub, o.userDataClient.StoreGitHubLoginInUserProfile)
}

func (o *OAuthClient) HandleFacebookCallback(w http.ResponseWriter, r *http.Request) {
	handleOAuthCallback(o, w, r, "https://graph.facebook.com/me?fields=link,name,email", o.confOAuthFacebook, o.userDataClient.StoreFacebookLoginInUserProfile)
}

// handleOAuthCallback handles the OAuth callback response, interpreting it as the specific OAuthUserInfo struct type.
func handleOAuthCallback[OAuthUserInfo any](o *OAuthClient, w http.ResponseWriter, r *http.Request, userInfoUrl string, conf *oauth2.Config, storeLogin func(context.Context, OAuthUserInfo, string, *oauth2.Token) (string, error)) {
	ctx := r.Context()

	checkStateResult, err := o.checkOAuthResponseStateAndGetBody(w, r, conf, userInfoUrl, ctx)
	if err != nil {
		o.loginFailed("checkOAuthResponseStateAndGetBody() failed", err, w, r)
		return
	}

	var userinfo OAuthUserInfo
	err = json.Unmarshal(checkStateResult.body, &userinfo)
	if err != nil {
		o.loginFailed("Unmarshalling of JSON from oauth2 callback failed", err, w, r)
		return
	}

	// Get the existing logged-in user's userId, if any, from the cookie, if any:
	userId, _, err := o.userSessionStore.GetUserIdAndOAuthTokenFromSession(r)
	if err != nil {
		o.loginFailed("getProfileFromSession() failed", err, w, r)
		return
	}

	userId, err = storeLogin(ctx, userinfo, userId, checkStateResult.token)
	if err != nil {
		o.loginFailed("storeLogin() failed", err, w, r)
		return
	}

	o.storeCookieAndRedirect(r, w, ctx, userId, checkStateResult.token)
}

type CheckStateResult struct {
	token *oauth2.Token
	body  []byte

	// For instance, this will be true if (but not only if) the token has expired.
	// (We should expect the OAuth token to expire quickly - for instance, within 30 minutes.)
	invalidToken bool
}

func (o *OAuthClient) checkOAuthResponseStateAndGetBody(w http.ResponseWriter, r *http.Request, conf *oauth2.Config, url string, ctx context.Context) (*CheckStateResult, error) {
	code, err := o.checkOAuthResponseStateAndGetCode(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("checkOAuthResponseStateAndGetCode() failed: %v", err)
	}

	return o.exchangeAndGetUserBody(w, r, conf, code, url, ctx)
}

func (o *OAuthClient) exchangeAndGetUserBody(w http.ResponseWriter, r *http.Request, conf *oauth2.Config, code string, url string, ctx context.Context) (*CheckStateResult, error) {
	// Extract the token, which will have the
	// - OAuth access token
	// - OAuth refresh code, because we specified oauth2.AccessTypeOffline to oauth2.Config.AuthCodeURL().
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("config.Exchange() failed: %v", err)
	}

	if !token.Valid() {
		return &CheckStateResult{
			token:        nil,
			body:         nil,
			invalidToken: true,
		}, fmt.Errorf("loginFailedUrl.Exchange() returned an invalid token")
	}

	client := conf.Client(ctx, token)
	infoResponse, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("client.Get() failed: %v", err)
	}

	defer func() {
		err := infoResponse.Body.Close()
		if err != nil {
			o.loginFailed("Body.Close() failed", err, w, r)
		}
	}()

	body, err := ioutil.ReadAll(infoResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll(body) failed: %v", err)
	}

	return &CheckStateResult{
		token: token,
		body:  body,
	}, nil
}

// Called after user info has been successfully stored in the database.
func (o *OAuthClient) storeCookieAndRedirect(r *http.Request, w http.ResponseWriter, ctx context.Context, strUserId string, token *oauth2.Token) {
	// Store the token in the cookie
	// so we can retrieve it from subsequent requests from the browser.
	session, err := o.userSessionStore.GetSession(r)
	if err != nil {
		o.loginFailed("Could not create new session", err, w, r)
		return
	}

	session.Values[usersessionstore.OAuthTokenSessionKey] = token
	session.Values[usersessionstore.UserIdSessionKey] = strUserId

	if err := session.Save(r, w); err != nil {
		o.loginFailed("Could not save session", err, w, r)
		return
	}

	// Redirect the user back to a page to show they are logged in:
	var userProfileUrl = o.config.BaseUrl + "/user"
	http.Redirect(w, r, userProfileUrl, http.StatusFound)
}

func (o *OAuthClient) loginFailed(message string, err error, w http.ResponseWriter, r *http.Request) {
	var loginFailedUrl = o.config.BaseUrl + "/login?failed=true"

	log.Printf(message+":'%v'\n", err)
	http.Redirect(w, r, loginFailedUrl, http.StatusTemporaryRedirect)
}

func (o *OAuthClient) RedirectToGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect the user to the Google login page:
	url, err := o.generateOAuthUrl(r, o.confOAuthGoogle)
	if err != nil {
		o.loginFailed("generateOAuthUrl() failed", err, w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (o *OAuthClient) RedirectToGitHubLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect the user to the GitHub login page:
	url, err := o.generateOAuthUrl(r, o.confOAuthGitHub)
	if err != nil {
		o.loginFailed("generateOAuthUrl() failed", err, w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (o *OAuthClient) RedirectToFacebookLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect the user to the Facebook login page:
	url, err := o.generateOAuthUrl(r, o.confOAuthFacebook)
	if err != nil {
		o.loginFailed("generateOAuthUrl() failed", err, w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
