package kq

import (
	"TravelHome/app/order/cmd/mq/internal/svc"
	"TravelHome/app/order/cmd/rpc/order"
	ordermodel "TravelHome/app/order/model"
	paymentmodel "TravelHome/app/payment/model"
	"TravelHome/pkg/kqueue"
	"TravelHome/pkg/xerr"
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

/*
*
Listening to the payment flow status change notification message queue
*/
type PaymentUpdateStatusMq struct {
	svcCtx *svc.ServiceContext
}

func NewPaymentUpdateStatusMq(svcCtx *svc.ServiceContext) *PaymentUpdateStatusMq {
	return &PaymentUpdateStatusMq{

		svcCtx: svcCtx,
	}
}

func (l *PaymentUpdateStatusMq) Consume(ctx context.Context, _, val string) error {

	var message kqueue.ThirdPaymentUpdatePayStatusNotifyMessage

	if err := json.Unmarshal([]byte(val), &message); err != nil {
		//logx
		return err

	}

	if err := l.execService(ctx, message); err != nil {
		//logx
		return err

	}

	return nil

}

func (l *PaymentUpdateStatusMq) execService(ctx context.Context, message kqueue.ThirdPaymentUpdatePayStatusNotifyMessage) error {

	orderTradeState := l.getOrderTradeStateByPaymentTradeState(message.PayStatus)

	if orderTradeState != -99 {
		_, err := l.svcCtx.OrderRpc.UpdateHomestayOrderTradeState(ctx, &order.UpdateHomestayOrderTradeStateReq{
			Sn:         message.OrderSn,
			TradeState: orderTradeState,
		})
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("update homestay order state fail"), "update homestay order state fail err : %v ,message:%+v", err, message)
		}
	}

	return nil
}

func (l *PaymentUpdateStatusMq) getOrderTradeStateByPaymentTradeState(thirdPaymentPayStatus int64) int64 {
	switch thirdPaymentPayStatus {
	case paymentmodel.ThirdPaymentPayTradeStateSuccess:
		return ordermodel.HomestayOrderTradeStateWaitUse
	case paymentmodel.ThirdPaymentPayTradeStateRefund:
		return ordermodel.HomestayOrderTradeStateRefund
	default:
		return -99
	}

}
