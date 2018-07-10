package db

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/wrfly/yasuser/config"
	"github.com/wrfly/yasuser/types"
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
	Store(URL types.URL) error
	GetShort(hashSum string) (URL types.URL, err error)
	GetLong(shortURL string) (URL types.URL, err error)
}

// New DB storage
func New(conf config.StoreConfig) (Database, error) {
	logrus.Infof("Creating DB with [ %v ]", conf.DBType)
	switch conf.DBType {
	case BOLT:
		return newBoltDB(conf.DBPath)
	case REDIS:
		// return newRedisDB(conf.Redis)
	}
	return nil, fmt.Errorf("Unknown DB Type: %s", conf.DBType)
}
