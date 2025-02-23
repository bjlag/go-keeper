package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidToken            = errors.New("invalid token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)

type Claims struct {
	jwt.RegisteredClaims
}

type Generator struct {
	secretKey       string
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration
}

func NewGenerator(secretKey string, accessTokenExp, refreshTokenExp time.Duration) *Generator {
	return &Generator{
		secretKey:       secretKey,
		accessTokenExp:  accessTokenExp,
		refreshTokenExp: refreshTokenExp,
	}
}

func (g *Generator) GetUserGUID(tokenString string) (string, error) {
	c := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, c, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigningMethod, t.Header["alg"])
		}
		return []byte(g.secretKey), nil
	})
	if err != nil {
		var e *jwt.ValidationError
		if errors.As(err, &e) {
			return "", ErrInvalidToken
		}
		return "", err
	}

	if !token.Valid || c.Issuer == "" {
		return "", ErrInvalidToken
	}

	return c.Issuer, nil
}

func (g *Generator) GenerateTokens(uuid string) (accessToken string, refreshToken string, err error) {
	const op = "GenerateTokens"

	var claim *Claims

	accessToken, claim, err = g.GenerateAccessToken(uuid)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		return
	}

	refreshToken, err = g.GenerateRefreshToken(claim)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
	}

	return
}

func (g *Generator) GenerateAccessToken(uuid string) (string, *Claims, error) {
	const op = "GenerateAccessToken"

	now := time.Now()
	claim := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    uuid,
			ExpiresAt: jwt.NewNumericDate(now.Add(g.accessTokenExp)),
			Subject:   "access_token",
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(g.secretKey))
	if err != nil {
		return "", nil, fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, claim, nil
}

func (g *Generator) GenerateRefreshToken(cl *Claims) (string, error) {
	const op = "GenerateRefreshToken"

	now := time.Now()
	claim := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cl.Issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(g.refreshTokenExp)),
			Subject:   "refresh_token",
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(g.secretKey))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, nil
}
