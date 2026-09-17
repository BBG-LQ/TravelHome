package logic

import (
	"context"
	"encoding/json"

	"TravelHome/app/mqueue/cmd/job/jobtype"
	"TravelHome/app/order/cmd/rpc/internal/svc"
	"TravelHome/app/order/cmd/rpc/pb"
	"TravelHome/app/order/model"
	"TravelHome/pkg/xerr"

	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateHomestayOrderTradeStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateHomestayOrderTradeStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateHomestayOrderTradeStateLogic {
	return &UpdateHomestayOrderTradeStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新民宿订单状态
func (l *UpdateHomestayOrderTradeStateLogic) UpdateHomestayOrderTradeState(in *pb.UpdateHomestayOrderTradeStateReq) (*pb.UpdateHomestayOrderTradeStateResp, error) {
	// todo: add your logic here and delete this line

	//查询homeOrder
	homestayOrder, err := l.svcCtx.HomestayOrderModel.FindOneBySn(l.ctx, in.Sn)

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "UpdateHomestayOrderTradeState FindOneBySn db err : %v , in:%+v", err, in)
	}
	if homestayOrder == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("order no exists"), "order no exists , in:%+v", in)
	}

	if homestayOrder.TradeState == in.TradeState {

		return &pb.UpdateHomestayOrderTradeStateResp{}, nil

	}

	//检查

	err = l.verifyOrderTradeState(in.TradeState, homestayOrder.TradeState)
	if err != nil {
		return nil, errors.WithMessagef(err, " , in : %+v", in)
	}

	//更新
	homestayOrder.TradeState = in.TradeState
	if err := l.svcCtx.HomestayOrderModel.UpdateWithVersion(l.ctx, nil, homestayOrder); err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("Failed to update homestay order status"), "Failed to update homestay order status db UpdateWithVersion err:%v , in : %v", err, in)

	}

	//	notify user

	if in.TradeState == model.HomestayOrderTradeStateWaitUse {
		payload, err := json.Marshal(jobtype.PaySuccessNotifyUserPayload{
			Order: homestayOrder,
		})
		if err != nil {
			// logx.WithContext(l.ctx).Errorf("pay success notify user task json Marshal fail, err :%+v , sn : %s", err, homestayOrder.Sn)
		} else {
			_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(jobtype.DeferCloseHomestayOrder, payload))
			if err != nil {
				// logx.WithContext(l.ctx).Errorf("pay success notify user  insert queue fail err :%+v , sn : %s", err, homestayOrder.Sn)
			}
		}
	}

	return &pb.UpdateHomestayOrderTradeStateResp{
		Id:              homestayOrder.Id,
		UserId:          homestayOrder.UserId,
		Sn:              homestayOrder.Sn,
		TradeCode:       homestayOrder.TradeCode,
		Title:           homestayOrder.Title,
		LiveStartDate:   homestayOrder.LiveStartDate.Unix(),
		LiveEndDate:     homestayOrder.LiveEndDate.Unix(),
		OrderTotalPrice: homestayOrder.OrderTotalPrice,
	}, nil
}

func (l *UpdateHomestayOrderTradeStateLogic) verifyOrderTradeState(newTradeState, oldTradeState int64) error {

	if newTradeState == model.HomestayOrderTradeStateWaitPay {
		return errors.Wrapf(xerr.NewErrMsg("Changing this status is not supported"),
			"Changing this status is not supported newTradeState:%d , oldTradeState:%d",
			newTradeState, oldTradeState)
	}
	if newTradeState == model.HomestayOrderTradeStateCancel {

		if newTradeState != model.HomestayOrderTradeStateWaitPay {

			return errors.Wrapf(xerr.NewErrMsg("只有待支付的订单才能被取消"),
				"Only orders pending payment can be cancelled newTradeState:%d , oldTradeState:%d",
				newTradeState, oldTradeState)
		}
	} else if newTradeState == model.HomestayOrderTradeStateWaitUse {
		if oldTradeState != model.HomestayOrderTradeStateWaitPay {
			return errors.Wrapf(xerr.NewErrMsg("Only orders pending payment can change this status"),
				"Only orders pending payment can change this status newTradeState: %d, oldTradeState: %d",
				newTradeState,
				oldTradeState)
		}
	} else if newTradeState == model.HomestayOrderTradeStateUsed {
		if oldTradeState != model.HomestayOrderTradeStateWaitUse {
			return errors.Wrapf(xerr.NewErrMsg("Only unused orders can be changed to this status"),
				"Only unused orders can be changed to this status newTradeState: %d, oldTradeState: %d",
				newTradeState,
				oldTradeState)
		}
	} else if newTradeState == model.HomestayOrderTradeStateRefund {
		if oldTradeState != model.HomestayOrderTradeStateWaitUse {
			return errors.Wrapf(xerr.NewErrMsg("Only unused orders can be changed to this status"),
				"Only unused orders can be changed to this status newTradeState: %d, oldTradeState: %d",
				newTradeState,
				oldTradeState)
		}
	} else if newTradeState == model.HomestayOrderTradeStateExpire {
		if oldTradeState != model.HomestayOrderTradeStateWaitUse {
			return errors.Wrapf(xerr.NewErrMsg("Only unused orders can be changed to this status"),
				"Only unused orders can be changed to this status newTradeState: %d, oldTradeState: %d",
				newTradeState,
				oldTradeState)
		}
	}

	return nil
}
