package token

type tokens struct {
	accessToken  string
	refreshToken string
}

type Store struct {
	tokens tokens
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) SaveTokens(accessToken, refreshToken string) {
	s.tokens.accessToken = accessToken
	s.tokens.refreshToken = refreshToken
}

func (s *Store) AccessToken() string {
	return s.tokens.accessToken
}

func (s *Store) RefreshToken() string {
	return s.tokens.refreshToken
}
