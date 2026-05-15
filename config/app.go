package config

import "github.com/suhrobdomoiZ/yadro-test-2026/pkg/logger"

type AppConfig struct {
	httpPort string
	envType  string

	DbConfig      *PostgresConfig
	TimeoutConfig *TimeoutConfig
}

func NewAppConfig() *AppConfig {
	httpPort := HTTPServerPort.MustGet()
	envType := EnvType.Get(logger.EnvLocal)

	return &AppConfig{httpPort, envType, NewPostgresConfig(), NewTimeoutConfig()}
}

func (c *AppConfig) HTTPPort() string {
	return c.httpPort
}

func (c *AppConfig) EnvType() string {
	return c.envType
}
