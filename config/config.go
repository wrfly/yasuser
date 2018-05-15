package config

type SrvConfig struct {
	Prefix string // https://url.kfd.me
	Port   int
}

type StoreConfig struct {
	DBPath string
	DBType string // file | redis
	Redis  string // redis hosts
}

type ShortenerConfig struct {
	Store StoreConfig
}

type Config struct {
	Debug     bool
	Shortener ShortenerConfig
	Server    SrvConfig
}
