package db

import (
	"fmt"
	"sync/atomic"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
	"github.com/wrfly/yasuser/types"
)

const (
	shortBucket = "s"
	longBucket  = "l"
)

type boltDB struct {
	db     *bolt.DB
	length *int64
}

func newBoltDB(path string) (*boltDB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	b := &boltDB{
		db:     db,
		length: &skipKeyNums,
	}
	b.createBucket(longBucket)  // short -> URL
	b.createBucket(shortBucket) // Hash  -> URL

	b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(shortBucket))
		*b.length += int64(bkt.Stats().KeyN)
		return nil
	})

	return b, nil
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
		if URL.Custom == "" {
			return nil
		}
		// store custom
		err = tx.Bucket([]byte(longBucket)).
			Put([]byte(URL.Custom), URL.Bytes())
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

// Len returns the currerent length of keys, and +1
// use atomic for concurrency conflict handling
func (b *boltDB) Len() int64 {
	return atomic.AddInt64(b.length, 1) - 1
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
