package storagepg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"guthub.com/serge64/joffer/internal/model"
	"guthub.com/serge64/joffer/internal/storage/storagepg"
)

func TestLetterRepository_Create(t *testing.T) {
	store, teardown := storagepg.TestDB(t, databaseURL)
	defer teardown("letters")

	letter := model.TestLetter(t)
	id, err := store.Letter().Create(letter)

	assert.NoError(t, err)
	assert.NotNil(t, letter)
	assert.EqualValues(t, letter.ID, id)
}

func TestLetterRepository_Update(t *testing.T) {
	store, teardown := storagepg.TestDB(t, databaseURL)
	defer teardown("letters")

	letter := model.TestLetter(t)
	_, err := store.Letter().Create(letter)

	assert.NoError(t, err)

	letter.Body = "Другое тело"
	err = store.Letter().Update(letter)

	assert.NoError(t, err)

	letters, err := store.Letter().Find(letter.UserID)

	assert.NoError(t, err)
	assert.NotNil(t, letters)
	assert.Equal(t, letter, &letters[0])
}

func TestLetterRepository_Find(t *testing.T) {
	store, teardown := storagepg.TestDB(t, databaseURL)
	defer teardown("letters")

	letter := model.TestLetter(t)
	id, err := store.Letter().Create(letter)

	assert.NoError(t, err)
	assert.NotNil(t, letter)
	assert.EqualValues(t, letter.ID, id)

	letters, err := store.Letter().Find(letter.UserID)

	assert.NoError(t, err)
	assert.NotNil(t, letters)
	assert.Equal(t, letter, &letters[0])
}

func TestLetterRepository_Delete(t *testing.T) {
	store, teardown := storagepg.TestDB(t, databaseURL)
	defer teardown("letters")

	letter := model.TestLetter(t)
	id, err := store.Letter().Create(letter)

	assert.NoError(t, err)
	assert.NotEqual(t, id, 0)

	letters, err := store.User().Find(letter.UserID)

	assert.Error(t, err)
	assert.Nil(t, letters)
}
