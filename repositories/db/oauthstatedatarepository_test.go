package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOAuthStateDataRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test which requires more setup.")
	}

	oauthStateDataRepository, err := NewOAuthStateDataRepository()
	assert.Nil(t, err)
	assert.NotNil(t, oauthStateDataRepository)
}

func TestOAuthStateDataRepositorySetAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test which requires more setup.")
	}

	oauthStateDataRepository, err := NewOAuthStateDataRepository()
	assert.Nil(t, err)
	assert.NotNil(t, oauthStateDataRepository)

	c := context.Background()
	const val = int64(123)

	err = oauthStateDataRepository.StoreOAuthState(c, val)
	assert.Nil(t, err)

	err = oauthStateDataRepository.CheckOAuthState(c, val)
	assert.Nil(t, err)
}

func TestOAuthStateDataRepositorySetAndRemoveAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test which requires more setup.")
	}

	oauthStateDataRepository, err := NewOAuthStateDataRepository()
	assert.Nil(t, err)
	assert.NotNil(t, oauthStateDataRepository)

	c := context.Background()
	const val = int64(345)
	err = oauthStateDataRepository.StoreOAuthState(c, val)
	assert.Nil(t, err)

	err = oauthStateDataRepository.RemoveOAuthState(c, val)
	assert.Nil(t, err)

	err = oauthStateDataRepository.CheckOAuthState(c, val)
	assert.NotNil(t, err)
}
