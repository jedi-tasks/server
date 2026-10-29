package config

import (
	"pkg/env"
	"pkg/log"

	"pkg/database/pgsql"
)

// Config - общая структура конфига
type Config struct {

	// Адрес для http-сервера
	Listen struct {
		HTTP        string `env:"LISTEN_HTTP"`
		ServiceHTTP string `env:"LISTEN_SERVICE_HTTP"`
		GRPC        string `env:"LISTEN_GRPC"`
	}

	// Данные базы данных
	Pgsql pgsql.PgsqlConfigEnv

	// Информация для JWT-токенов
	Auth struct {
		AccessTokenTTL  string `env:"AUTH_ACCESS_TOKEN_TTL"`
		RefreshTokenTTL string `env:"AUTH_REFRESH_TOKEN_TTL"`
		SigningKey      string `env:"AUTH_TOKEN_SIGNING_KEY"`
		// Информация для шифрования паролей
		GeneralSalt string `env:"SHA_SALT"`

		// Открыта ли регистрация новых пользователей
		SignUpEnabled bool `env:"AUTH_SIGN_UP_ENABLED"`
	}

	ServiceName string `env:"SERVICE_NAME"`
	Environment string `env:"ENVIRONMENT"`

	Logger log.LoggerSettingsEnv
}

func Load() Config {
	return env.Load[Config]()
}
