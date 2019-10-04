package db

import (
	"cloud.google.com/go/datastore"
	"fmt"
	"golang.org/x/net/context"
	"time"
)

type OAuthStateDataRepository struct {
	Client *datastore.Client
}

func NewOAuthStateDataRepository() (*OAuthStateDataRepository, error) {
	result := &OAuthStateDataRepository{}

	c := context.Background()
	var err error
	result.Client, err = datastore.NewClient(c, "bigoquiz")
	if err != nil {
		return nil, fmt.Errorf("datastore.NewClient() failed: %v", err)
	}

	return result, nil
}

/**
 */
type OAuthState struct {
	// The datastore ID (int64) is the oauth2 state.
	timestamp time.Time
}

func stateKey(state int64) *datastore.Key {
	return datastore.IDKey(DB_KIND_OAUTH_STATE, state, nil)
}

func (db *OAuthStateDataRepository) StoreOAuthState(c context.Context, state int64) error {
	key := stateKey(state)

	var stateObj OAuthState

	// Store a timestamp so a cron job can periodically remove old states.
	stateObj.timestamp = time.Now().UTC()

	_, err := db.Client.Put(c, key, &stateObj)
	if err != nil {
		return fmt.Errorf("datastore.Put() failed: %v", err)
	}

	return err
}

func (db *OAuthStateDataRepository) CheckOAuthState(c context.Context, state int64) error {
	key := stateKey(state)

	var stateObj OAuthState
	err := db.Client.Get(c, key, &stateObj)
	if err != nil {
		return fmt.Errorf("datastore Get() failed: %v", err)
	}

	return nil
}

func (db *OAuthStateDataRepository) RemoveOAuthState(c context.Context, state int64) error {
	key := stateKey(state)
	return db.Client.Delete(c, key)
}
