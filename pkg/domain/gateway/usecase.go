package gateway

import (
	gatewayEntities "github.com/xyedo/blindate/pkg/domain/gateway/entities"
)

type Session interface {
	SetUserSocket(string, gatewayEntities.Conn)
	GetUserSocket(string) (gatewayEntities.Conn, bool)
	DeleteUserSocket(string)
}