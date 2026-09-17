// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"TravelHome/app/mqueue/cmd/job/internal/config"
	"TravelHome/app/order/cmd/rpc/order"
	"TravelHome/app/usercenter/cmd/rpc/usercenter"

	"github.com/hibiken/asynq"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc      order.Order
	UsercenterRpc usercenter.Usercenter

	AsynqServer *asynq.Server
	MiniProgram *miniprogram.MiniProgram
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		OrderRpc:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
		AsynqServer:   newAsynqServer(c),
		MiniProgram:   newMiniprogramClient(c),
	}
}
