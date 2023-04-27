package authentication

type Usecase interface {
	Login(email, password string) (accessToken string, refreshToken string, err error)
	RevalidateRefreshToken(refreshToken string) (string, error)
	Logout(refreshToken string) error
}
