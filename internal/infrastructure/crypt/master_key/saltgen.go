package master_key

// SaltGenerator генератор соли Salt.
// В length длина соли.
type SaltGenerator struct {
	length int
}

// NewSaltGenerator создает генератор соли с указанной длиной в length.
func NewSaltGenerator(length int) *SaltGenerator {
	return &SaltGenerator{
		length: length,
	}
}

// GenerateSalt генерирует соль.
func (g SaltGenerator) GenerateSalt() (*Salt, error) {
	return NewSalt(g.length)
}
