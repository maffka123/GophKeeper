package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

// InitLogger configures zap logger
func InitLogger(debug bool, projectID string) (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig.LevelKey = "severity"
	zapConfig.EncoderConfig.MessageKey = "message"

	if debug {
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, err := zapConfig.Build(zap.Fields(
		zap.String("projectID", projectID),
	))

	if err != nil {
		return nil, err
	}

	return logger, nil
}

type Config struct {
	Endpoint string `env:"SERVER_ADDRESS"`
	Debug    bool   `env:"SERVER_DEBUG"`
	DBpath   string `env:"DATABASE_URI"`
	Key      string `env:"KEY"`
}

// InitConfig initialises config, first from flags, then from env, so that env overwrites flags
func InitConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.Endpoint, "a", "127.0.0.1:8082", "server address as host:port")
	flag.StringVar(&cfg.DBpath, "d", "postgres://postgres:pass@localhost:5431/secrets?pool_max_conns=10", "path for connection with pg: postgres://postgres:pass@localhost:5431/secrets?pool_max_conns=10")
	flag.BoolVar(&cfg.Debug, "debug", true, "key for hash function")
	flag.StringVar(&cfg.Key, "k", "someweirdkey!", "key for hash function")

	flag.Parse()

	err := env.Parse(&cfg)

	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
