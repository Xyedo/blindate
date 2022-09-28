package tokenizer

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	accessExpiresAtStr = os.Getenv("JWT_ACCESS_EXPIRES")
	accessSecret       = os.Getenv("JWT_ACCESS_SECRET_KEY")

	refreshExpiresAtStr = os.Getenv("JWT_REFRESH_EXPIRES")
	refreshSecret       = os.Getenv("JWT_REFRESH_SECRET_KEY")
)

var (
	ErrNotValidCredential = errors.New("jwt: credential is not valid")
	ErrTokenExpired       = errors.New("jwt: token expired")
)

type customClaims struct {
	CredentialId string `json:"credId,omitempty"`
	jwt.RegisteredClaims
}

type Jwt struct{}

func (j *Jwt) GenerateAccessToken(id string) (string, error) {
	return generateToken(id, accessSecret)
}

func (j *Jwt) GenerateRefreshToken(id string) (string, error) {
	return generateToken(id, refreshSecret)
}

func (j *Jwt) ValidateRefreshToken(token string) (string, error) {
	return validateToken(token, refreshSecret)

}

func (j *Jwt) ValidateAccessToken(token string) (string, error) {
	return validateToken(token, accessSecret)
}

func generateToken(id, secret string) (string, error) {
	duration, err := time.ParseDuration(refreshExpiresAtStr)
	if err != nil {
		panic(err)
	}
	claims := generateCustomClaims(id, duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encodedToken, err := token.SignedString(secret)
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
			return nil, ErrNotValidCredential
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
		return "", ErrNotValidCredential
	}
	claims, ok := decodedToken.Claims.(*customClaims)
	if !ok || !decodedToken.Valid {
		return "", ErrNotValidCredential
	}

	return claims.CredentialId, nil
}
