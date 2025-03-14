// Package token in-memory хранилище токенов и мастер ключа на стороне клиента.
package token

type tokens struct {
	accessToken  string
	refreshToken string
	masterKey    []byte
}

type Store struct {
	tokens tokens
}

func NewStore() *Store {
	return &Store{}
}

// SaveTokens запомнить переданные токены.
func (s *Store) SaveTokens(accessToken, refreshToken string) {
	s.tokens.accessToken = accessToken
	s.tokens.refreshToken = refreshToken
}

// SaveMasterKey запомнить мастер ключ.
func (s *Store) SaveMasterKey(key []byte) {
	s.tokens.masterKey = key
}

// AccessToken получить access токен.
func (s *Store) AccessToken() string {
	return s.tokens.accessToken
}

// RefreshToken получить refresh токен.
func (s *Store) RefreshToken() string {
	return s.tokens.refreshToken
}

// MasterKey получить мастер ключ.
func (s *Store) MasterKey() []byte {
	return s.tokens.masterKey
}
