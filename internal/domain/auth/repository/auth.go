package repository

import (
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/xyedo/blindate/internal/infrastructure"
)

func New() *Auth {
	client, err := clerk.NewClient(infrastructure.Config.ClrekToken)
	if err != nil {
		panic(err)
	}

	return &Auth{
		conn: client,
	}
}

type Auth struct {
	conn clerk.Client
}
