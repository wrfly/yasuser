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
	d.createBucket("LONG")  // shortURL -> longURL
	d.createBucket("SHORT") // longURL's MD5 -> shortURL
	return d, nil
}

func (d *BoltDB) Close() error {
	d.db.Close()
	return nil
}

func (d *BoltDB) SetShort(index, shortURL string) error {
	// Start a writable transaction.
	tx, err := d.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Use the transaction...
	b := tx.Bucket([]byte("SHORT"))
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

func (d *BoltDB) GetShort(index string) (string, error) {
	// Start a writable transaction.
	tx, err := d.db.Begin(false)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// Use the transaction...
	b := tx.Bucket([]byte("SHORT"))
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
	// Start a writable transaction.
	tx, err := d.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Use the transaction...
	b := tx.Bucket([]byte("LONG"))
	err = b.Put([]byte(shortURL), []byte(longURL))
	if err != nil {
		return err
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *BoltDB) GetLong(shortURL string) (string, error) {
	// Start a writable transaction.
	tx, err := d.db.Begin(false)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// Use the transaction...
	b := tx.Bucket([]byte("LONG"))
	byteURL := b.Get([]byte(shortURL))

	longURL := string(byteURL)

	return longURL, nil
}
