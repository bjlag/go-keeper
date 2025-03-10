package master_key

type SaltGenerator struct {
	length int
}

func NewSaltGenerator(length int) *SaltGenerator {
	return &SaltGenerator{
		length: length,
	}
}

func (g SaltGenerator) GenerateSalt() (*Salt, error) {
	return NewSalt(g.length)
}
