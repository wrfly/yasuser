package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoltDB(t *testing.T) {
	key := "5d41402abc4b2a76b9719d911017c592"
	URL := "http://kfd.me"
	db, err := NewDB("/tmp/myyxxy.db")
	assert.NoError(t, err)
	defer db.Close()

	err = db.Set(key, URL)
	assert.NoError(t, err)

	u, err := db.Get(key)
	assert.NoError(t, err)

	if u != URL {
		t.Error("!=")
	}
	fmt.Println(u)

	u, err = db.Get("nonono")
	assert.NoError(t, err)
	fmt.Println(u)

}

func TestBoltDBLen(t *testing.T) {
	db, err := NewDB("/tmp/myyxxy.db")
	assert.NoError(t, err)
	defer db.Close()

	l, err := db.Len()
	assert.NoError(t, err)
	fmt.Println(l)
}
