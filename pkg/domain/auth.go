package domain

type AuthRepository interface {
	AddRefreshToken(token string) (int64, error)
	VerifyRefreshToken(token string) error
	DeleteRefreshToken(token string) (int64, error)
}

type AuthService interface {
	AddRefreshToken(token string) error
	VerifyRefreshToken(token string) error
	DeleteRefreshToken(token string) error
}
