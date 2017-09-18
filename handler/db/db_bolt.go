package db

import (
	"fmt"

	"github.com/boltdb/bolt"
)

type BoltDB struct {
	db *bolt.DB
}

func NewDB(path string) (*BoltDB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	d := &BoltDB{
		db: db,
	}
	d.createBucket("URL")
	return d, nil
}

func (d *BoltDB) Close() error {
	d.db.Close()
	return nil
}

func (d *BoltDB) Set(index, shortURL string) error {
	// Start a writable transaction.
	tx, err := d.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Use the transaction...
	b := tx.Bucket([]byte("URL"))
	err = b.Put([]byte(index), []byte(shortURL))
	if err != nil {
		return err
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *BoltDB) Get(index string) (string, error) {
	// Start a writable transaction.
	tx, err := d.db.Begin(false)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// Use the transaction...
	b := tx.Bucket([]byte("URL"))
	byteURL := b.Get([]byte(index))

	shortURL := string(byteURL)

	return shortURL, nil
}

func (d *BoltDB) createBucket(bucketName string) error {
	// Start a writable transaction.
	tx, err := d.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Use the transaction...
	_, err = tx.CreateBucket([]byte(bucketName))
	if err != nil {
		return err
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *BoltDB) Len() (uint64, error) {
	// Start a writable transaction.
	tx, err := d.db.Begin(true)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Use the transaction...
	b := tx.Bucket([]byte("URL"))

	return b.NextSequence()
}
