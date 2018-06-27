package db

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/wrfly/yasuser/config"
)

const (
	BOLT  = "bolt"
	REDIS = "redis"
)

var skipKeyNums int64 = 99

// Database is a KV storage, there two relationships
// hashSum -> short & short -> long
// hashSum is the URL's hashSum
type Database interface {
	Close() error
	Len() int64
	SetShort(hashSum, shortURL string) error
	GetShort(hashSum string) (short string, err error)
	SetLong(shortURL, longURL string) error
	GetLong(hashSum string) (long string, err error)
}

// New DB storage
func New(conf config.StoreConfig) (Database, error) {
	logrus.Infof("Creating DB with [ %v ]", conf.DBType)
	switch conf.DBType {
	case BOLT:
		return newBoltDB(conf.DBPath)
	case REDIS:
		return newRedisDB(conf.Redis)
	}
	return nil, fmt.Errorf("Unknown DB Type: %s", conf.DBType)
}
