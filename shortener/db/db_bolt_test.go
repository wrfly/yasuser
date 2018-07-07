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

	key := "5d41402abc4b2a76b9719d911017c592"
	URL := "http://kfd.me"
	db, err := newBoltDB(tempDBPath)
	assert.NoError(t, err)
	defer db.Close()

	err = db.SetShort(key, URL)
	assert.NoError(t, err)

	u, err := db.GetShort(key)
	assert.NoError(t, err)

	assert.Equal(t, u, URL)

	u, err = db.GetShort("nonono")
	assert.Error(t, types.ErrNotFound)
}

func TestBoltDBLen(t *testing.T) {
	removeTempDB()

	db, err := newBoltDB(tempDBPath)
	assert.NoError(t, err)
	defer db.Close()

	skipped := skipKeyNums
	count := 99
	for index := 0; index < count; index++ {
		long := fmt.Sprintf("http://u.kfd.me/index-%d", index)
		hash := fmt.Sprintf("%d", index)
		assert.NoError(t, db.SetShort(hash, long))
		db.Len()
	}

	assert.Equal(t, int64(count)+skipped, db.Len())
}
