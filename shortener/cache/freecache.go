package cache

import (
	"runtime/debug"

	"github.com/coocood/freecache"
)

type Cacher struct {
	cache *freecache.Cache
}

func NewCacher(cacheSize int) Cacher {
	cache := freecache.NewCache(cacheSize)
	debug.SetGCPercent(20)
	return Cacher{cache: cache}
}

func (c Cacher) Set(key, val string) error {
	return c.cache.Set([]byte(key), []byte(val), -1)
}

func (c Cacher) SetWithExpire(key, val string, expireSeconds int) error {
	return c.cache.Set([]byte(key), []byte(val), expireSeconds)
}

func (c Cacher) Get(key string) (string, error) {
	bVal, err := c.cache.Get([]byte(key))
	if err != nil {
		return "", err
	}
	return string(bVal), nil
}

func (c Cacher) Del(key string) bool {
	return c.cache.Del([]byte(key))
}
