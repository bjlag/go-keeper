package config

import "time"

// Config хранит конфигурацию тестового окружения.
type Config struct {
	// Env окружение.
	Env string `yaml:"env" env:"ENV" env-default:"dev" env-description:"Environment" json:"env"`

	// Migration настройки магратора.
	Migration struct {
		// SourcePath путь до файлов миграций.
		SourcePath string `yaml:"sourcePath" env:"MIGRATION_SOURCE_PATH" env-description:"Path to migration source" json:"source_path"`
		// Table название таблицы с примененными миграциями.
		Table string `yaml:"table" env:"MIGRATION_TABLE" env-description:"Migration table" json:"table"`
	} `yaml:"migration" json:"migration"`

	// Auth настройки авторизации.
	Auth struct {
		// AccessTokenExp время жизни access токена.
		AccessTokenExp time.Duration `yaml:"accessTokenExp" env:"ACCESS_TOKEN_EXP" env-description:"Access token expiration" json:"access_token_exp"`
		// RefreshTokenExp время жизни refresh токена.
		RefreshTokenExp time.Duration `yaml:"refreshTokenExp" env:"REFRESH_TOKEN_EXP" env-description:"Refresh token expiration" json:"refresh_token_exp"`
		// SecretKey секретный ключ.
		SecretKey string `yaml:"secretKey" env:"SECRET_KEY" env-description:"Secret key" json:"secret_key"`
	} `yaml:"auth" json:"auth"`

	// Container настройки Docker контейнеров.
	Container struct {
		// PG настройки для контейнера под PostgreSQL.
		PG struct {
			// Tag контейнера
			Tag string `yaml:"tag" env:"CONTAINER_PG_TAG" env-description:"Container image tag" json:"tag"`
			// DBName значение для переменной окружения POSTGRES_DB
			DBName string `yaml:"dbName" env:"CONTAINER_PG_DB_NAME" env-description:"DB name" json:"db_name"`
			// DBUser значение для переменной окружения POSTGRES_USER
			DBUser string `yaml:"dbUser" env:"CONTAINER_PG_DB_USER" env-description:"DB user" json:"db_user"`
			// DBPassword значение для переменной окружения POSTGRES_PASSWORD
			DBPassword string `yaml:"dbPassword" env:"CONTAINER_PG_DB_PASSWORD" env-description:"DB password" json:"db_password"`
		} `yaml:"pg" json:"pg"`
	} `yaml:"container" json:"container"`
}
