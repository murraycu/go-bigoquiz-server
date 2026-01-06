package loginserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/murraycu/go-bigoquiz-server/config"
	"github.com/murraycu/go-bigoquiz-server/repositories/db"
	"github.com/murraycu/go-bigoquiz-server/server/usersessionstore"
)

type LoginServer struct {
	userDataClient db.UserDataRepository

	// Session cookie store.
	userSessionStore usersessionstore.UserSessionStore

	oauthClient *OAuthClient
}

func NewLoginServer(userSessionStore usersessionstore.UserSessionStore, conf *config.Config) (*LoginServer, error) {
	result := &LoginServer{}

	result.userSessionStore = userSessionStore

	userDataClient, err := db.NewUserDataRepository()
	if err != nil {
		return nil, fmt.Errorf("NewUserDataRepository() failed: %v", err)
	}

	result.oauthClient, err = NewOAuthClient(userSessionStore, userDataClient, conf)
	if err != nil {
		return nil, fmt.Errorf("NewOAuthClient() failed: %v", err)
	}

	return result, nil
}

func (s *LoginServer) HandleGoogleLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.oauthClient.RedirectToGoogleLogin(w, r)
}

func (s *LoginServer) HandleGoogleCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.oauthClient.HandleGoogleCallback(w, r)
}

func (s *LoginServer) HandleGitHubCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.oauthClient.HandleGitHubCallback(w, r)
}

func (s *LoginServer) HandleFacebookCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.oauthClient.HandleFacebookCallback(w, r)
}

func (s *LoginServer) HandleLogout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Wipe the cookie:
	session, err := s.userSessionStore.GetSession(r)
	if err != nil {
		logoutError("could not get default session", err, w)
		return
	}

	session.Options.MaxAge = -1 // Clear session.

	if err := session.Save(r, w); err != nil {
		logoutError("Could not save session", err, w)
		return
	}

	redirectURL := r.FormValue("redirect")
	if redirectURL == "" {
		redirectURL = "/"
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (s *LoginServer) HandleGitHubLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.oauthClient.RedirectToGitHubLogin(w, r)
}

func (s *LoginServer) HandleFacebookLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.oauthClient.RedirectToFacebookLogin(w, r)
}

func logoutError(message string, err error, w http.ResponseWriter) {
	handleErrorAsHttpError(w, http.StatusInternalServerError, "message: %v", err)
}

func handleErrorAsHttpError(w http.ResponseWriter, code int, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	log.Print(msg)

	http.Error(w, msg, code)
}
