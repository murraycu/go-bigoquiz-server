package db

import (
	"context"
	"github.com/murraycu/go-bigoquiz-server/server/loginserver/oauthparsers"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"testing"
)

func TestNewRestServerInstantiate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test which requires more setup.")
	}

	userDataClient, err := NewUserDataRepository()
	assert.Nil(t, err)
	assert.NotNil(t, userDataClient)
}

func TestNewRestServerGetUserProfileByIdForNonExistantUser(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test which requires more setup.")
	}

	userDataClient, err := NewUserDataRepository()
	assert.Nil(t, err)
	assert.NotNil(t, userDataClient)

	c := context.Background()

	// This must be decodable with datastore.DecodeKey().
	userId := "EhYKC1VzZXJQcm9maWxlEICAgICw2IIK"

	userProfile, err := userDataClient.GetUserProfileById(c, userId)
	assert.Nil(t, err)
	assert.Nil(t, userProfile)
}

func TestNewRestServerStoreAndGetUserProfileById(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test which requires more setup.")
	}

	userDataClient, err := NewUserDataRepository()
	assert.Nil(t, err)
	assert.NotNil(t, userDataClient)

	c := context.Background()

	// This must be decodable with datastore.DecodeKey().
	userId := "EhYKC1VzZXJQcm9maWxlEICAgICw2IIL"

	userInfo := oauthparsers.GoogleUserInfo{
		Sub: "some-google-user-id",
		Email: "example@example.com",
		EmailVerified: true,
		Name: "Example McExample",
	}
	strUserId := "" // Create a new user.
	token  := oauth2.Token{
		AccessToken: "some-access-token",
	}
	userId, err = userDataClient.StoreGoogleLoginInUserProfile(c, userInfo, strUserId, &token)
	assert.Nil(t, err)
	assert.NotNil(t, userId)

	userProfile, err := userDataClient.GetUserProfileById(c, userId)
	assert.Nil(t, err)
	assert.NotNil(t, userProfile)

	assert.Equal(t, userInfo.Email, userProfile.Email)
	assert.Equal(t, userInfo.Name, userProfile.Name)
}