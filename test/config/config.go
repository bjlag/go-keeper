package config

type Config struct {
	Migration struct {
		SourcePath string `yaml:"sourcePath" env:"MIGRATION_SOURCE_PATH" env-description:"Path to migration source" json:"source_path"`
		Table      string `yaml:"table" env:"MIGRATION_TABLE" env-description:"Migration table" json:"table"`
	} `yaml:"migration" json:"migration"`

	Container struct {
		PG struct {
			Tag        string `yaml:"tag" env:"CONTAINER_PG_TAG" env-description:"Container image tag" json:"tag"`
			DBName     string `yaml:"dbName" env:"CONTAINER_PG_DB_NAME" env-description:"DB name" json:"db_name"`
			DBUser     string `yaml:"dbUser" env:"CONTAINER_PG_DB_USER" env-description:"DB user" json:"db_user"`
			DBPassword string `yaml:"dbPassword" env:"CONTAINER_PG_DB_PASSWORD" env-description:"DB password" json:"db_password"`
		} `yaml:"pg" json:"pg"`
	} `yaml:"container" json:"container"`
}
