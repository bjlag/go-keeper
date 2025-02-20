package pg

import "fmt"

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
