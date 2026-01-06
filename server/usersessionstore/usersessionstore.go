package usersessionstore

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

const OAuthTokenSessionKey = "oauth_token"
const DefaultSessionID = "default"
const UserIdSessionKey = "id" // A generic user ID, not a google user ID.

type UserSessionStore interface {
	GetSession(r *http.Request) (*sessions.Session, error)
	GetUserIdAndOAuthTokenFromSession(r *http.Request) (string, *oauth2.Token, error)
}

type UserSessionStoreImpl struct {
	// Session cookie store.
	store *sessions.CookieStore
}

func NewUserSessionStore(cookieKey string) (UserSessionStore, error) {
	result := &UserSessionStoreImpl{}

	// Create the session cookie store,
	// using the secret key from the configuration file.
	result.store = sessions.NewCookieStore([]byte(cookieKey))
	result.store.Options.HttpOnly = true
	result.store.Options.Secure = true // Only send via HTTPS connections, not HTTP.

	return result, nil
}

func (s *UserSessionStoreImpl) GetSession(r *http.Request) (*sessions.Session, error) {
	result, err := s.store.Get(r, DefaultSessionID)
	if err != nil {
		return nil, fmt.Errorf("store.Get() failed: %v", err)
	}

	return result, nil
}

func (s *UserSessionStoreImpl) GetUserIdAndOAuthTokenFromSession(r *http.Request) (string, *oauth2.Token, error) {
	session, err := s.GetSession(r)
	if err != nil {
		return "", nil, fmt.Errorf("GetUserIdAndOAuthTokenFromSession(): store.Get() failed: %v", err)
	}

	// Get the oauth2 token from the cookie:
	// (If the cookie has no token then the user is not logged in.)
	tokenVal, ok := session.Values[OAuthTokenSessionKey]
	if !ok {
		// Not an error.
		// It's just not in the cookie.
		return "", nil, nil
	}

	// Try casting it to the expected type:
	var token *oauth2.Token
	token, ok = tokenVal.(*oauth2.Token)
	if !ok {
		return "", nil, fmt.Errorf("oauthTokenSessionKey is not a *Token")
	}

	// Get the userID from the cookie:
	userIdVal, ok := session.Values[UserIdSessionKey]
	if !ok {
		// Not an error.
		// It's just not in the cookie.
		// (the user is not logged in.)
		return "", nil, nil
	}

	// Try casting it to the expected type:
	var strUserId string
	strUserId, ok = userIdVal.(string)
	if !ok {
		// Not an error.
		// It's just not (correctly) in the cookie.
		// (We changed its format in 2019/10.)
		return "", nil, nil
	}

	userId, err := datastore.DecodeKey(strUserId)
	if err != nil {
		return "", nil, fmt.Errorf("datastore.DecodeKey() failed: %v", err)
	}

	if userId == nil {
		return "", nil, fmt.Errorf("userId is null")
	}

	return userId.Encode(), token, nil
}
