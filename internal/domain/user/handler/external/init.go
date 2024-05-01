package external

import "github.com/xyedo/blindate/internal/domain/user"

func New(userUsecase user.Usecase) *User {
	return &User{
		usecase: userUsecase,
	}
}

type User struct {
	usecase user.Usecase
}
