package server

import "time"

type Config struct {
	Env string `yaml:"env" env:"ENV" env-default:"dev" env-description:"Environment" json:"env"`

	Address struct {
		Host string `yaml:"host" env:"ADDRESS_HOST" env-description:"Server address host" json:"host"`
		Port int    `yaml:"port" env:"ADDRESS_PORT" env-description:"Server address port" json:"port"`
	} `yaml:"address" json:"address"`

	Auth struct {
		AccessTokenExp  time.Duration `yaml:"accessTokenExp" env:"ACCESS_TOKEN_EXP" env-description:"Access token expiration" json:"access_token_exp"`
		RefreshTokenExp time.Duration `yaml:"refreshTokenExp" env:"REFRESH_TOKEN_EXP" env-description:"Refresh token expiration" json:"refresh_token_exp"`
		SecretKey       string        `yaml:"secretKey" env:"SECRET_KEY" env-description:"Secret key" json:"secret_key"`
	} `yaml:"auth" json:"auth"`

	Database struct {
		Host     string `yaml:"host" env:"DB_HOST" env-description:"Database host" json:"host"`
		Port     string `yaml:"port" env:"DB_PORT" env-description:"Database port" json:"port"`
		Name     string `yaml:"name" env:"DB_NAME" env-description:"Database name" json:"name"`
		User     string `yaml:"user" env:"DB_USER" env-description:"Database user" json:"user"`
		Password string `yaml:"password" env:"DB_PASSWORD" env-description:"Database password" json:"password"`
	} `yaml:"database" json:"database"`
}
