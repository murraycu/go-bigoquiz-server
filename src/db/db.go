package db

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"user"
	"golang.org/x/oauth2"
	"fmt"
)

const (
	// These are like database table names.
	DB_KIND_PROFILE = "UserProfile"
	DB_KIND_USER_STATS = "UserStats"
)

// Get the UserProfile via the GoogleID, adding it if necessary.
func StoreGoogleLoginInUserProfile(c context.Context, userInfo GoogleUserInfo, token *oauth2.Token) (*datastore.Key, error) {
	q := datastore.NewQuery(DB_KIND_PROFILE).
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

		key = datastore.NewIncompleteKey(c, DB_KIND_PROFILE, nil)
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
	err := datastore.Get(c, userId, &profile);
	if err == nil {
		return &profile, nil
	}

	// It's OK if no profile was found.
	// The caller can just create one.
	if err == datastore.ErrNoSuchEntity {
		return nil, nil
	}

	// Ignore errors caused by old fields in the datastore that are no longer mentioned in our Go struct.
	// TODO: The documentation does not clearly state that all matching fields will still be extracted.
	_, ok := err.(*datastore.ErrFieldMismatch)
	if ok {
		return &profile, nil
	}

	return nil, fmt.Errorf("datastore.Get() failed with key: %v: %v", userId, err)
}


// Get a map of stats by section ID, for all quizzes, from the database.
func GetUserStats(c context.Context, userId *datastore.Key) (map[string]*user.Stats, error) {
	// In case a nil value could lead to getting all users' stats:
	if userId == nil {
		return nil, fmt.Errorf("GetUserStatsForQuiz(): userId is nil.")
	}

	// Get all the Stats from the db, for each section:
	q := getQueryForUserStats(userId)
	iter := q.Run(c)

	if iter == nil {
		return nil, fmt.Errorf("datastore query for Stats failed.")
	}

	// Build a map of the stats by section ID:
	var result = make(map[string]*user.Stats)
	var stats user.Stats
	for {
		_, err := iter.Next(&stats)
		if err == datastore.Done {
			break;
		}

		if err != nil {
			return nil, fmt.Errorf("iter.Next() failed: %v", err)
		}

		result[stats.SectionId] = &stats
	}

	return result, nil
}

// Get a map of stats by section ID, for a specific quiz, from the database.
func GetUserStatsForQuiz(c context.Context, userId *datastore.Key, quizId string) (map[string]*user.Stats, error) {
	// In case a nil value could lead to getting all users' stats:
	if userId == nil {
		return nil, fmt.Errorf("GetUserStatsForQuiz(): userId is nil.")
	}

	// In case an empty value could lead to getting all quizzes' stats:
	if len(quizId) == 0 {
		return nil, fmt.Errorf("GetUserStatsForQuiz(): quizId is nil or empty.")
	}

	// Get all the Stats from the db, for each section:
	q := getQueryForUserStats(userId).
		Filter("QuizId = ", quizId)
	iter := q.Run(c)

	if iter == nil {
		return nil, fmt.Errorf("datastore query for Stats failed.")
	}

	// Build a map of the stats by section ID:
	var result = make(map[string]*user.Stats)
	var stats user.Stats
	for {
		_, err := iter.Next(&stats)
		if err == datastore.Done {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("iter.Next() failed: %v", err)
		}

		result[stats.SectionId] = &stats


	}

	return result, nil
}


// Get the stats for a specific section ID, from the database.
func GetUserStatsForSection(c context.Context, userId *datastore.Key, quizId string, sectionId string) (*user.Stats, error) {
	// Get all the Stats from the db, for each section:
	q := getQueryForUserStats(userId).
		Filter("SectionId =", sectionId).
	    Limit(1)
	iter := q.Run(c)

	if iter == nil {
		return nil, fmt.Errorf("datastore query for Stats failed.")
	}

	var stats user.Stats
	_, err := iter.Next(&stats)
	if err != nil && err != datastore.Done {
		return nil, fmt.Errorf("iter.Next() failed: %v", err)
	}

	return &stats, nil
}

func StoreUserStats(c context.Context, stats *user.Stats) error {
	if stats.Key == nil {
		// It hasn't been updated yet.
		stats.Key = datastore.NewIncompleteKey(c, DB_KIND_USER_STATS, nil)
	}

	if _, err := datastore.Put(c, stats.Key, stats); err != nil {
		return fmt.Errorf("StoreUserStats(): datastore.Put() failed: %v", err)
	}

	return nil;
}

func getQueryForUserStats(userId *datastore.Key) *datastore.Query {
	return datastore.NewQuery(DB_KIND_USER_STATS).
		Filter("Id =", userId)
}

func updateProfileFromGoogleUserInfo(profile *user.Profile, userInfo *GoogleUserInfo) {
	profile.GoogleId = userInfo.Sub
	profile.Name = userInfo.Name

	if userInfo.EmailVerified {
		profile.Email = userInfo.Email
	}
}
