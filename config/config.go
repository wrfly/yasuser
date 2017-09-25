package config

type Config struct {
	Prefix string // https://url.kfd.me
	Port   int
	DBPath string
	DBType string // file | redis
	Redis  string // redis hosts
	Debug  bool
}
