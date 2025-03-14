package client

// Config хранит конфигурацию клиента.
type Config struct {
	// Env окружение.
	Env string `yaml:"env" env:"ENV" env-default:"dev" env-description:"Environment" json:"env"`

	// Server содержит настройки подключения к серверу.
	Server struct {
		// Host хост сервера.
		Host string `yaml:"host" env:"SERVER_HOST" env-description:"Server host" json:"host"`
		// Port порт сервера.
		Port int `yaml:"port" env:"SERVER_PORT" env-description:"Server port" json:"port"`
	} `yaml:"server" json:"server"`

	// MasterKey настройки мастер ключа.
	MasterKey struct {
		// SaltLength длина соли.
		SaltLength int `yaml:"saltLength" env:"SALT_LENGTH" env-description:"Salt length" json:"salt_length"`
		// IterCount количество итерация при генерации мастер ключа.
		IterCount int `yaml:"iterCount" env:"MASTER_KEY_ITER_COUNT" env-description:"Master key iteration" json:"iter_count"`
		// Length длина мастер ключа.
		Length int `yaml:"length" env:"MASTER_KEY_LENGTH" env-description:"Master key length" json:"length"`
	} `yaml:"masterKey" json:"master_key"`

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
