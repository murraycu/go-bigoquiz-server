package db

import (
	"cloud.google.com/go/datastore"
	"fmt"
	domainuser "github.com/murraycu/go-bigoquiz-server/domain/user"
	dtouser "github.com/murraycu/go-bigoquiz-server/repositories/db/dtos/user"
	"github.com/murraycu/go-bigoquiz-server/server/loginserver/oauthparsers"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/iterator"
)

const (
	// These are like database table names.
	DB_KIND_PROFILE     = "UserProfile"
	DB_KIND_USER_STATS  = "UserStats"
	DB_KIND_OAUTH_STATE = "OAuthState"
)

type UserDataRepository struct {
	client *datastore.Client
}

func NewUserDataRepository() (*UserDataRepository, error) {
	result := &UserDataRepository{}

	c := context.Background()
	var err error
	result.client, err = datastore.NewClient(c, "bigoquiz")
	if err != nil {
		return nil, fmt.Errorf("datastore.NewClient() failed: %v", err)
	}

	return result, nil
}

func (db *UserDataRepository) getProfileFromDbQuery(c context.Context, q *datastore.Query) (*datastore.Key, *dtouser.Profile, error) {
	iter := db.client.Run(c, q)
	if iter == nil {
		return nil, nil, fmt.Errorf("datastore query for googleId failed")
	}

	var profile dtouser.Profile
	userId, err := iter.Next(&profile)
	if err == iterator.Done {
		// This is not an error.
		return nil, nil, nil
	} else if err != nil {
		// An unexpected error.
		return nil, nil, fmt.Errorf("datastore iter.Next() failed: %v", err)
	}

	return userId, &profile, nil
}

func (db *UserDataRepository) getProfileFromDbByGitHubID(c context.Context, id int) (*datastore.Key, *dtouser.Profile, error) {
	q := datastore.NewQuery(DB_KIND_PROFILE).
		Filter("githubId =", id).
		Limit(1)
	return db.getProfileFromDbQuery(c, q)
}

// Create, or update (if strUsedId is not empty) a user profile,
// using the provided Google login data.
// Returns the user ID.
func (db *UserDataRepository) StoreGitHubLoginInUserProfile(c context.Context, userInfo oauthparsers.GitHubUserInfo, strUserId string, token *oauth2.Token) (string, error) {
	userIdFound, profile, err := db.getProfileFromDbByGitHubID(c, userInfo.Id)
	if err != nil {
		// An unexpected error.
		return "", fmt.Errorf("getProfileFromDbByGitHubID() failed: %v", err)
	}

	var userId *datastore.Key
	if userIdFound != nil {
		// Use the found user ID,
		// ignoring any user id from the caller.
		userId = userIdFound
	} else if len(strUserId) != 0 {
		userId, err := datastore.DecodeKey(strUserId)
		if err != nil {
			return "", fmt.Errorf("datastore.DecodeKey() failed: %v", err)
		}

		// Try getting it via the supplied userID instead:
		profile, err = db.getProfileFromDbByUserID(c, userId)
		if err != nil {
			return "", fmt.Errorf("getProfileFromDbByUserID() failed")
		}
	}

	if profile == nil {
		// It is not in the datastore yet, so we add it.
		profile = new(dtouser.Profile)
		if err := db.updateProfileFromGitHubUserInfo(profile, &userInfo, token); err != nil {
			return "", fmt.Errorf("updateProfileFromGitHubUserInfo() failed (new profile): %v", err)
		}

		userId = datastore.IncompleteKey(DB_KIND_PROFILE, nil)
		if userId, err = db.client.Put(c, userId, profile); err != nil {
			return "", fmt.Errorf("datastore Put(with incomplete userId %v) failed: %v", userId, err)
		}
	} else if userId != nil {
		// Update the Profile:
		if err := db.updateProfileFromGitHubUserInfo(profile, &userInfo, token); err != nil {
			return "", fmt.Errorf("updateProfileFromGitHubUserInfo() failed: %v", err)
		}

		if userId, err = db.client.Put(c, userId, profile); err != nil {
			return "", fmt.Errorf("datastore Put(with userId %v) failed: %v", userId, err)
		}
	}

	return userId.Encode(), nil
}

func (db *UserDataRepository) getProfileFromDbByFacebookID(c context.Context, id string) (*datastore.Key, *dtouser.Profile, error) {
	q := datastore.NewQuery(DB_KIND_PROFILE).
		Filter("facebookId =", id).
		Limit(1)
	return db.getProfileFromDbQuery(c, q)
}

// Create, or update (if strUsedId is not empty) a user profile,
// using the provided Facebook login data.
// Returns the user ID.
func (db *UserDataRepository) StoreFacebookLoginInUserProfile(c context.Context, userInfo oauthparsers.FacebookUserInfo, strUserId string, token *oauth2.Token) (string, error) {
	userIdFound, profile, err := db.getProfileFromDbByFacebookID(c, userInfo.Id)
	if err != nil {
		// An unexpected error.
		return "", fmt.Errorf("getProfileFromDbByFacebookID() failed: %v", err)
	}

	var userId *datastore.Key
	if userIdFound != nil {
		// Use the found user ID,
		// ignoring any user id from the caller.
		userId = userIdFound
	} else if len(strUserId) != 0 {
		userId, err := datastore.DecodeKey(strUserId)
		if err != nil {
			return "", fmt.Errorf("datastore.DecodeKey() failed: %v", err)
		}

		// Try getting it via the supplied userID instead:
		profile, err = db.getProfileFromDbByUserID(c, userId)
		if err != nil {
			return "", fmt.Errorf("getProfileFromDbByUserID() failed")
		}
	}

	if profile == nil {
		// It is not in the datastore yet, so we add it.
		profile = new(dtouser.Profile)
		if err := db.updateProfileFromFacebookUserInfo(profile, &userInfo, token); err != nil {
			return "", fmt.Errorf("updateProfileFromFacebookUserInfo() failed (new profile): %v", err)
		}

		userId = datastore.IncompleteKey(DB_KIND_PROFILE, nil)
		if userId, err = db.client.Put(c, userId, profile); err != nil {
			return "", fmt.Errorf("datastore Put(with incomplete userId %v) failed: %v", userId, err)
		}
	} else if userId != nil {
		// Update the Profile:
		if err := db.updateProfileFromFacebookUserInfo(profile, &userInfo, token); err != nil {
			return "", fmt.Errorf("updateProfileFromFacebookUserInfo() failed: %v", err)
		}

		if userId, err = db.client.Put(c, userId, profile); err != nil {
			return "", fmt.Errorf("datastore Put(with userId %v) failed: %v", userId, err)
		}
	}

	return userId.Encode(), nil
}

func (db *UserDataRepository) getProfileFromDbByGoogleID(c context.Context, sub string) (*datastore.Key, *dtouser.Profile, error) {
	q := datastore.NewQuery(DB_KIND_PROFILE).
		Filter("googleId =", sub).
		Limit(1)
	return db.getProfileFromDbQuery(c, q)
}

func (db *UserDataRepository) getProfileFromDbByUserID(c context.Context, userId *datastore.Key) (*dtouser.Profile, error) {
	var profile dtouser.Profile
	err := db.client.Get(c, userId, &profile)
	if err != nil {
		// This is not an error.
		return nil, nil
	}

	return &profile, nil
}

// Create, or update (if strUsedId is not empty) a user profile,
// using the provided Google login data.
// Returns the user ID.
//
// TODO: Make this function generic, parameterizing on GoogleUserInfo/GithubUserInfo,
// if Go ever has generics.
// Get the UserProfile via the GoogleID, adding it if necessary.
func (db *UserDataRepository) StoreGoogleLoginInUserProfile(c context.Context, userInfo oauthparsers.GoogleUserInfo, strUserId string, token *oauth2.Token) (string, error) {
	userIdFound, profile, err := db.getProfileFromDbByGoogleID(c, userInfo.Sub)
	if err != nil {
		// An unexpected error.
		return "", fmt.Errorf("getProfileFromDbByGoogleID() failed: %v", err)
	}

	var userId *datastore.Key
	if userIdFound != nil {
		// Use the found user ID,
		// ignoring any user id from the caller.
		userId = userIdFound
	} else if len(strUserId) != 0 {
		// Try getting it via the supplied userID instead:
		userId, err := datastore.DecodeKey(strUserId)
		if err != nil {
			return "", fmt.Errorf("datastore.DecodeKey() failed for key: %v: %v", strUserId, err)
		}

		profile, err = db.getProfileFromDbByUserID(c, userId)
		if err != nil {
			return "", fmt.Errorf("getProfileFromDbByUserID() failed")
		}
	}

	if profile == nil {
		// It is not in the datastore yet, so we add it.
		profile = new(dtouser.Profile)
		if err := db.updateProfileFromGoogleUserInfo(profile, &userInfo, token); err != nil {
			return "", fmt.Errorf("updateProfileFromGoogleUserInfo() failed (new profile): %v", err)
		}

		userId = datastore.IncompleteKey(DB_KIND_PROFILE, nil)
		if userId, err = db.client.Put(c, userId, profile); err != nil {
			return "", fmt.Errorf("datastore. ut(with incomplete userId %v) failed: %v", userId, err)
		}
	} else if userId != nil {
		// Update the Profile:
		if err := db.updateProfileFromGoogleUserInfo(profile, &userInfo, token); err != nil {
			return "", fmt.Errorf("updateProfileFromGoogleUserInfo() failed: %v", err)
		}

		if userId, err = db.client.Put(c, userId, profile); err != nil {
			return "", fmt.Errorf("datastore Put(with userId %v) failed: %v", userId, err)
		}
	}

	return userId.Encode(), nil
}

func (db *UserDataRepository) GetUserProfileById(c context.Context, strUserId string) (*domainuser.Profile, error) {
	userId, err := datastore.DecodeKey(strUserId)
	if err != nil {
		return nil, fmt.Errorf("datastore.DecodeKey() failed: %v", err)
	}

	var profile dtouser.Profile
	err = db.client.Get(c, userId, &profile)
	if err == nil {
		return convertDtoProfileToDomainProfile(&profile), nil
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
		return convertDtoProfileToDomainProfile(&profile), nil
	}

	return nil, fmt.Errorf("datastore.Get() failed with key: %v: %v", userId, err)
}

/** Add the values from userStat to this instance,
* returning a combined UserStats,
* ignoring the question histories,
* without changing this instance.
 */
func createCombinedUserStatsWithoutQuestionHistories(self *domainuser.Stats, stats *dtouser.Stats) *domainuser.Stats {
	if stats == nil {
		return self
	}

	var result domainuser.Stats
	result.QuizId = self.QuizId

	result.Answered = self.Answered + stats.Answered
	result.Correct = self.Correct + stats.Correct

	result.CountQuestionsAnsweredOnce = self.CountQuestionsAnsweredOnce + stats.CountQuestionsAnsweredOnce
	result.CountQuestionsCorrectOnce = self.CountQuestionsCorrectOnce + stats.CountQuestionsCorrectOnce

	return &result
}

/** Get a map of stats by quiz ID, for all quizzes, from the database.
 * userId may be nil.
 */
func (db *UserDataRepository) GetUserStats(c context.Context, strUserId string) (map[string]*domainuser.Stats, error) {
	userId, err := datastore.DecodeKey(strUserId)
	if err != nil {
		return nil, fmt.Errorf("datastore.DecodeKey() failed: %v", err)
	}

	var result = make(map[string]*domainuser.Stats)

	// In case a nil value could lead to getting all users' stats:
	if userId == nil {
		return result, nil
	}

	// Get all the Stats from the db, for each section:
	q := db.getQueryForUserStats(userId)

	iter := db.client.Run(c, q)

	if iter == nil {
		return nil, fmt.Errorf("datastore query for Stats failed")
	}

	// Build a map of the stats by section ID:
	var stats dtouser.Stats
	for {
		_, err := iter.Next(&stats)
		if err == iterator.Done {
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
			result[quizId] = convertDtoStatsToDomainStats(&stats)
		} else {
			combinedStats := createCombinedUserStatsWithoutQuestionHistories(existing, &stats)
			result[stats.QuizId] = combinedStats
		}

		// This does not correspond to a dtouser.Stats in the datastore.
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
func (db *UserDataRepository) GetUserStatsForQuiz(c context.Context, strUserId string, quizId string) (map[string]*domainuser.Stats, error) {
	userId, err := datastore.DecodeKey(strUserId)
	if err != nil {
		return nil, fmt.Errorf("datastore.DecodeKey() failed: %v", err)
	}

	var result = make(map[string]*domainuser.Stats)

	// In case a nil value could lead to getting all users' stats:
	if userId == nil {
		return result, nil
	}

	// In case an empty value could lead to getting all quizzes' stats:
	if len(quizId) == 0 {
		return nil, fmt.Errorf("GetUserStatsForQuiz(): quizId is nil or empty")
	}

	// Get all the Stats from the db, for each section:
	q := db.getQueryForUserStatsForQuiz(userId, quizId)

	iter := db.client.Run(c, q)

	if iter == nil {
		return nil, fmt.Errorf("datastore query for Stats failed")
	}

	// Build a map of the stats by section ID:
	for {
		var stats dtouser.Stats
		key, err := iter.Next(&stats)
		if err == iterator.Done {
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
		result[stats.SectionId] = convertDtoStatsToDomainStats(&stats)
	}

	return result, nil
}

// Get the stats for a specific section ID, from the database.
func (db *UserDataRepository) GetUserStatsForSection(c context.Context, strUserId string, quizId string, sectionId string) (*domainuser.Stats, error) {
	userId, err := datastore.DecodeKey(strUserId)
	if err != nil {
		return nil, fmt.Errorf("datastore.DecodeKey() failed: %v", err)
	}

	// Get all the Stats from the db, for each section:
	q := db.getQueryForUserStatsForQuiz(userId, quizId).
		Filter("sectionId =", sectionId).
		Limit(1)

	iter := db.client.Run(c, q)

	if iter == nil {
		return nil, fmt.Errorf("datastore query for Stats failed")
	}

	var stats dtouser.Stats
	key, err := iter.Next(&stats)
	if err != nil {
		if err == iterator.Done {
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
	return convertDtoStatsToDomainStats(&stats), nil
}

func (db *UserDataRepository) StoreUserStats(c context.Context, userID string, stats *domainuser.Stats) error {
	if len(stats.QuizId) == 0 {
		return fmt.Errorf("StoreUserStats(): QuizId is empty")
	}

	if len(stats.SectionId) == 0 {
		return fmt.Errorf("StoreUserStats(): SectionId is empty")
	}

	dtoStats, err := convertDomainStatsToDtoStats(stats, userID)
	if err != nil {
		return fmt.Errorf("convertDomainStatsToDtoStats() failed: %v", err)
	}

	key := dtoStats.Key
	if key == nil {
		// It hasn't been updated yet.
		//
		// Note: Don't store the key in stats.Key - that confuses the datastore API,
		// (but without any error being returned to our code.)
		// so we won't be able to read the entity back later.
		// That also results in an error when trying to list the UserStats entities in dev_server.py's
		// Datastore Viewer:
		// "in ValidatePropertyKey 'Incomplete key found for reference property %s.' % name)
		// BadValueError: Incomplete key found for reference property Key."
		key = datastore.IncompleteKey(DB_KIND_USER_STATS, nil)
	}

	if key, err = db.client.Put(c, key, dtoStats); err != nil {
		return fmt.Errorf("StoreUserStats(): datastore Put() failed: %v", err)
	}

	// TODO: stats.UserId = key // See the comment on Stats.Key.

	return nil
}

func (db *UserDataRepository) getQueryForUserStats(userId *datastore.Key) *datastore.Query {
	return datastore.NewQuery(DB_KIND_USER_STATS).
		Filter("userId =", userId)
}

func (db *UserDataRepository) getQueryForUserStatsForQuiz(userId *datastore.Key, quizId string) *datastore.Query {
	return db.getQueryForUserStats(userId).
		Filter("quizId = ", quizId)
}

func (db *UserDataRepository) DeleteUserStatsForQuiz(c context.Context, strUserId string, quizId string) error {
	userId, err := datastore.DecodeKey(strUserId)
	if err != nil {
		return fmt.Errorf("datastore.DecodeKey() failed: %v", err)
	}

	// In case a nil value could lead to deleting all users' stats:
	if userId == nil {
		return fmt.Errorf("DeleteUserStatsForQuiz(): userId is nil")
	}

	// In case an empty value could lead to deleting all quizzes' stats:
	if len(quizId) == 0 {
		return fmt.Errorf("DeleteUserStatsForQuiz(): quizId is nil or empty")
	}

	q := db.getQueryForUserStatsForQuiz(userId, quizId)
	iter := db.client.Run(c, q)

	if iter == nil {
		return fmt.Errorf("datastore query for Stats failed")
	}

	var stats dtouser.Stats
	for {
		_, err := iter.Next(&stats)
		if err == iterator.Done {
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

		if stats.Key == nil {
			return fmt.Errorf("The retrieved Stats's key is nil: %v", err)
		}

		// TODO: Batch these with datastore.DeleteMulti().
		err = db.client.Delete(c, stats.Key)
		if err != nil {
			return fmt.Errorf("datastore Delete() failed: %v", err)
		}
	}

	return nil
}

func (db *UserDataRepository) updateProfileFromGoogleUserInfo(profile *dtouser.Profile, userInfo *oauthparsers.GoogleUserInfo, token *oauth2.Token) error {
	if profile == nil {
		return fmt.Errorf("profile is nil")
	}

	if userInfo == nil {
		return fmt.Errorf("userInfo is nil")
	}

	if token == nil {
		return fmt.Errorf("token is nil")
	}

	profile.GoogleId = userInfo.Sub
	profile.Name = userInfo.Name

	if userInfo.EmailVerified {
		profile.Email = userInfo.Email
	}

	profile.GoogleAccessToken = *token
	profile.GoogleProfileUrl = userInfo.ProfileUrl

	return nil
}

func (db *UserDataRepository) updateProfileFromGitHubUserInfo(profile *dtouser.Profile, userInfo *oauthparsers.GitHubUserInfo, token *oauth2.Token) error {
	if profile == nil {
		return fmt.Errorf("profile is nil")
	}

	if userInfo == nil {
		return fmt.Errorf("userInfo is nil")
	}

	if token == nil {
		return fmt.Errorf("token is nil")
	}

	profile.GitHubId = userInfo.Id
	profile.Name = userInfo.Name
	// TODO: Get a verified email address, to compare with the other account?
	profile.GitHubAccessToken = *token
	profile.GitHubProfileUrl = userInfo.ProfileUrl

	return nil
}

func (db *UserDataRepository) updateProfileFromFacebookUserInfo(profile *dtouser.Profile, userInfo *oauthparsers.FacebookUserInfo, token *oauth2.Token) error {
	if profile == nil {
		return fmt.Errorf("profile is nil")
	}

	if userInfo == nil {
		return fmt.Errorf("userInfo is nil")
	}

	if token == nil {
		return fmt.Errorf("token is nil")
	}

	profile.FacebookId = userInfo.Id
	profile.Name = userInfo.Name
	// TODO: Get a verified email address, to compare with the other account?
	profile.FacebookAccessToken = *token
	profile.FacebookProfileUrl = userInfo.ProfileUrl

	return nil
}
