package auth

import (
	"sync"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/xyedo/blindate/internal/infrastructure"
)

type requestId string

const (
	RequestId requestId = "RequestId"
)

var (
	Auth clerk.Client
	once sync.Once
)

func initConn() {
	client, err := clerk.NewClient(infrastructure.Config.ClrekToken)
	if err != nil {
		panic(err)
	}

	Auth = client
}
func Get() clerk.Client {
	once.Do(func() {
		initConn()
	})

	return Auth
}
