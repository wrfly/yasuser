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

	initLength := int64(99)
	b := &boltDB{
		db:     db,
		length: &initLength,
	}
	b.createBucket(longBucket)  // shortURL -> longURL
	b.createBucket(shortBucket) // longURL's MD5 -> shortURL

	return b, nil
}

func (b *boltDB) Close() error {
	b.db.Close()
	return nil
}

func (b *boltDB) SetShort(md5sum, shortURL string) error {
	err := b.set(shortBucket, md5sum, shortURL)
	if err == nil {
		atomic.AddInt64(b.length, 1)
	}
	return err
}

func (b *boltDB) GetShort(md5sum string) (string, error) {
	return b.get(shortBucket, md5sum)
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

func (b *boltDB) Len() int64 {
	return atomic.LoadInt64(b.length)
}

func (b *boltDB) SetLong(shortURL, longURL string) error {
	return b.set(longBucket, shortURL, longURL)
}

func (b *boltDB) GetLong(shortURL string) (string, error) {
	return b.get(longBucket, shortURL)
}

func (b *boltDB) set(bkName, key, value string) error {
	logrus.Debugf("bolt set [%s]: '%s'='%s'", bkName, key, value)
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bkName))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
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
