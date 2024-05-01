package handler

import "github.com/xyedo/blindate/internal/domain/match"

func New(matchUsecase match.Usecase) *Match {
	return &Match{
		usecase: matchUsecase,
	}
}

type Match struct {
	usecase match.Usecase
}
