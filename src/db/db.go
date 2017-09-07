package db

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"user"
	"golang.org/x/oauth2"
	"fmt"
)

// Get the UserProfile via the GoogleID, adding it if necessary.
func StoreGoogleLoginInUserProfile(c context.Context, userInfo GoogleUserInfo, token *oauth2.Token) (*datastore.Key, error) {
	q := datastore.NewQuery("user.Profile").
		Filter("GoogleId =", userInfo.Sub).
		Limit(1)
	iter := q.Run(c)
	if iter == nil {
		return nil, fmt.Errorf("datastore query for GoogleId failed.")
	}

	var profile user.Profile
	var key *datastore.Key
	var err error
	key, err = iter.Next(&profile)
	if err == datastore.Done {
		// It is not in the datastore yet, so we add it.
		updateProfileFromGoogleUserInfo(&profile, &userInfo)
		profile.GoogleAccessToken = *token

		key = datastore.NewIncompleteKey(c, "user.Profile", nil)
		if key, err = datastore.Put(c, key, &profile); err != nil {
			return nil, fmt.Errorf("datastore.Put(with incomplete key %v) failed: %v", key, err)
		}
	} else if err != nil {
		// An unexpected error.
		return nil, err
	} else {
		// Update the Profile:
		updateProfileFromGoogleUserInfo(&profile, &userInfo)
		profile.GoogleAccessToken = *token

		if key, err = datastore.Put(c, key, &profile); err != nil {
			return nil, fmt.Errorf("datastore.Put(with key %v) failed: %v", key, err)
		}
	}

	return key, nil
}

func GetUserProfileById(c context.Context, userId *datastore.Key) (*user.Profile, error) {
	var profile user.Profile
	if err := datastore.Get(c, userId, &profile); err != nil {
		return nil, fmt.Errorf("datastore.Get() failed with key: %v: %v", userId, err)
	}

	return &profile, nil
}

func updateProfileFromGoogleUserInfo(profile *user.Profile, userInfo *GoogleUserInfo) {
	profile.GoogleId = userInfo.Sub
	profile.Name = userInfo.Name

	if userInfo.EmailVerified {
		profile.Email = userInfo.Email
	}
}
