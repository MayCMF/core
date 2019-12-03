package jwtauth

import (
	"context"
	"testing"

	"github.com/MayCMF/core/src/common/auth/jwtauth/store/buntdb"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	store, err := buntdb.NewStore(":memory:")
	assert.Nil(t, err)

	jwtAuth := New(store)

	defer jwtAuth.Release()

	ctx := context.Background()
	userUUID := "test"
	token, err := jwtAuth.GenerateToken(ctx, userUUID)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	id, err := jwtAuth.ParseUserUUID(ctx, token.GetAccessToken())
	assert.Nil(t, err)
	assert.Equal(t, userUUID, id)

	err = jwtAuth.DestroyToken(ctx, token.GetAccessToken())
	assert.Nil(t, err)

	id, err = jwtAuth.ParseUserUUID(ctx, token.GetAccessToken())
	assert.NotNil(t, err)
	assert.EqualError(t, err, "invalid token")
	assert.Empty(t, id)
}
