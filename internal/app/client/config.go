package client

type Config struct {
	Env string `yaml:"env" env:"ENV" env-default:"dev" env-description:"Environment" json:"env"`

	Server struct {
		Host string `yaml:"host" env:"SERVER_HOST" env-description:"Server host" json:"host"`
		Port int    `yaml:"port" env:"SERVER_PORT" env-description:"Server port" json:"port"`
	} `yaml:"server" json:"server"`

	Database struct {
		Host     string `yaml:"host" env:"DB_HOST" env-description:"Database host" json:"host"`
		Port     string `yaml:"port" env:"DB_PORT" env-description:"Database port" json:"port"`
		Name     string `yaml:"name" env:"DB_NAME" env-description:"Database name" json:"name"`
		User     string `yaml:"user" env:"DB_USER" env-description:"Database user" json:"user"`
		Password string `yaml:"password" env:"DB_PASSWORD" env-description:"Database password" json:"password"`
	} `yaml:"database" json:"database"`
}
