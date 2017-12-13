package db

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

/**
 */
type OAuthState struct {
	// The datastore ID (int64) is the oauth2 state.
	timestamp time.Time
}

func StoreOAuthState(c context.Context, state int64) error {
	key := datastore.NewKey(c, DB_KIND_OAUTH_STATE, "", state, nil)

	var stateObj OAuthState;

	// Store a timestamp so a cron job can periodically remove old states.
	stateObj.timestamp = time.Now().UTC()

	_, err := datastore.Put(c, key, &stateObj)

	/*
	if err != nil {
		log.Errorf(c, "datastore.Put() failed: %v", err)
	}
	*/

	return err
}


func CheckOAuthState(c context.Context, state int64) bool {
	key := datastore.NewKey(c, DB_KIND_OAUTH_STATE, "", state, nil)

	var stateObj OAuthState
	err := datastore.Get(c, key, &stateObj)

	/*
	if err != nil {
		log.Errorf(c, "datastore.Get() failed: %v", err)
	}
	*/

	return err == nil
}