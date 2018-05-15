package db

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/wrfly/short-url/config"
)

const (
	BOLT  = "bolt"
	REDIS = "redis"
)

// Database is a KV storage, there two relationships
// md5sum -> short & short -> long
// md5sum is the URL's md5sum
type Database interface {
	Close() error
	Len() (int64, error)
	SetShort(md5sum, shortURL string) error
	GetShort(md5sum string) (short string, err error)
	SetLong(shortURL, longURL string) error
	GetLong(md5sum string) (long string, err error)
}

// New DB storage
func New(conf config.StoreConfig) (Database, error) {
	logrus.Infof("Creating DB with [ %v ]", conf.DBType)
	switch conf.DBType {
	case BOLT:
		return newBoltDB(conf.DBPath)
	case REDIS:
		// TODO: newRedisDB()
	}
	return nil, fmt.Errorf("Unknown DB Type: %s", conf.DBType)
}
