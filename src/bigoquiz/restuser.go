package bigoquiz

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"user"
)

func restHandleUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: Use actual authentication.
	var loginInfo user.LoginInfo
	loginInfo.Nickname = "example@example.com"

	jsonStr, err := json.Marshal(loginInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonStr)
}
