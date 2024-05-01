package user

import (
	"context"

	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

type Repository interface {
	StoreUser(ctx context.Context, conn pg.Querier, id string) error
	FindUserDetailByIds(ctx context.Context, conn pg.Querier, ids []string) (entities.UserDetails, error)
	GetUserById(ctx context.Context, conn pg.Querier, id string, opts ...entities.GetUserOption) (entities.User, error)
	SoftDeleteUserById(ctx context.Context, conn pg.Querier, id string) error

	StoreUserDetail(ctx context.Context, conn pg.Querier, payload entities.UserDetail) (string, error)
	GetUserDetailById(ctx context.Context, conn pg.Querier, id string, opts ...entities.GetUserDetailOption) (entities.UserDetail, error)
	UpdateUserDetailById(ctx context.Context, conn pg.Querier, id string, payload entities.UpdateUserDetail) error

	FindNonMatchClosestUser(ctx context.Context, conn pg.Querier, payload entities.FindClosestUser) ([]string, error)

	StoreHobbiesByUserId(ctx context.Context, conn pg.Querier, userId string, hobbies []entities.Hobbie) error
	UpdateHobbies(ctx context.Context, conn pg.Querier, hobbies []entities.UpdateHobbie) error
	DeleteHobbiesByIds(ctx context.Context, conn pg.Querier, ids []string) error

	StoreMovieSeriesByUserId(ctx context.Context, conn pg.Querier, userId string, movieSeries []entities.MovieSerie) error
	UpdateMovieSeries(ctx context.Context, conn pg.Querier, movieSeries []entities.UpdateMovieSeries) error
	DeleteMovieSeriesByIds(ctx context.Context, conn pg.Querier, ids []string) error

	StoreSportsByUserId(ctx context.Context, conn pg.Querier, userId string, sports []entities.Sport) error
	DeleteSportByIds(ctx context.Context, conn pg.Querier, ids []string) error
	UpdateSports(ctx context.Context, conn pg.Querier, sports []entities.UpdateSport) error

	StoreTravelingsByUserId(ctx context.Context, conn pg.Querier, userId string, travels []entities.Travel) error
	UpdateTravelings(ctx context.Context, conn pg.Querier, travels []entities.UpdateTravel) error
	DeleteTravelingByIds(ctx context.Context, conn pg.Querier, ids []string) error

	UpdateProfilePictureSelectedToFalseByUserId(ctx context.Context, conn pg.Querier, id string) error
	InsertProfilePicture(ctx context.Context, conn pg.Querier, profilePicture entities.ProfilePicture) (string, error)
}
