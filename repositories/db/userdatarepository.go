package db

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/datastore"
	domainuser "github.com/murraycu/go-bigoquiz-server/domain/user"
	dtouser "github.com/murraycu/go-bigoquiz-server/repositories/db/dtos/user"
	"github.com/murraycu/go-bigoquiz-server/server/loginserver/oauthparsers"
	"golang.org/x/oauth2"
	"google.golang.org/api/iterator"
)

const (
	// These are like database table names.
	DB_KIND_PROFILE     = "UserProfile"
	DB_KIND_USER_STATS  = "UserStats"
	DB_KIND_OAUTH_STATE = "OAuthState"
)

type UserDataRepository interface {
	GetUserProfileById(c context.Context, strUserId string) (*domainuser.Profile, error)
	GetUserStats(c context.Context, strUserId string) (map[string]*domainuser.Stats, error)
	GetUserStatsForQuiz(c context.Context, strUserId string, quizId string) (map[string]*domainuser.Stats, error)
	GetUserStatsForSection(c context.Context, strUserId string, quizId string, sectionId string) (*domainuser.Stats, error)
	StoreUserStats(c context.Context, userID string, stats *domainuser.Stats) error
	DeleteUserStatsForQuiz(c context.Context, strUserId string, quizId string) error

	StoreGoogleLoginInUserProfile(c context.Context, userInfo oauthparsers.GoogleUserInfo, strUserId string, token *oauth2.Token) (string, error)
	StoreGitHubLoginInUserProfile(c context.Context, userInfo oauthparsers.GitHubUserInfo, strUserId string, token *oauth2.Token) (string, error)
	StoreFacebookLoginInUserProfile(c context.Context, userInfo oauthparsers.FacebookUserInfo, strUserId string, token *oauth2.Token) (string, error)

	StoreGoogleTokenInUserProfile(c context.Context, userId string, token *oauth2.Token) error
	StoreGitHubTokenInUserProfile(c context.Context, userId string, token *oauth2.Token) error
	StoreFacebookTokenInUserProfile(c context.Context, userId string, token *oauth2.Token) error
}

type UserDataRepositoryImpl struct {
	client *datastore.Client
}

func NewUserDataRepository() (UserDataRepository, error) {
	result := &UserDataRepositoryImpl{}

	c := context.Background()
	var err error
	result.client, err = datastore.NewClient(c, "bigoquiz")
	if err != nil {
		return nil, fmt.Errorf("datastore.NewClient() failed: %v", err)
	}

	return result, nil
}

func (db *UserDataRepositoryImpl) getProfileFromDbQuery(c context.Context, q *datastore.Query) (*datastore.Key, *dtouser.Profile, error) {
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

func (db *UserDataRepositoryImpl) StoreGitHubLoginInUserProfile(c context.Context, userInfo oauthparsers.GitHubUserInfo, strUserId string, token *oauth2.Token) (string, error) {
	return storeOAuthLoginInUserProfile(db, c, userInfo, userInfo.Id, strUserId, token, db.getProfileFromDbByGitHubID, db.updateProfileFromGitHubUserInfo)
}

func (db *UserDataRepositoryImpl) StoreFacebookLoginInUserProfile(c context.Context, userInfo oauthparsers.FacebookUserInfo, strUserId string, token *oauth2.Token) (string, error) {
	return storeOAuthLoginInUserProfile(db, c, userInfo, userInfo.Id, strUserId, token, db.getProfileFromDbByFacebookID, db.updateProfileFromFacebookUserInfo)
}

func (db *UserDataRepositoryImpl) StoreGoogleLoginInUserProfile(c context.Context, userInfo oauthparsers.GoogleUserInfo, strUserId string, token *oauth2.Token) (string, error) {
	return storeOAuthLoginInUserProfile(db, c, userInfo, userInfo.Sub, strUserId, token, db.getProfileFromDbByGoogleID, db.updateProfileFromGoogleUserInfo)
}

func storeOAuthLoginInUserProfile[OAuthUserInfo any, ID any](db *UserDataRepositoryImpl, c context.Context, userInfo OAuthUserInfo, id ID, strUserId string, token *oauth2.Token, getProfileByOAuthId func(context.Context, ID) (*datastore.Key, *dtouser.Profile, error), updateProfile func(*dtouser.Profile, *OAuthUserInfo, *oauth2.Token) error) (string, error) {
	userIdFound, profile, err := getProfileByOAuthId(c, id)
	if err != nil {
		return "", fmt.Errorf("getProfileByOAuthId() failed: %v", err)
	}

	var userId *datastore.Key
	if userIdFound != nil {
		// Use the found user ID,
		// ignoring any user id from the caller.
		userId = userIdFound
	} else if len(strUserId) != 0 {
		userId, err = datastore.DecodeKey(strUserId)
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
		if err := updateProfile(profile, &userInfo, token); err != nil {
			return "", fmt.Errorf("updateProfile() failed (new profile): %v", err)
		}

		userId = datastore.IncompleteKey(DB_KIND_PROFILE, nil)
		if userId, err = db.client.Put(c, userId, profile); err != nil {
			return "", fmt.Errorf("datastore Put(with incomplete userId %v) failed: %v", userId, err)
		}
	} else if userId != nil {
		// Update the Profile:
		if err := updateProfile(profile, &userInfo, token); err != nil {
			return "", fmt.Errorf("updateProfile() failed: %v", err)
		}

		if userId, err = db.client.Put(c, userId, profile); err != nil {
			return "", fmt.Errorf("datastore Put(with userId %v) failed: %v", userId, err)
		}
	}

	return userId.Encode(), nil
}

func (db *UserDataRepositoryImpl) StoreGitHubTokenInUserProfile(c context.Context, userId string, token *oauth2.Token) error {
	return storeOAuthTokenInUserProfile(db, c, userId, token, db.updateProfileFromGitHubOAuthToken)
}

func (db *UserDataRepositoryImpl) StoreFacebookTokenInUserProfile(c context.Context, userId string, token *oauth2.Token) error {
	return storeOAuthTokenInUserProfile(db, c, userId, token, db.updateProfileFromFacebookOAuthToken)
}

func (db *UserDataRepositoryImpl) StoreGoogleTokenInUserProfile(c context.Context, userId string, token *oauth2.Token) error {
	return storeOAuthTokenInUserProfile(db, c, userId, token, db.updateProfileFromGoogleOAuthToken)
}

func storeOAuthTokenInUserProfile(db *UserDataRepositoryImpl, c context.Context, strUserId string, token *oauth2.Token, updateProfileFromOAuthToken func(*dtouser.Profile, *oauth2.Token) error) error {
	if len(strUserId) == 0 {
		return fmt.Errorf("storeOAuthTokenInUserProfile(): strUserId is empty")
	}

	userId, err := datastore.DecodeKey(strUserId)
	if err != nil {
		return fmt.Errorf("datastore.DecodeKey() failed: %v", err)
	}

	if userId == nil {
		return errors.New("userId is nil")
	}

	profile, err := db.getProfileFromDbByUserID(c, userId)
	if err != nil {
		return fmt.Errorf("getProfileFromDbByUserID() failed")
	}

	// Update the Profile:
	if err := updateProfileFromOAuthToken(profile, token); err != nil {
		return fmt.Errorf("updateProfileFromOAuthToken() failed: %v", err)
	}

	if userId, err = db.client.Put(c, userId, profile); err != nil {
		return fmt.Errorf("datastore Put(with userId %v) failed: %v", userId, err)
	}

	return nil
}

func (db *UserDataRepositoryImpl) getProfileFromDbByGitHubID(c context.Context, id int) (*datastore.Key, *dtouser.Profile, error) {
	q := datastore.NewQuery(DB_KIND_PROFILE).
		Filter("githubId =", id).
		Limit(1)
	return db.getProfileFromDbQuery(c, q)
}

func (db *UserDataRepositoryImpl) getProfileFromDbByFacebookID(c context.Context, id string) (*datastore.Key, *dtouser.Profile, error) {
	q := datastore.NewQuery(DB_KIND_PROFILE).
		Filter("facebookId =", id).
		Limit(1)
	return db.getProfileFromDbQuery(c, q)
}

func (db *UserDataRepositoryImpl) getProfileFromDbByGoogleID(c context.Context, sub string) (*datastore.Key, *dtouser.Profile, error) {
	q := datastore.NewQuery(DB_KIND_PROFILE).
		Filter("googleId =", sub).
		Limit(1)
	return db.getProfileFromDbQuery(c, q)
}

func (db *UserDataRepositoryImpl) getProfileFromDbByUserID(c context.Context, userId *datastore.Key) (*dtouser.Profile, error) {
	var profile dtouser.Profile
	err := db.client.Get(c, userId, &profile)
	if err != nil {
		// This is not an error.
		return nil, nil
	}

	return &profile, nil
}

func (db *UserDataRepositoryImpl) GetUserProfileById(c context.Context, strUserId string) (*domainuser.Profile, error) {
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

/** Get a map of stats by quiz ID, for all quizzes that have ever been used, from the database.
 * userId may be nil.
 */
func (db *UserDataRepositoryImpl) GetUserStats(c context.Context, strUserId string) (map[string]*domainuser.Stats, error) {
	userId, err := datastore.DecodeKey(strUserId)
	if err != nil {
		return nil, fmt.Errorf("datastore.DecodeKey() failed: %v", err)
	}

	var result = make(map[string]*domainuser.Stats)

	// In case a nil value could lead to getting all users' stats:
	if userId == nil {
		return result, nil
	}

	// Get all the Stats from the db, for each section, for each quiz:
	q := db.getQueryForUserStats(userId)

	iter := db.client.Run(c, q)

	if iter == nil {
		return nil, fmt.Errorf("datastore query for Stats failed")
	}

	// Build a map of the stats by quiz ID:
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
			// Start with this:
			existing = &domainuser.Stats{}
			existing.QuizId = quizId
		}

		combinedStats := createCombinedUserStatsWithoutQuestionHistories(existing, &stats)
		result[stats.QuizId] = combinedStats
	}

	return result, nil
}

/** Get a map of stats by section ID, for a specific quiz, from the database.
 * userId may be nil.
 * quizId may not be nil.
 */
func (db *UserDataRepositoryImpl) GetUserStatsForQuiz(c context.Context, strUserId string, quizId string) (map[string]*domainuser.Stats, error) {
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
func (db *UserDataRepositoryImpl) GetUserStatsForSection(c context.Context, strUserId string, quizId string, sectionId string) (*domainuser.Stats, error) {
	stats, err := db.getUserStatsForSectionAsDto(c, quizId, sectionId, strUserId)
	if err != nil {
		return nil, fmt.Errorf("getUserStatsForSectionAsDto() failed: %v", err)
	}

	if stats == nil {
		// This is not an error.
		// There are just no stats stored yet for this section.
		return nil, nil
	}

	return convertDtoStatsToDomainStats(stats), nil
}

func (db *UserDataRepositoryImpl) getUserStatsForSectionAsDto(c context.Context, quizId string, sectionId string, strUserId string) (*dtouser.Stats, error) {
	userId, err := datastore.DecodeKey(strUserId)
	if err != nil {
		return nil, fmt.Errorf("datastore.DecodeKey() failed: %v", err)
	}

	// Get the Stats from the db, for this section:
	// TODO: Remove duplicates if there is more than one?
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

	stats.Key = key
	// See the comment on user.Stats.Key
	return &stats, nil
}

func (db *UserDataRepositoryImpl) StoreUserStats(c context.Context, userID string, stats *domainuser.Stats) error {
	if len(stats.QuizId) == 0 {
		return fmt.Errorf("StoreUserStats(): QuizId is empty")
	}

	if len(stats.SectionId) == 0 {
		return fmt.Errorf("StoreUserStats(): SectionId is empty")
	}

	dtoOldStats, err := db.getUserStatsForSectionAsDto(c, stats.QuizId, stats.SectionId, userID)
	if err != nil {
		return fmt.Errorf("getUserStatsForSectionAsDto() failed: %v", err)
	}

	dtoStats, err := convertDomainStatsToDtoStats(stats, userID)
	if err != nil {
		return fmt.Errorf("convertDomainStatsToDtoStats() failed: %v", err)
	}

	// Use an existing key, to replace the existing entity, if one exists,
	// instead of just adding another one.
	if dtoOldStats != nil {
		dtoStats.Key = dtoOldStats.Key
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

func (db *UserDataRepositoryImpl) getQueryForUserStats(userId *datastore.Key) *datastore.Query {
	return datastore.NewQuery(DB_KIND_USER_STATS).
		Filter("userId =", userId)
}

func (db *UserDataRepositoryImpl) getQueryForUserStatsForQuiz(userId *datastore.Key, quizId string) *datastore.Query {
	return db.getQueryForUserStats(userId).
		Filter("quizId = ", quizId)
}

func (db *UserDataRepositoryImpl) DeleteUserStatsForQuiz(c context.Context, strUserId string, quizId string) error {
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
	q = q.KeysOnly()
	iter := db.client.Run(c, q)

	if iter == nil {
		return fmt.Errorf("datastore query for Stats failed")
	}

	for {
		key, err := iter.Next(nil)
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

		// Note that stats.key is nil, for some reason, but we have it from Next().

		// TODO: Batch these with datastore.DeleteMulti().
		err = db.client.Delete(c, key)
		if err != nil {
			return fmt.Errorf("datastore Delete() failed: %v", err)
		}
	}

	return nil
}
func (db *UserDataRepositoryImpl) updateProfileFromGoogleOAuthToken(profile *dtouser.Profile, token *oauth2.Token) error {
	profile.GoogleAccessToken = *token

	return nil
}

func (db *UserDataRepositoryImpl) updateProfileFromGoogleUserInfo(profile *dtouser.Profile, userInfo *oauthparsers.GoogleUserInfo, token *oauth2.Token) error {
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

	if err := db.updateProfileFromGoogleOAuthToken(profile, token); err != nil {
		return fmt.Errorf("failed to update Google OAuth token: %v", err)
	}

	profile.GoogleProfileUrl = userInfo.ProfileUrl

	return nil
}

func (db *UserDataRepositoryImpl) updateProfileFromGitHubOAuthToken(profile *dtouser.Profile, token *oauth2.Token) error {
	profile.GitHubAccessToken = *token

	return nil
}

func (db *UserDataRepositoryImpl) updateProfileFromGitHubUserInfo(profile *dtouser.Profile, userInfo *oauthparsers.GitHubUserInfo, token *oauth2.Token) error {
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

	if err := db.updateProfileFromGitHubOAuthToken(profile, token); err != nil {
		return fmt.Errorf("failed to update GitHub OAuth token: %v", err)
	}

	profile.GitHubProfileUrl = userInfo.ProfileUrl

	return nil
}

func (db *UserDataRepositoryImpl) updateProfileFromFacebookOAuthToken(profile *dtouser.Profile, token *oauth2.Token) error {
	profile.FacebookAccessToken = *token

	return nil
}

func (db *UserDataRepositoryImpl) updateProfileFromFacebookUserInfo(profile *dtouser.Profile, userInfo *oauthparsers.FacebookUserInfo, token *oauth2.Token) error {
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

	if err := db.updateProfileFromFacebookOAuthToken(profile, token); err != nil {
		return fmt.Errorf("failed to update Facebook OAuth token: %v", err)
	}

	profile.FacebookProfileUrl = userInfo.ProfileUrl

	return nil
}
