package db

import (
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
	"github.com/wrfly/yasuser/types"
)

const (
	shortBucket = "s"
	longBucket  = "l"
	statsBucket = "st"
)

var visitedKey = []byte("visited")

type boltDB struct {
	db      *bolt.DB
	length  *int64
	visited *int64
}

func newBoltDB(path string) (*boltDB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	visited := new(int64)
	b := &boltDB{
		db:      db,
		length:  &skipKeyNum,
		visited: visited,
	}
	b.createBucket(longBucket)  // short -> URL
	b.createBucket(shortBucket) // Hash  -> URL
	b.createBucket(statsBucket) // status

	err = b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(shortBucket))
		atomic.AddInt64(b.length, int64(bkt.Stats().KeyN))

		statsBkt := tx.Bucket([]byte(statsBucket))
		v := string(statsBkt.Get(visitedKey))
		if v != "" {
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}
			atomic.AddInt64(b.visited, i)
		}
		return nil
	})

	return b, err
}

func (b *boltDB) createBucket(bucketName string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func (b *boltDB) Close() error {
	b.db.Close()
	return nil
}

func (b *boltDB) Store(URL *types.URL) error {
	logrus.Debugf("store %v", URL)

	return b.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte(shortBucket)).
			Put(URL.HashSum(), URL.Bytes())
		if err != nil {
			return err
		}
		err = tx.Bucket([]byte(longBucket)).
			Put(URL.ShortURL(), URL.Bytes())
		if err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
}

func (b *boltDB) GetShort(hashSum string) (*types.URL, error) {
	return b.get(shortBucket, hashSum)
}

// Len returns how many keys in store
func (b *boltDB) Keys() (int64, error) {
	return atomic.LoadInt64(b.length), nil
}

// IncKey returns the current length of keys, and +1
// use atomic for concurrency conflict handling
func (b *boltDB) IncKey() (int64, error) {
	return atomic.AddInt64(b.length, 1) - 1, nil
}

func (b *boltDB) Visited() (int64, error) {
	return atomic.LoadInt64(b.visited), nil
}

func (b *boltDB) IncVisited() (int64, error) {
	err := b.db.Update(func(tx *bolt.Tx) error {
		statsBkt := tx.Bucket([]byte(statsBucket))
		return statsBkt.Put(visitedKey, []byte(fmt.Sprint(
			atomic.AddInt64(b.visited, 1),
		)))
	})
	return atomic.LoadInt64(b.visited), err
}

func (b *boltDB) GetLong(shortURL string) (*types.URL, error) {
	return b.get(longBucket, shortURL)
}

func (b *boltDB) get(bkName, key string) (*types.URL, error) {
	u := new(types.URL)
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bkName))
		v := b.Get([]byte(key))
		if len(v) == 0 {
			return types.ErrNotFound
		}
		u.Decode(v)
		return nil
	})
	logrus.Debugf("bolt get [%s]: %s=%v", bkName, key, u)

	return u, err
}
