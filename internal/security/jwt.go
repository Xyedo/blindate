package security

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidCred = errors.New("jwt:invalid credentials")
	ErrExpired     = errors.New("jwt:expired")
)

type customClaims struct {
	CredentialId string `json:"credId,omitempty"`
	jwt.RegisteredClaims
}

func NewJwt(accessSecret, refreshSecret, accessExpires, refreshExpires string) *Jwt {
	return &Jwt{
		accessSecret:   accessSecret,
		refreshSecret:  refreshSecret,
		accessExpires:  accessExpires,
		refreshExpires: refreshExpires,
	}
}

type Jwt struct {
	accessSecret   string
	refreshSecret  string
	accessExpires  string
	refreshExpires string
}

func (j *Jwt) GenerateAccessToken(id string) string {
	token, err := generateToken(id, j.accessSecret, j.accessExpires)
	if err != nil {
		log.Panic(err)
	}
	return token
}

func (j *Jwt) GenerateRefreshToken(id string) string {
	token, err := generateToken(id, j.refreshSecret, j.refreshExpires)
	if err != nil {
		log.Panic(err)
	}
	return token
}

func (j *Jwt) ValidateRefreshToken(token string) (string, error) {
	return validateToken(token, j.refreshSecret)

}

func (j *Jwt) ValidateAccessToken(token string) (string, error) {
	return validateToken(token, j.accessSecret)
}

func generateToken(id, secret string, expires string) (string, error) {
	duration, err := time.ParseDuration(expires)
	if err != nil {
		panic(err)
	}
	claims := generateCustomClaims(id, duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encodedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return encodedToken, nil
}

func validateToken(token, secret string) (string, error) {
	decodedToken, err := jwt.ParseWithClaims(token, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidCred
		}
		return []byte(secret), nil
	})
	if err != nil {
		var jwtErr *jwt.ValidationError
		if errors.As(err, &jwtErr) {
			if jwtErr.Errors == jwt.ValidationErrorExpired {
				return "", ErrExpired
			}
		}
		return "", err
	}
	claims, ok := decodedToken.Claims.(*customClaims)
	if !ok || !decodedToken.Valid {
		return "", ErrInvalidCred
	}

	return claims.CredentialId, nil
}

func generateCustomClaims(id string, duration time.Duration) customClaims {
	return customClaims{
		CredentialId: id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
}
