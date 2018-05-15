package db

import (
	"fmt"
	"sync/atomic"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
	"github.com/wrfly/short-url/types"
)

type BoltDB struct {
	db     *bolt.DB
	length *int64
}

func newBoltDB(path string) (*BoltDB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	initLength := int64(-1)
	boltDB := &BoltDB{
		db:     db,
		length: &initLength,
	}
	boltDB.createBucket("LONG")  // shortURL -> longURL
	boltDB.createBucket("SHORT") // longURL's MD5 -> shortURL

	return boltDB, nil
}

func (boltDB *BoltDB) Close() error {
	boltDB.db.Close()
	return nil
}

func (boltDB *BoltDB) SetShort(md5sum, shortURL string) error {
	err := boltDB.set("SHORT", md5sum, shortURL)
	if err == nil {
		atomic.AddInt64(boltDB.length, 1)
	}
	return err
}

func (boltDB *BoltDB) GetShort(md5sum string) (string, error) {
	return boltDB.get("SHORT", md5sum)
}

func (boltDB *BoltDB) createBucket(bucketName string) error {
	return boltDB.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func (boltDB *BoltDB) Len() (int64, error) {
	if atomic.LoadInt64(boltDB.length) != -1 {
		return *boltDB.length + 1, nil
	}

	err := boltDB.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("SHORT"))
		*boltDB.length = int64(b.Stats().KeyN)
		return nil
	})

	return atomic.LoadInt64(boltDB.length) + 1, err
}

func (boltDB *BoltDB) SetLong(shortURL, longURL string) error {
	return boltDB.set("LONG", shortURL, longURL)
}

func (boltDB *BoltDB) GetLong(shortURL string) (string, error) {
	return boltDB.get("LONG", shortURL)
}

func (boltDB *BoltDB) set(bkName, key, value string) error {
	logrus.Debugf("bolt set [%s]: '%s'='%s'", bkName, key, value)
	return boltDB.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bkName))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
}

func (boltDB *BoltDB) get(bkName, key string) (value string, err error) {
	err = boltDB.db.View(func(tx *bolt.Tx) error {
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
