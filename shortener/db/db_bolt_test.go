package db

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wrfly/yasuser/types"
)

var tempDBPath = "/tmp/myyxxy.db"

func removeTempDB() {
	os.Remove(tempDBPath)
}

func TestBoltDB(t *testing.T) {
	removeTempDB()

	hashSum := "5d41402abc4b2a76b9719d911017c592"
	URL := "http://kfd.me"
	shortURL := "1B"
	db, err := newBoltDB(tempDBPath)
	assert.NoError(t, err)
	defer db.Close()

	err = db.Store(hashSum, shortURL, URL)
	assert.NoError(t, err)

	uShort, err := db.GetShort(hashSum)
	assert.NoError(t, err)
	assert.Equal(t, uShort, shortURL)

	lShort, err := db.GetLong(shortURL)
	assert.NoError(t, err)
	assert.Equal(t, lShort, URL)

	_, err = db.GetShort("nonono")
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
		assert.NoError(t, db.Store(hash, hash, long))
		db.Len()
	}

	assert.Equal(t, int64(count)+skipped, db.Len())
}

func TestBoltDBStoreWithTTL(t *testing.T) {
	removeTempDB()

	hashSum := "5d41402abc4b2a76b9719d911017c592"
	URL := "http://kfd.me"
	shortURL := "_1B"
	db, err := newBoltDB(tempDBPath)
	assert.NoError(t, err)
	defer db.Close()

	err = db.StoreWithTTL(hashSum, shortURL, URL, time.Second)
	assert.NoError(t, err)

	time.Sleep(time.Second * 2)

	_, err = db.GetLong(shortURL)
	assert.EqualError(t, err, "not found")
}
