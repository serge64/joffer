package storagepg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"guthub.com/serge64/joffer/internal/model"
	"guthub.com/serge64/joffer/internal/storage"
	"guthub.com/serge64/joffer/internal/storage/storagepg"
)

func TestUserRepository_Create(t *testing.T) {
	store, teardown := storagepg.TestDB(t, databaseURL)
	defer teardown("users")

	u := model.TestUser(t)
	id, err := store.User().Create(u)

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.EqualValues(t, u.ID, id)
}

func TestUserRepository_Find(t *testing.T) {
	store, teardown := storagepg.TestDB(t, databaseURL)
	defer teardown("users")

	u1 := model.TestUser(t)
	id, err := store.User().Create(u1)

	assert.NoError(t, err)
	assert.NotNil(t, u1)
	assert.EqualValues(t, u1.ID, id)

	u2, err := store.User().Find(u1.ID)

	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	store, teardown := storagepg.TestDB(t, databaseURL)
	defer teardown("users")

	u1 := model.TestUser(t)

	_, err := store.User().FindByEmail(u1.Email)
	assert.EqualError(t, err, storage.ErrRecordNotFound.Error())

	_, err = store.User().Create(u1)
	assert.NoError(t, err)

	u2, err := store.User().FindByEmail(u1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
