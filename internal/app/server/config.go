package server

import "time"

// Config хранит конфигурацию сервера.
type Config struct {
	// Env окружение.
	Env string `yaml:"env" env:"ENV" env-default:"dev" env-description:"Environment" json:"env"`

	// Address содержит адрес сервера.
	Address struct {
		// Host хост сервера.
		Host string `yaml:"host" env:"ADDRESS_HOST" env-description:"Server host" json:"host"`
		// Port порт сервера.
		Port int `yaml:"port" env:"ADDRESS_PORT" env-description:"Server port" json:"port"`
	} `yaml:"address" json:"address"`

	// Auth настройки авторизации.
	Auth struct {
		// AccessTokenExp время жизни access токена.
		AccessTokenExp time.Duration `yaml:"accessTokenExp" env:"ACCESS_TOKEN_EXP" env-description:"Access token expiration" json:"access_token_exp"`
		// RefreshTokenExp время жизни refresh токена.
		RefreshTokenExp time.Duration `yaml:"refreshTokenExp" env:"REFRESH_TOKEN_EXP" env-description:"Refresh token expiration" json:"refresh_token_exp"`
		// SecretKey секретный ключ.
		SecretKey string `yaml:"secretKey" env:"SECRET_KEY" env-description:"Secret key" json:"secret_key"`
	} `yaml:"auth" json:"auth"`

	// Database настройки подключения к базе данных клиента.
	Database struct {
		// Host хост базы.
		Host string `yaml:"host" env:"DB_HOST" env-description:"Database host" json:"host"`
		// Port порт базы.
		Port string `yaml:"port" env:"DB_PORT" env-description:"Database port" json:"port"`
		// Name название базы данных.
		Name string `yaml:"name" env:"DB_NAME" env-description:"Database name" json:"name"`
		// User пользователь.
		User string `yaml:"user" env:"DB_USER" env-description:"Database user" json:"user"`
		// Password пароль.
		Password string `yaml:"password" env:"DB_PASSWORD" env-description:"Database password" json:"password"`
	} `yaml:"database" json:"database"`
}
