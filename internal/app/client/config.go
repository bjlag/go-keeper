package client

// Config хранит конфигурацию клиента.
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
		// Dir директория, где будут лежать файлы БД.
		Dir string `yaml:"dir" env:"DB_DIR" env-description:"Database directory" json:"dir"`
		// Prefix префикс в названии файла базы данных.
		Prefix string `yaml:"prefix" env:"DB_PREFIX" env-description:"Database prefix" json:"prefix"`
	} `yaml:"database" json:"database"`
}
