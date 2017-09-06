package bigoquiz

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"net/http"
	"user"
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

	token, ok := session.Values[oauthTokenSessionKey].(*oauth2.Token)
	if !ok {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (no oauthTokenSessionKey)"
		return &loginInfo, nil
	}

	// TODO: Get the name from the database, not from the cookie.
	nameVal, ok := session.Values[nameSessionKey]
	if !ok {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (no name as value)."
		return &loginInfo, nil
	}

	var name string
	name, ok = nameVal.(string)
	if !ok {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (no name as string). nameVal is not a string."
		return &loginInfo, nil
	}

	if token.Valid() {
		loginInfo.LoggedIn = true
		loginInfo.Nickname = name
	} else {
		loginInfo.LoggedIn = false
		loginInfo.ErrorMessage = "not logged in user (invalid token)"
	}

	return &loginInfo, err
}
