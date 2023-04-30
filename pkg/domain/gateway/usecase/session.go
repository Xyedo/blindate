package usecase

import (
	"github.com/xyedo/blindate/internal/rwmap"
	"github.com/xyedo/blindate/pkg/domain/gateway"
	gatewayEntities "github.com/xyedo/blindate/pkg/domain/gateway/entities"
)

func NewSession() gateway.Session {
	return &gatewayUC{
		clients:    rwmap.New[string, gatewayEntities.Conn](),
	}
}

type gatewayUC struct {
	clients    *rwmap.RwMap[string, gatewayEntities.Conn]
}


func (g *gatewayUC) SetUserSocket(id string, conn gatewayEntities.Conn) {
	g.clients.Set(id, conn)
}

func (g *gatewayUC) GetUserSocket(id string) (gatewayEntities.Conn, bool) {
	return g.clients.Get(id)
}

func (g *gatewayUC) DeleteUserSocket(id string) {
	g.clients.Delete(id)

}
