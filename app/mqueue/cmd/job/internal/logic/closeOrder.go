package logic

import (
	"TravelHome/app/mqueue/cmd/job/internal/svc"
	"TravelHome/app/mqueue/cmd/job/jobtype"
	"TravelHome/app/order/cmd/rpc/order"
	"TravelHome/app/order/model"
	"TravelHome/pkg/xerr"
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
)

var ErrCloseOrderFal = xerr.NewErrMsg("close order fail")

type CloseOrderHander struct {
	svcCtx *svc.ServiceContext
}

func NewCloseOrderHander(svcCtx *svc.ServiceContext) *CloseOrderHander {
	return &CloseOrderHander{

		svcCtx: svcCtx,
	}
}

func (l *CloseOrderHander) ProcessTask(ctx context.Context, t *asynq.Task) error {

	// close logic

	//jiema
	var p jobtype.DeferCloseHomestayOrderPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return errors.Wrapf(ErrCloseOrderFal, "closeHomestayOrderStateMqHandler payload err:%v, payLoad:%+v", err, t.Payload())
	}

	//order detail
	resp, err := l.svcCtx.OrderRpc.HomestayOrderDetail(ctx, &order.HomestayOrderDetailReq{
		Sn: p.Sn,
	})
	if err != nil || resp.HomestayOrder == nil {
		return errors.Wrapf(ErrCloseOrderFal, "closeHomestayOrderStateMqHandler  get order fail or order no exists err:%v, sn:%s ,HomestayOrder : %+v", err, p.Sn, resp.HomestayOrder)
	}
	//order update

	if resp.HomestayOrder.TradeState == model.HomestayOrderTradeStateWaitPay {

		_, err := l.svcCtx.OrderRpc.UpdateHomestayOrderTradeState(ctx, &order.UpdateHomestayOrderTradeStateReq{
			Sn:         p.Sn,
			TradeState: model.HomestayOrderTradeStateCancel,
		})
		if err != nil {
			return errors.Wrapf(ErrCloseOrderFal, "CloseHomestayOrderHandler close order fail  err:%v, sn:%s ", err, p.Sn)
		}

	}
	return nil
}
