package config

import loadconfig "github.com/suhrobdomoiZ/yadro-test-2026/pkg/config"

const (
	HTTPServerPort loadconfig.Key = "HTTP_PORT"
	EnvType        loadconfig.Key = "ENV_TYPE"

	DatabaseName     loadconfig.Key = "DATABASE_NAME"
	DatabaseHost     loadconfig.Key = "DATABASE_HOST"
	DatabasePort     loadconfig.Key = "DATABASE_PORT"
	DatabaseUser     loadconfig.Key = "DATABASE_USER"
	DatabasePassword loadconfig.Key = "DATABASE_PASSWORD"
)
