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

type Database interface {
	Close() error
	Len() (int, error)
	SetShort(index, shortURL string) error
	GetShort(index string) (string, error)
	SetLong(index, longURL string) error
	GetLong(index string) (string, error)
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
