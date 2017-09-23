package bigoquiz

import (
	"db"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
	"user"
)

func restHandleUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	loginInfo, err := getLoginInfoFromSessionAndDb(r, w)
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

func getLoginInfoFromSessionAndDb(r *http.Request, w http.ResponseWriter) (*user.LoginInfo, error) {
	session, err := store.Get(r, defaultSessionID)
	if err != nil {
		return nil, fmt.Errorf("getLoginInfoFromSessionAndDb(): store.Get() failed: %v", err)
	}

	var loginInfo user.LoginInfo

	// Get the token from the cookie:
	tokenVal, ok := session.Values[oauthTokenSessionKey]
	if !ok {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (no oauthTokenSessionKey)"
		return &loginInfo, nil
	}

	// Try casting it to the expected type:
	var token *oauth2.Token
	token, ok = tokenVal.(*oauth2.Token)
	if !ok {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (oauthTokenSessionKey is not a *Token)"
		return &loginInfo, nil
	}

	// Get the name from the database, via the userID from the cookie:
	userIdVal, ok := session.Values[userIdSessionKey]
	if !ok {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (no name as value)."
		return &loginInfo, nil
	}

	// Try casting it to the expected type:
	var userId *datastore.Key
	userId, ok = userIdVal.(*datastore.Key)
	if !ok {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (no name as *Key). userIdVal is not a *Key."
		return &loginInfo, nil
	}

	if userId == nil {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (userId is null)."
		return &loginInfo, nil
	}

	c := appengine.NewContext(r)
	profile, err := db.GetUserProfileById(c, userId)
	if err != nil {
		return nil, fmt.Errorf("getUserProfileById() failed: %v", err)
	}

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
	}

	return &loginInfo, err
}

/** Get the user ID.
 * Returns a nil Key, and a nil error, if the user is not logged in.
 */
func getUserIdFromSessionAndDb(r *http.Request, w http.ResponseWriter) (*datastore.Key, error) {
	loginInfo, err := getLoginInfoFromSessionAndDb(r, w)
	if err != nil {
		return nil, fmt.Errorf("getLoginInfoFromSessionAndDb() failed: %v", err)
	}

	return loginInfo.UserId, nil
}
