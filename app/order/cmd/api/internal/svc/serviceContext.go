// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"TravelHome/app/order/cmd/api/internal/config"
	"TravelHome/app/order/cmd/rpc/order"
	"TravelHome/app/travel/cmd/rpc/travel"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc  order.Order
	TravelRpc travel.Travel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		OrderRpc:  order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		TravelRpc: travel.NewTravel(zrpc.MustNewClient(c.TravelRpcConf)),
	}
}
