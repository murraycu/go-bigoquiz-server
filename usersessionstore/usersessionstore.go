package usersessionstore

import (
	"cloud.google.com/go/datastore"
	"fmt"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"net/http"
)

const OAuthTokenSessionKey = "oauth_token"
const DefaultSessionID = "default"
const UserIdSessionKey = "id" // A generic user ID, not a google user ID.

type UserSessionStore struct {
	// Session cookie store.
	Store *sessions.CookieStore
}

func NewUserSessionStore(cookieKey string) (*UserSessionStore, error) {
	result := &UserSessionStore{}

	// Create the session cookie store,
	// using the secret key from the configuration file.
	result.Store = sessions.NewCookieStore([]byte(cookieKey))
	result.Store.Options.HttpOnly = true
	result.Store.Options.Secure = true // Only send via HTTPS connections, not HTTP.

	return result, nil
}

func (s *UserSessionStore) GetProfileFromSession(r *http.Request) (string, *oauth2.Token, error) {
	session, err := s.Store.Get(r, DefaultSessionID)
	if err != nil {
		return "", nil, fmt.Errorf("GetProfileFromSession(): store.Get() failed: %v", err)
	}

	// Get the token from the cookie:
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

	// Get the name from the database, via the userID from the cookie:
	userIdVal, ok := session.Values[UserIdSessionKey]
	if !ok {
		// Not an error.
		// It's just not in the cookie.
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
