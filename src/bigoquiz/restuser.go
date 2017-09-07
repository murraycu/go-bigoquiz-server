package bigoquiz

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"net/http"
	"user"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine"
	"fmt"
)

func restHandleUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	loginInfo, err := loginInfoFromSession(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.Marshal(loginInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonStr)
}

func loginInfoFromSession(r *http.Request, w http.ResponseWriter) (*user.LoginInfo, error) {
	session, err := store.Get(r, defaultSessionID)
	if err != nil {
		return nil, err
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

	c := appengine.NewContext(r)
	profile, err := getUserProfileById(c, userId)
	if err != nil {
		return nil, fmt.Errorf("getUserProfileById() failed: %v", err)
	}

	if token.Valid() {
		loginInfo.LoggedIn = true
		loginInfo.Nickname = profile.Name
	} else {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (invalid token)"
	}

	return &loginInfo, err
}
