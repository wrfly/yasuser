package db

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/wrfly/yasuser/config"
)

const (
	BOLT  = "bolt"
	REDIS = "redis"
)

var skipKeyNums int64 = 99

// Database is a KV storage, there are two relationships
// hashSum -> short & short -> long
// hashSum is the URL's hashSum
type Database interface {
	Close() error
	Len() int64
	Store(hashSum, shortURL, longURL string) error
	StoreWithTTL(hashSum, shortURL, longURL string, ttl time.Duration) error
	GetShort(hashSum string) (short string, err error)
	GetLong(shortURL string) (long string, err error)
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
