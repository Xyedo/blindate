package authentication

type Repository interface {
	AddRefreshToken(token string) error
	VerifyRefreshToken(token string) error
	DeleteRefreshToken(token string) error
}
