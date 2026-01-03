package loginserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/murraycu/go-bigoquiz-server/config"
	"github.com/murraycu/go-bigoquiz-server/server"
	"github.com/murraycu/go-bigoquiz-server/server/usersessionstore"
)

type LoginServer struct {
	// Session cookie store.
	userSessionStore usersessionstore.UserSessionStore

	oauthLogins *server.OAuthLogins
}

func NewLoginServer(userSessionStore usersessionstore.UserSessionStore, conf *config.Config) (*LoginServer, error) {
	result := &LoginServer{}

	result.userSessionStore = userSessionStore

	var err error
	result.oauthLogins, err = server.NewOAuthLogins(conf, userSessionStore)
	if err != nil {
		return nil, fmt.Errorf("NewOAuthLogins() failed: %v", err)
	}

	return result, nil
}

func (s *LoginServer) HandleGoogleLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.oauthLogins.RedirectToGoogleLogin(w, r)
}

func (s *LoginServer) HandleGitHubLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.oauthLogins.RedirectToGitHubLogin(w, r)
}

func (s *LoginServer) HandleFacebookLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.oauthLogins.RedirectToFacebookLogin(w, r)
}

func (s *LoginServer) HandleGoogleCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.oauthLogins.HandleGoogleCallback(w, r)
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

/** Get an oauth2 URL based on the secret .json file.
 * See githubConfigCredentialsFilename.
 */
func (s *LoginServer) HandleGitHubCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.oauthLogins.HandleGitHubCallback(w, r)
}

func (s *LoginServer) HandleFacebookCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.oauthLogins.HandleFacebookCallback(w, r)
}

func logoutError(message string, err error, w http.ResponseWriter) {
	handleErrorAsHttpError(w, http.StatusInternalServerError, "messsage: %v", err)
}

func handleErrorAsHttpError(w http.ResponseWriter, code int, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	log.Print(msg)

	http.Error(w, msg, code)
}
