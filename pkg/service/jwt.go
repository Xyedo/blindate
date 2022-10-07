package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/xyedo/blindate/pkg/domain"
)

var (
	ErrTokenExpired = errors.New("jwt: token expired")
)

type customClaims struct {
	CredentialId string `json:"credId,omitempty"`
	jwt.RegisteredClaims
}

type Jwt interface {
	GenerateAccessToken(id string) (string, error)
	GenerateRefreshToken(id string) (string, error)
	ValidateRefreshToken(token string) (string, error)
	ValidateAccessToken(token string) (string, error)
}

func NewJwt(accessSecret, refreshSecret, accessExpires, refreshExpires string) *jwtSvc {
	return &jwtSvc{
		accessSecret:   accessSecret,
		refreshSecret:  refreshSecret,
		accessExpires:  accessExpires,
		refreshExpires: refreshExpires,
	}
}

type jwtSvc struct {
	accessSecret   string
	refreshSecret  string
	accessExpires  string
	refreshExpires string
}

func (j *jwtSvc) GenerateAccessToken(id string) (string, error) {
	return generateToken(id, j.accessSecret, j.accessExpires)
}

func (j *jwtSvc) GenerateRefreshToken(id string) (string, error) {
	return generateToken(id, j.refreshSecret, j.refreshExpires)
}

func (j *jwtSvc) ValidateRefreshToken(token string) (string, error) {
	return validateToken(token, j.refreshSecret)

}

func (j *jwtSvc) ValidateAccessToken(token string) (string, error) {
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

func generateCustomClaims(id string, duration time.Duration) customClaims {
	return customClaims{
		CredentialId: id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
}

func validateToken(token, secret string) (string, error) {
	decodedToken, err := jwt.ParseWithClaims(token, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrNotMatchCredential
		}
		return []byte(secret), nil
	})
	if err != nil {
		var jwtErr *jwt.ValidationError
		if errors.As(err, &jwtErr) {
			if jwtErr.Errors == jwt.ValidationErrorExpired {
				return "", ErrTokenExpired
			}
		}
		return "", domain.ErrNotMatchCredential
	}
	claims, ok := decodedToken.Claims.(*customClaims)
	if !ok || !decodedToken.Valid {
		return "", domain.ErrNotMatchCredential
	}

	return claims.CredentialId, nil
}
