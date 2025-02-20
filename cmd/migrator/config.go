package main

type Config struct {
	Env             string `yaml:"env" env:"ENV" env-default:"dev" env-description:"Environment" json:"env"`
	SourcePath      string `yaml:"sourcePath" env:"MIGRATIONS_SOURCE_PATH" env-default:"./migrations" env-description:"Migrations source path" json:"source_path"`
	MigrationsTable string `yaml:"migrationsTable" env:"MIGRATIONS_TABLE" env-default:"migrations" env-description:"Migrations table name" json:"migrations_table"`

	Database struct {
		Host     string `yaml:"host" env:"DB_HOST" env-description:"Database host" json:"host"`
		Port     string `yaml:"port" env:"DB_PORT" env-description:"Database port" json:"port"`
		Name     string `yaml:"name" env:"DB_NAME" env-description:"Database name" json:"name"`
		User     string `yaml:"user" env:"DB_USER" env-description:"Database user" json:"user"`
		Password string `yaml:"password" env:"DB_PASSWORD" env-description:"Database password" json:"password"`
	} `yaml:"database" json:"database"`
}
