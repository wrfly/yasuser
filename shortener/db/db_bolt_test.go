package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wrfly/yasuser/types"
)

var tempDBPath = "/tmp/myyxxy.db"

func removeTempDB() {
	os.Remove(tempDBPath)
}

func TestBoltDB(t *testing.T) {
	removeTempDB()

	u := &types.URL{
		Short:    "_1B",
		Original: "http://kfd.me",
		Hash:     "5d41402abc4b2a76b9719d911017c592",
	}

	db, err := newBoltDB(tempDBPath)
	assert.NoError(t, err)
	defer db.Close()

	assert.NoError(t, db.Store(u))

	uShort, err := db.GetShort(u.Hash)
	assert.NoError(t, err)
	assert.Equal(t, uShort.Bytes(), u.Bytes())

	uLong, err := db.GetLong(u.Short)
	assert.NoError(t, err)
	assert.Equal(t, uLong.Bytes(), u.Bytes())

	_, err = db.GetShort("nonono")
	assert.Error(t, types.ErrNotFound)

	_, err = db.GetLong("nonono")
	assert.Error(t, types.ErrNotFound)
}

func TestBoltDBLen(t *testing.T) {
	removeTempDB()

	db, err := newBoltDB(tempDBPath)
	assert.NoError(t, err)
	defer db.Close()

	skipped := skipKeyNum
	count := 99
	for index := 0; index < count; index++ {
		u := &types.URL{
			Short:    fmt.Sprintf("%d", index),
			Original: fmt.Sprintf("http://u.kfd.me/index-%d", index),
			Hash:     fmt.Sprintf("%d", index),
		}
		assert.NoError(t, db.Store(u))
		db.IncKey()
	}

	k, _ := db.IncKey()
	assert.Equal(t, int64(count)+skipped, k)
}
