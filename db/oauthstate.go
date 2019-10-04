package db

import (
	"cloud.google.com/go/datastore"
	"fmt"
	"golang.org/x/net/context"
	"time"
)

/**
 */
type OAuthState struct {
	// The datastore ID (int64) is the oauth2 state.
	timestamp time.Time
}

func stateKey(c context.Context, state int64) *datastore.Key {
	return datastore.IDKey(DB_KIND_OAUTH_STATE, state, nil)
}

func StoreOAuthState(c context.Context, state int64) error {
	key := stateKey(c, state)

	var stateObj OAuthState

	// Store a timestamp so a cron job can periodically remove old states.
	stateObj.timestamp = time.Now().UTC()

	client, err := dataStoreClient(c)
	if err != nil {
		return fmt.Errorf("datastoreClient() failed: %v", err)
	}
	defer client.Close()

	_, err = client.Put(c, key, &stateObj)
	if err != nil {
		return fmt.Errorf("datastore.Put() failed: %v", err)
	}

	return err
}

func CheckOAuthState(c context.Context, state int64) bool {
	key := stateKey(c, state)

	client, err := dataStoreClient(c)
	if err != nil {
		return false
		// TODO: return fmt.Errorf("datastoreClient() failed: %v", err)
	}
	defer client.Close()

	var stateObj OAuthState
	err = client.Get(c, key, &stateObj)

	/*
		if err != nil {
			log.Errorf(c, "datastore.Get() failed: %v", err)
		}
	*/

	return err == nil
}

func RemoveOAuthState(c context.Context, state int64) error {
	client, err := dataStoreClient(c)
	if err != nil {
		return fmt.Errorf("datastoreClient() failed: %v", err)
	}
	defer client.Close()

	key := stateKey(c, state)
	return client.Delete(c, key)
}
