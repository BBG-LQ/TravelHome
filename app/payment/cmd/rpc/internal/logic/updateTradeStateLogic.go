package logic

import (
	"context"
	"time"

	"TravelHome/app/payment/cmd/rpc/internal/svc"
	"TravelHome/app/payment/cmd/rpc/pb"
	"TravelHome/app/payment/model"
	"TravelHome/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTradeStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTradeStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTradeStateLogic {
	return &UpdateTradeStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新交易状态
func (l *UpdateTradeStateLogic) UpdateTradeState(in *pb.UpdateTradeStateReq) (*pb.UpdateTradeStateResp, error) {
	// todo: add your logic here and delete this line

	//查询第三方订单
	thirdPayment, err := l.svcCtx.ThirdPaymentModel.FindOneBySn(l.ctx, in.Sn)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "getpayment fail findbysn error err%v,sn:%s", err, in.Sn)
	}

	if thirdPayment == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg(" thirdPayment record no exists "), "getpayment fail sn:%s", in.Sn)
	}
	//验证状态

	if in.PayStatus == model.ThirdPaymentPayTradeStateSuccess || in.PayStatus == model.ThirdPaymentPayTradeStateFAIL {
		if thirdPayment.PayStatus != model.ThirdPaymentPayTradeStateWait {
			return &pb.UpdateTradeStateResp{}, nil
		}

	} else if in.PayStatus == model.ThirdPaymentPayTradeStateRefund {
		if thirdPayment.PayStatus != model.ThirdPaymentPayTradeStateSuccess {
			return nil, errors.Wrapf(xerr.NewErrMsg(" Only orders with successful payment can be refunded"), "Only orders with successful payment can be refunded in : %+v", in)
		}
	} else {
		return nil, errors.Wrapf(xerr.NewErrMsg(" This status is not currently supported"), "Modify payment flow status is not supported  in : %+v", in)
	}

	//update 状态

	thirdPayment.TradeState = in.TradeState
	thirdPayment.TradeType = in.TradeType
	thirdPayment.TransactionId = in.TransactionId
	thirdPayment.TradeStateDesc = in.TradeStateDesc
	thirdPayment.PayStatus = in.PayStatus
	thirdPayment.PayTime = time.Unix(in.PayTime, 0)

	if err := l.svcCtx.ThirdPaymentModel.UpdateWithVersion(l.ctx, nil, thirdPayment); err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "UpdateWithVersion error err%v,hirdPayment %+v,in:%+v", err, thirdPayment, in)
	}

	//4、notify  sub "payment-update-paystatus-topic"  services(order-mq ..), pub、sub use kq

	return &pb.UpdateTradeStateResp{}, nil
}

/*
function pubKqPaySuccess

*/
