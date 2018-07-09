package db

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
	"github.com/wrfly/yasuser/types"
)

const (
	shortBucket  = "s"
	longBucket   = "l"
	expireBucket = "e"
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
	b.createBucket(longBucket)   // shortURL -> longURL
	b.createBucket(shortBucket)  // longURL's MD5 -> shortURL
	b.createBucket(expireBucket) // shortURL -> TeL

	b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(shortBucket))
		*b.length += int64(bkt.Stats().KeyN)
		return nil
	})

	return b, nil
}

func (b *boltDB) Close() error {
	b.db.Close()
	return nil
}

func (b *boltDB) Store(hashSum, shortURL, longURL string) error {
	return b.StoreWithTTL(hashSum, shortURL, longURL, -1)
}

func (b *boltDB) GetShort(hashSum string) (string, error) {
	return b.get(shortBucket, hashSum)
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

// Len returns the currerent length of keys, and +1
// use atomic for concurrency conflict handling
func (b *boltDB) Len() int64 {
	return atomic.AddInt64(b.length, 1) - 1
}

func (b *boltDB) GetLong(shortURL string) (string, error) {
	// check whether it's has ttl
	if strings.HasPrefix(shortURL, "_") {
		if expire, err := b.get(expireBucket, shortURL); err == nil {
			// err == nil -> got the expire data
			expireInt, _ := strconv.Atoi(expire)
			if time.Now().Unix()-int64(expireInt) > 0 {
				// expired
				return "", types.ErrNotFound
			}
		}
	}
	return b.get(longBucket, shortURL)
}

func (b *boltDB) get(bkName, key string) (value string, err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bkName))
		v := b.Get([]byte(key))
		value = string(v)
		if value == "" {
			return types.ErrNotFound
		}
		return nil
	})
	logrus.Debugf("bolt get [%s]: '%s'='%s'", bkName, key, value)

	return
}

func (b *boltDB) StoreWithTTL(hashSum, shortURL, longURL string, ttl time.Duration) error {
	storageInfo := fmt.Sprintf("store [%s]: '%s'='%s'", hashSum, shortURL, longURL)
	if ttl > 0 {
		storageInfo = fmt.Sprintf("%s; ttl: %s",
			storageInfo, ttl.String())
	}
	logrus.Debugf(storageInfo)

	return b.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte(shortBucket)).
			Put([]byte(hashSum), []byte(shortURL))
		if err != nil {
			return err
		}
		err = tx.Bucket([]byte(longBucket)).
			Put([]byte(shortURL), []byte(longURL))
		if err != nil {
			tx.Rollback()
			return err
		}
		if ttl > 0 {
			ttlInt := time.Now().Add(ttl).Unix()
			ttlStr := strconv.FormatInt(ttlInt, 10)

			err = tx.Bucket([]byte(expireBucket)).
				Put([]byte(shortURL), []byte(ttlStr))
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		return nil
	})
}
