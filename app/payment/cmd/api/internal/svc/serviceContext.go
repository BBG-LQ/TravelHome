// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"TravelHome/app/order/cmd/rpc/order"
	"TravelHome/app/payment/cmd/api/internal/config"
	"TravelHome/app/payment/cmd/rpc/payment"
	"TravelHome/app/usercenter/cmd/rpc/usercenter"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	WxPayClient *core.Client

	OrderRpc      order.Order
	UsercenterRpc usercenter.Usercenter
	PaymentRpc    payment.Payment
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config: c,

		OrderRpc:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
		PaymentRpc:    payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
	}
}
