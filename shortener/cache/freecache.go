package cache

import (
	"runtime/debug"
	"time"

	"github.com/coocood/freecache"
	"github.com/sirupsen/logrus"
	"github.com/wrfly/yasuser/types"
)

type Cacher struct {
	cache *freecache.Cache
}

func NewCacher(cacheSize int) Cacher {
	cache := freecache.NewCache(cacheSize)
	debug.SetGCPercent(20)
	return Cacher{cache: cache}
}

func (c Cacher) Store(u *types.URL) {
	exp := -1
	if u.Expire != nil && !u.Expire.IsZero() {
		exp = int(u.Expire.Sub(time.Now()).Seconds())
		if exp < 0 {
			// less than 1s
			return
		}
		logrus.Debugf("cache store %v with ttl %d", u, exp)
	}
	c.cache.Set(u.ShortURL(), u.Bytes(), exp)
	c.cache.Set(u.HashSum(), u.Bytes(), exp)
	if u.Custom != "" {
		c.cache.Set([]byte(u.Custom), u.Bytes(), exp)
	}
}

func (c Cacher) Get(key string) (*types.URL, error) {
	bVal, err := c.cache.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	if len(bVal) == 0 {
		return nil, types.ErrNotFound
	}
	u := new(types.URL)
	u.Decode(bVal)
	return u, nil
}

func (c Cacher) Del(key string) bool {
	return c.cache.Del([]byte(key))
}
