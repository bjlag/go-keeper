package client

const (
	// OptSaltKey слаг под опцию, которая хранит соль для генерации мастер ключа.
	OptSaltKey = "salt"
)

// Option описывает опцию, которую храним в БД клиента.
type Option struct {
	// Slug какое-то название опции.
	Slug string
	// Value значение опции.
	Value string
}
