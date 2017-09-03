package bigoquiz

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"user"
)

func restHandleUserHistoryAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: Use actual authentication.
	var info user.UserHistoryOverall
	info.LoginInfo.Nickname = "example@example.com"

	jsonStr, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonStr)
}

func restHandleUserHistoryByQuizId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: Use actual authentication.
	var info user.UserHistorySections
	info.LoginInfo.Nickname = "example@example.com"

	jsonStr, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonStr)
}
