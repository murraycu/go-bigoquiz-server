package db

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/appengine/datastore"
	"user"
)

const (
	// These are like database table names.
	DB_KIND_PROFILE    = "UserProfile"
	DB_KIND_USER_STATS = "UserStats"
	DB_KIND_OAUTH_STATE = "OAuthState"
)

func getProfileFromDbQuery(c context.Context, q *datastore.Query) (*datastore.Key, *user.Profile, error) {
	iter := q.Run(c)
	if iter == nil {
		return nil, nil, fmt.Errorf("datastore query for googleId failed")
	}

	var profile user.Profile
	userId, err := iter.Next(&profile)
	if err == datastore.Done {
		// This is not an error.
		return nil, nil, nil
	} else if err != nil {
		// An unexpected error.
		return nil, nil, fmt.Errorf("datastore iter.Next() failed: %v", err)
	}

	return userId, &profile, nil
}

func getProfileFromDbByGitHubID(c context.Context, id int) (*datastore.Key, *user.Profile, error) {
	q := datastore.NewQuery(DB_KIND_PROFILE).
		Filter("githubId =", id).
		Limit(1)
	return getProfileFromDbQuery(c, q)
}

func StoreGitHubLoginInUserProfile(c context.Context, userInfo GitHubUserInfo, userId *datastore.Key, token *oauth2.Token) (*datastore.Key, error) {
	userIdFound, profile, err := getProfileFromDbByGitHubID(c, userInfo.Id)
	if err != nil {
		// An unexpected error.
		return nil, fmt.Errorf("getProfileFromDbByGitHubID() failed: %v", err)
	}

	if userIdFound != nil {
		// Use the found user ID,
		// ignoring any user id from the caller.
		userId = userIdFound
	} else if userId != nil {
		// Try getting it via the supplied userID instead:
		profile, err = getProfileFromDbByUserID(c, userId)
		if err != nil {
			return nil, fmt.Errorf("getProfileFromDbByUserID() failed")
		}
	}

	if profile == nil {
		// It is not in the datastore yet, so we add it.
		updateProfileFromGitHubUserInfo(profile, &userInfo, token)

		userId = datastore.NewIncompleteKey(c, DB_KIND_PROFILE, nil)
		if userId, err = datastore.Put(c, userId, profile); err != nil {
			return nil, fmt.Errorf("datastore.Put(with incomplete userId %v) failed: %v", userId, err)
		}
	} else if userId != nil {
		// Update the Profile:
		updateProfileFromGitHubUserInfo(profile, &userInfo, token)

		if userId, err = datastore.Put(c, userId, profile); err != nil {
			return nil, fmt.Errorf("datastore.Put(with userId %v) failed: %v", userId, err)
		}
	}

	return userId, nil
}

func getProfileFromDbByFacebookID(c context.Context, id string) (*datastore.Key, *user.Profile, error) {
	q := datastore.NewQuery(DB_KIND_PROFILE).
		Filter("facebookId =", id).
		Limit(1)
	return getProfileFromDbQuery(c, q)
}

func StoreFacebookLoginInUserProfile(c context.Context, userInfo FacebookUserInfo, userId *datastore.Key, token *oauth2.Token) (*datastore.Key, error) {
	userIdFound, profile, err := getProfileFromDbByFacebookID(c, userInfo.Id)
	if err != nil {
		// An unexpected error.
		return nil, fmt.Errorf("getProfileFromDbByFacebookID() failed: %v", err)
	}

	if userIdFound != nil {
		// Use the found user ID,
		// ignoring any user id from the caller.
		userId = userIdFound
	} else if userId != nil {
		// Try getting it via the supplied userID instead:
		profile, err = getProfileFromDbByUserID(c, userId)
		if err != nil {
			return nil, fmt.Errorf("getProfileFromDbByUserID() failed")
		}
	}

	if profile == nil {
		// It is not in the datastore yet, so we add it.
		profile = new(user.Profile)
		updateProfileFromFacebookUserInfo(profile, &userInfo, token)

		userId = datastore.NewIncompleteKey(c, DB_KIND_PROFILE, nil)
		if userId, err = datastore.Put(c, userId, profile); err != nil {
			return nil, fmt.Errorf("datastore.Put(with incomplete userId %v) failed: %v", userId, err)
		}
	} else if userId != nil {
		// Update the Profile:
		updateProfileFromFacebookUserInfo(profile, &userInfo, token)

		if userId, err = datastore.Put(c, userId, profile); err != nil {
			return nil, fmt.Errorf("datastore.Put(with userId %v) failed: %v", userId, err)
		}
	}

	return userId, nil
}

func getProfileFromDbByGoogleID(c context.Context, sub string) (*datastore.Key, *user.Profile, error) {
	q := datastore.NewQuery(DB_KIND_PROFILE).
		Filter("googleId =", sub).
		Limit(1)
	return getProfileFromDbQuery(c, q)
}

func getProfileFromDbByUserID(c context.Context, userId *datastore.Key) (*user.Profile, error) {
	var profile user.Profile
	err := datastore.Get(c, userId, &profile)
	if err != nil {
		// This is not an error.
		return nil, nil
	}

	return &profile, nil
}

// TODO: Make this function generic, parameterizing on GoogleUserInfo/GithubUserInfo,
// if Go ever has generics.
// Get the UserProfile via the GoogleID, adding it if necessary.
func StoreGoogleLoginInUserProfile(c context.Context, userInfo GoogleUserInfo, userId *datastore.Key, token *oauth2.Token) (*datastore.Key, error) {
	userIdFound, profile, err := getProfileFromDbByGoogleID(c, userInfo.Sub)
	if err != nil {
		// An unexpected error.
		return nil, fmt.Errorf("getProfileFromDbByGoogleID() failed: %v", err)
	}

	if userIdFound != nil {
		// Use the found user ID,
		// ignoring any user id from the caller.
		userId = userIdFound
	} else if userId != nil {
		// Try getting it via the supplied userID instead:
		profile, err = getProfileFromDbByUserID(c, userId)
		if err != nil {
			return nil, fmt.Errorf("getProfileFromDbByUserID() failed")
		}
	}

	if profile == nil {
		// It is not in the datastore yet, so we add it.
		profile = new(user.Profile)
		updateProfileFromGoogleUserInfo(profile, &userInfo, token)

		userId = datastore.NewIncompleteKey(c, DB_KIND_PROFILE, nil)
		if userId, err = datastore.Put(c, userId, profile); err != nil {
			return nil, fmt.Errorf("datastore.Put(with incomplete userId %v) failed: %v", userId, err)
		}
	} else if userId != nil {
		// Update the Profile:
		updateProfileFromGoogleUserInfo(profile, &userInfo, token)

		if userId, err = datastore.Put(c, userId, profile); err != nil {
			return nil, fmt.Errorf("datastore.Put(with userId %v) failed: %v", userId, err)
		}
	}

	return userId, nil
}

func GetUserProfileById(c context.Context, userId *datastore.Key) (*user.Profile, error) {
	var profile user.Profile
	err := datastore.Get(c, userId, &profile)
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

/** Get a map of stats by quiz ID, for all quizzes, from the database.
 * userId may be nil.
 */
func GetUserStats(c context.Context, userId *datastore.Key) (map[string]*user.Stats, error) {
	var result = make(map[string]*user.Stats)

	// In case a nil value could lead to getting all users' stats:
	if userId == nil {
		return result, nil
	}

	// Get all the Stats from the db, for each section:
	q := getQueryForUserStats(userId)
	iter := q.Run(c)

	if iter == nil {
		return nil, fmt.Errorf("datastore query for Stats failed")
	}

	// Build a map of the stats by section ID:
	var stats user.Stats
	for {
		_, err := iter.Next(&stats)
		if err == datastore.Done {
			break
		}

		if err != nil {
			if _, ok := err.(*datastore.ErrFieldMismatch); ok {
				// Ignore these errors during development,
				// TODO: Remove this for production,
				// because it then gives us a Stats instance in an unpredictable state.
				continue
			}

			return nil, fmt.Errorf("iter.Next() failed: %v", err)
		}

		quizId := stats.QuizId

		existing, exists := result[quizId]
		if !exists {
			result[quizId] = &stats
		} else {
			combinedStats := existing.CreateCombinedUserStatsWithoutQuestionHistories(&stats)
			result[stats.QuizId] = combinedStats
		}

		// This does not correspond to a user.Stats in the datastore.
		// Instead this one is for the whole quiz, not just a section in a quiz.
		// So we wipe the Key to make sure that we don't try to write it back sometime.
		stats.Key = nil // See the comment on Stats.Key.

	}

	return result, nil
}

/** Get a map of stats by section ID, for a specific quiz, from the database.
 * userId may be nil.
 * quizId may not be nil.
 */
func GetUserStatsForQuiz(c context.Context, userId *datastore.Key, quizId string) (map[string]*user.Stats, error) {
	var result = make(map[string]*user.Stats)

	// In case a nil value could lead to getting all users' stats:
	if userId == nil {
		return result, nil
	}

	// In case an empty value could lead to getting all quizzes' stats:
	if len(quizId) == 0 {
		return nil, fmt.Errorf("GetUserStatsForQuiz(): quizId is nil or empty")
	}

	// Get all the Stats from the db, for each section:
	q := getQueryForUserStatsForQuiz(userId, quizId)
	iter := q.Run(c)

	if iter == nil {
		return nil, fmt.Errorf("datastore query for Stats failed")
	}

	// Build a map of the stats by section ID:
	for {
		var stats user.Stats
		key, err := iter.Next(&stats)
		if err == datastore.Done {
			break
		}

		if err != nil {
			if _, ok := err.(*datastore.ErrFieldMismatch); ok {
				// Ignore these errors during development,
				// TODO: Remove this for production,
				// because it then gives us a Stats instance in an unpredictable state.
				continue
			}

			return nil, fmt.Errorf("iter.Next() failed: %v", err)
		}

		stats.Key = key // See the comment on user.Stats.Key
		result[stats.SectionId] = &stats
	}

	return result, nil
}

// Get the stats for a specific section ID, from the database.
func GetUserStatsForSection(c context.Context, userId *datastore.Key, quizId string, sectionId string) (*user.Stats, error) {
	// Get all the Stats from the db, for each section:
	q := getQueryForUserStatsForQuiz(userId, quizId).
		Filter("sectionId =", sectionId).
		Limit(1)
	iter := q.Run(c)

	if iter == nil {
		return nil, fmt.Errorf("datastore query for Stats failed")
	}

	var stats user.Stats
	key, err := iter.Next(&stats)
	if err != nil {
		if err == datastore.Done {
			// It was not found.
			return nil, nil
		} else {
			if _, ok := err.(*datastore.ErrFieldMismatch); ok {
				// Ignore these errors during development,
				// TODO: Remove this for production,
				// because it then gives us a Stats instance in an unpredictable state.
			} else {
				return nil, fmt.Errorf("iter.Next() failed: %v", err)
			}
		}
	}

	stats.Key = key // See the comment on user.Stats.Key
	return &stats, nil
}

func StoreUserStats(c context.Context, stats *user.Stats) error {
	if len(stats.QuizId) == 0 {
		return fmt.Errorf("StoreUserStats(): QuizId is empty")
	}

	if len(stats.SectionId) == 0 {
		return fmt.Errorf("StoreUserStats(): SectionId is empty")
	}

	key := stats.Key
	if key == nil {
		// It hasn't been updated yet.
		//
		// Not: Don't store the key in stats.Key - that confuses the datastore API,
		// (but without any error being returned to our code.)
		// so we won't be able to read the entity back later.
		// That also results in an error when trying to list the UserStats entities in dev_server.py's
		// Datastore Viewer:
		// "in ValidatePropertyKey 'Incomplete key found for reference property %s.' % name)
		// BadValueError: Incomplete key found for reference property Key."
		key = datastore.NewIncompleteKey(c, DB_KIND_USER_STATS, nil)
	}

	var err error
	if key, err = datastore.Put(c, key, stats); err != nil {
		return fmt.Errorf("StoreUserStats(): datastore.Put() failed: %v", err)
	}

	stats.Key = key // See the comment on Stats.Key.

	return nil
}

func getQueryForUserStats(userId *datastore.Key) *datastore.Query {
	return datastore.NewQuery(DB_KIND_USER_STATS).
		Filter("userId =", userId)
}

func getQueryForUserStatsForQuiz(userId *datastore.Key, quizId string) *datastore.Query {
	return getQueryForUserStats(userId).
		Filter("quizId = ", quizId)
}

func DeleteUserStatsForQuiz(c context.Context, userId *datastore.Key, quizId string) error {
	// In case a nil value could lead to deleting all users' stats:
	if userId == nil {
		return fmt.Errorf("DeleteUserStatsForQuiz(): userId is nil")
	}

	// In case an empty value could lead to deleting all quizzes' stats:
	if len(quizId) == 0 {
		return fmt.Errorf("DeleteUserStatsForQuiz(): quizId is nil or empty")
	}

	q := getQueryForUserStatsForQuiz(userId, quizId)
	iter := q.Run(c)

	if iter == nil {
		return fmt.Errorf("datastore query for Stats failed")
	}

	var stats user.Stats
	for {
		_, err := iter.Next(&stats)
		if err == datastore.Done {
			break
		}

		if err != nil {
			if _, ok := err.(*datastore.ErrFieldMismatch); ok {
				// Ignore these errors during development,
				// TODO: Remove this for production,
				// because it then gives us a Stats instance in an unpredictable state.
				continue
			}

			return fmt.Errorf("iter.Next() failed: %v", err)
		}

		// TODO: Batch these with datastore.DeleteMulti().
		err = datastore.Delete(c, stats.Key)
		if err != nil {
			return fmt.Errorf("datastore.Delete() failed: %v", err)
		}
	}

	return nil
}

func updateProfileFromGoogleUserInfo(profile *user.Profile, userInfo *GoogleUserInfo, token *oauth2.Token) {
	profile.GoogleId = userInfo.Sub
	profile.Name = userInfo.Name

	if userInfo.EmailVerified {
		profile.Email = userInfo.Email
	}

	profile.GoogleAccessToken = *token
	profile.GoogleProfileUrl = userInfo.ProfileUrl
}

func updateProfileFromGitHubUserInfo(profile *user.Profile, userInfo *GitHubUserInfo, token *oauth2.Token) {
	profile.GitHubId = userInfo.Id
	profile.Name = userInfo.Name
	// TODO: Get a verified email address, to compare with the other account?
	profile.GitHubAccessToken = *token
	profile.GitHubProfileUrl = userInfo.ProfileUrl
}

func updateProfileFromFacebookUserInfo(profile *user.Profile, userInfo *FacebookUserInfo, token *oauth2.Token) {
	profile.FacebookId = userInfo.Id
	profile.Name = userInfo.Name
	// TODO: Get a verified email address, to compare with the other account?
	profile.FacebookAccessToken = *token
	profile.FacebookProfileUrl = userInfo.ProfileUrl
}
