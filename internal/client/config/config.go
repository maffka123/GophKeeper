package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"time"
)

type Config struct {
	Endpoint       string        `env:"APP_ADDRESS"`
	ServerEndpoint string        `env:"SERVER_ADDRESS"`
	AppName        string        `env:"APP_NAME" envDefault:"KeeperApp"`
	Debug          bool          `env:"APP_DEBUG"`
	DBpath         string        `env:"DATABASE_URI"`
	SyncInterval   time.Duration `env:"SYNC_PERIOD"`
}

// InitConfig initialises config, first from flags, then from env, so that env overwrites flags
func InitConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.Endpoint, "a", "127.0.0.1:8081", "server address as host:port")
	flag.StringVar(&cfg.ServerEndpoint, "s", "127.0.0.1:8082", "server address as host:port")
	flag.BoolVar(&cfg.Debug, "debug", true, "key for hash function")
	flag.StringVar(&cfg.DBpath, "d", "postgres://postgres:pass@localhost:5433/secrets?pool_max_conns=10", "path for connection with pg: postgres://postgres:pass@localhost:5431/secrets?pool_max_conns=10")
	flag.DurationVar(&cfg.SyncInterval, "si", 10*time.Second, "how often to sync with db")
	flag.Parse()

	err := env.Parse(&cfg)

	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
