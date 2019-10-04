package usersessionstore

import (
	"cloud.google.com/go/datastore"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/murraycu/go-bigoquiz-server/quiz"
	"github.com/murraycu/go-bigoquiz-server/repositories/db"
	"golang.org/x/oauth2"
	"net/http"
)

const OAuthTokenSessionKey = "oauth_token"
const DefaultSessionID = "default"
const UserIdSessionKey = "id" // A generic user ID, not a google user ID.

type UserSessionStore struct {
	Quizzes           map[string]*quiz.Quiz
	QuizzesListSimple []*quiz.Quiz
	QuizzesListFull   []*quiz.Quiz

	UserDataClient *db.UserDataRepository

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

func (s *UserSessionStore) GetProfileFromSession(r *http.Request) (*datastore.Key, *oauth2.Token, error) {
	session, err := s.Store.Get(r, DefaultSessionID)
	if err != nil {
		return nil, nil, fmt.Errorf("getLoginInfoFromSessionAndDb(): store.Get() failed: %v", err)
	}

	// Get the token from the cookie:
	tokenVal, ok := session.Values[OAuthTokenSessionKey]
	if !ok {
		// Not an error.
		// It's just not in the cookie.
		return nil, nil, nil
	}

	// Try casting it to the expected type:
	var token *oauth2.Token
	token, ok = tokenVal.(*oauth2.Token)
	if !ok {
		return nil, nil, fmt.Errorf("oauthTokenSessionKey is not a *Token")
	}

	// Get the name from the database, via the userID from the cookie:
	userIdVal, ok := session.Values[UserIdSessionKey]
	if !ok {
		return nil, nil, fmt.Errorf("no name as value")
	}

	// Try casting it to the expected type:
	var userId *datastore.Key
	userId, ok = userIdVal.(*datastore.Key)
	if !ok {
		return nil, nil, fmt.Errorf("no name as *Key. userIdVal is not a *Key")
	}

	if userId == nil {
		return nil, nil, fmt.Errorf("userId is null")
	}

	return userId, token, nil
}
