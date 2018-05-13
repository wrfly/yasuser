package db

import (
	"fmt"

	"github.com/boltdb/bolt"
)

type BoltDB struct {
	db *bolt.DB
}

func newBoltDB(path string) (*BoltDB, error) {
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
	d.createBucket("LONG")  // shortURL -> longURL
	d.createBucket("SHORT") // longURL's MD5 -> shortURL
	return d, nil
}

func (d *BoltDB) Close() error {
	d.db.Close()
	return nil
}

func (d *BoltDB) SetShort(index, shortURL string) error {
	return d.set("SHORT", index, shortURL)
}

func (d *BoltDB) GetShort(index string) (string, error) {
	return d.get("SHORT", index)
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

func (d *BoltDB) Len() (int, error) {
	// Start a writable transaction.
	tx, err := d.db.Begin(true)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Use the transaction...
	b := tx.Bucket([]byte("SHORT"))

	return b.Stats().KeyN, nil
}

func (d *BoltDB) SetLong(shortURL, longURL string) error {
	return d.set("LONG", shortURL, longURL)
}

func (d *BoltDB) GetLong(shortURL string) (string, error) {
	return d.get("LONG", shortURL)
}

func (d *BoltDB) set(bkName, key, value string) error {
	// Start a writable transaction.
	tx, err := d.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Use the transaction...
	b := tx.Bucket([]byte(bkName))
	err = b.Put([]byte(key), []byte(value))
	if err != nil {
		return err
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *BoltDB) get(bkName, key string) (string, error) {
	// Start a writable transaction.
	tx, err := d.db.Begin(false)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// Use the transaction...
	b := tx.Bucket([]byte(bkName))
	byteURL := b.Get([]byte(key))

	longURL := string(byteURL)

	return longURL, nil
}
