package db

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/wrfly/yasuser/types"
)

type redisDB struct {
	cli *redis.Client
}

const (
	KEY_NUMS = "KEY_NUMS"
	VISITED  = "VISITED"
)

func newRedisDB(redisAddr string) (*redisDB, error) {
	opts, err := redis.ParseURL(redisAddr)
	if err != nil {
		return nil, err
	}
	cli := redis.NewClient(opts)

	if err := cli.Ping().Err(); err != nil {
		return nil, err
	}
	if cli.Get(KEY_NUMS).Err() == redis.Nil {
		cli.Set(KEY_NUMS, skipKeyNum, -1)
	}

	return &redisDB{
		cli: cli,
	}, nil
}

func (r *redisDB) Close() error {
	return r.cli.Close()
}

// IncKey returns the current length of keys, and +1
// use atomic for concurrency conflict handling
func (r *redisDB) IncKey() (int64, error) {
	return r.incKey(KEY_NUMS)

}

func (r *redisDB) Keys() (int64, error) {
	return r.cli.Get(KEY_NUMS).Int64()
}

func (r *redisDB) Visited() (int64, error) {
	return r.cli.Get(VISITED).Int64()
}

func (r *redisDB) IncVisited() (int64, error) {
	return r.incKey(VISITED)
}

func (r *redisDB) incKey(key string) (int64, error) {
	length, err := r.cli.Incr(key).Result()
	if err != nil {
		return -1, err
	}
	return length, nil
}

func (r *redisDB) Store(URL *types.URL) error {
	var ttl time.Duration
	if URL.Expire != nil {
		ttl = URL.Expire.Sub(time.Now())
	}
	err := r.cli.Set(URL.Hash, URL.Bytes(), ttl).Err()
	if err != nil {
		return err
	}

	return r.cli.Set(URL.Short, URL.Bytes(), ttl).Err()
}

func (r *redisDB) GetShort(hashSum string) (URL *types.URL, err error) {
	return r.get(hashSum)
}

func (r *redisDB) GetLong(short string) (URL *types.URL, err error) {
	return r.get(short)
}

func (r *redisDB) get(key string) (URL *types.URL, err error) {
	stringCmd := r.cli.Get(key)
	if err := stringCmd.Err(); err != nil {
		if err == redis.Nil {
			return nil, types.ErrNotFound
		}
		return nil, err
	}
	bs, err := stringCmd.Bytes()
	if err != nil {
		return nil, err
	}
	u := new(types.URL)
	return u.Decode(bs), nil
}
