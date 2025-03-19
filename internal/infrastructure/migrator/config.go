package migrator

// Config хранит конфигурацию мигратора.
type Config struct {
	// Env окружение.
	Env string `yaml:"env" env:"ENV" env-default:"dev" env-description:"Environment" json:"env"`
	// SourcePath путь до файлов миграций.
	SourcePath string `yaml:"sourcePath" env:"MIGRATIONS_SOURCE_PATH" env-default:"./migrations" env-description:"Migrations source path" json:"source_path"`
	// MigrationsTable название таблицы с примененными миграциями.
	MigrationsTable string `yaml:"migrationsTable" env:"MIGRATIONS_TABLE" env-default:"migrations" env-description:"Migrations table name" json:"migrations_table"`

	// Database настройки подключения к базе данной, к которой применяются миграции.
	Database struct {
		// Type тип базы данных: pg - PostgreSQL, sqlite - SQLite.
		Type DBType `yaml:"type" env:"DB_TYPE" env-default:"pg" env-description:"Database type" json:"type"`
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
