package main

import (
	"cloud.google.com/go/datastore"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/murraycu/go-bigoquiz-server/db"
	"github.com/murraycu/go-bigoquiz-server/user"
	"golang.org/x/oauth2"
	"net/http"
)

func restHandleUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	loginInfo, err := getLoginInfoFromSessionAndDb(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.Marshal(loginInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getProfileFromSession(r *http.Request) (*datastore.Key, *oauth2.Token, error) {
	session, err := store.Get(r, defaultSessionID)
	if err != nil {
		return nil, nil, fmt.Errorf("getLoginInfoFromSessionAndDb(): store.Get() failed: %v", err)
	}

	// Get the token from the cookie:
	tokenVal, ok := session.Values[oauthTokenSessionKey]
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
	userIdVal, ok := session.Values[userIdSessionKey]
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

func getProfileFromSessionAndDb(r *http.Request) (*user.Profile, *datastore.Key, *oauth2.Token, error) {
	userId, token, err := getProfileFromSession(r)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("getProfileFromSession() failed: %v", err)
	}

	if userId == nil {
		// Not an error.
		// It's just not in the session cookie.
		return nil, nil, nil, nil
	}

	c := r.Context()
	profile, err := db.GetUserProfileById(c, userId)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("GetUserProfileById() failed: %v", err)
	}

	return profile, userId, token, nil
}

func getLoginInfoFromSessionAndDb(r *http.Request) (*user.LoginInfo, error) {
	var loginInfo user.LoginInfo

	profile, userId, token, err := getProfileFromSessionAndDb(r)
	if err != nil {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = fmt.Sprintf("not logged in (%v)", err)
	}

	updateLoginInfoFromProfile(&loginInfo, profile, token, userId)

	return &loginInfo, err
}

func updateLoginInfoFromProfile(loginInfo *user.LoginInfo, profile *user.Profile, token *oauth2.Token, userId *datastore.Key) {
	if profile == nil {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (no profile found)"
	} else if !token.Valid() {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (invalid token)"
	} else {
		loginInfo.LoggedIn = true
		loginInfo.Nickname = profile.Name
		loginInfo.UserId = userId // Not for the JSON, but useful to callers.

		loginInfo.GoogleLinked = profile.GoogleId != ""
		loginInfo.GoogleProfileUrl = profile.GoogleProfileUrl
		loginInfo.GitHubLinked = profile.GitHubId != 0
		loginInfo.GitHubProfileUrl = profile.GitHubProfileUrl
		loginInfo.FacebookLinked = profile.FacebookId != ""
		loginInfo.FacebookProfileUrl = profile.FacebookProfileUrl
	}
}

/** Get the user ID.
 * Returns a nil Key, and a nil error, if the user is not logged in.
 */
func getUserIdFromSessionAndDb(r *http.Request, w http.ResponseWriter) (*datastore.Key, error) {
	loginInfo, err := getLoginInfoFromSessionAndDb(r)
	if err != nil {
		return nil, fmt.Errorf("getLoginInfoFromSessionAndDb() failed: %v", err)
	}

	return loginInfo.UserId, nil
}
