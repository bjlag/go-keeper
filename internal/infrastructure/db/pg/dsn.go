package pg

import "fmt"

// GetDSN вспомогательная функция для получения валидного DSN для подключения БД PostgreSQL.
func GetDSN(host, port, name, user, pass string) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		pass,
		host,
		port,
		name,
	)
}
