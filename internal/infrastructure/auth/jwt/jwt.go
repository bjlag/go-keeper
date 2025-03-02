package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	prefixOp = "auth.jwt."

	subjectAccessToken  = "access_token"
	subjectRefreshToken = "refresh_token"
)

var (
	ErrInvalidToken            = errors.New("invalid token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrUnexpectedTypeToken     = errors.New("unexpected type token")
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

func (g *Generator) GetUserGUIDFromAccessToken(tokenString string) (uuid.UUID, error) {
	const op = prefixOp + "GetUserGUIDFromAccessToken"

	token, clams, err := g.ParseToken(tokenString, subjectAccessToken)
	if err != nil {
		var e *jwt.ValidationError
		if errors.As(err, &e) {
			return uuid.Nil, ErrInvalidToken
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	if !token.Valid || clams.Issuer == "" {
		return uuid.Nil, ErrInvalidToken
	}

	guid, err := uuid.Parse(clams.Issuer)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return guid, nil
}

func (g *Generator) GetUserGUIDFromRefreshToken(tokenString string) (uuid.UUID, error) {
	const op = prefixOp + "GetUserGUIDFromRefreshToken"

	token, clams, err := g.ParseToken(tokenString, subjectRefreshToken)
	if err != nil {
		var e *jwt.ValidationError
		if errors.As(err, &e) {
			return uuid.Nil, ErrInvalidToken
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	if !token.Valid || clams.Issuer == "" {
		return uuid.Nil, ErrInvalidToken
	}

	guid, err := uuid.Parse(clams.Issuer)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return guid, nil
}

func (g *Generator) GenerateTokens(uuid string) (accessToken string, refreshToken string, err error) {
	const op = prefixOp + "GenerateTokens"

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
	const op = prefixOp + "GenerateAccessToken"

	now := time.Now()
	claim := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    uuid,
			ExpiresAt: jwt.NewNumericDate(now.Add(g.accessTokenExp)),
			Subject:   subjectAccessToken,
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
	const op = prefixOp + "GenerateRefreshToken"

	now := time.Now()
	claim := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cl.Issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(g.refreshTokenExp)),
			Subject:   subjectRefreshToken,
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

func (g *Generator) ParseToken(tokenString, subjectClaim string) (*jwt.Token, *Claims, error) {
	const op = prefixOp + "ParseToken"

	c := &Claims{}
	fn := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigningMethod, t.Header["alg"])
		}
		claim, ok := t.Claims.(*Claims)
		if ok && claim.Subject != subjectClaim {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedTypeToken, claim.Subject)
		}
		return []byte(g.secretKey), nil
	}

	t, err := jwt.ParseWithClaims(tokenString, c, fn)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return t, c, nil
}
