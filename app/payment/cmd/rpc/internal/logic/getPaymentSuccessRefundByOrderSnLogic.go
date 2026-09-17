package logic

import (
	"context"

	"TravelHome/app/payment/cmd/rpc/internal/svc"
	"TravelHome/app/payment/cmd/rpc/pb"
	"TravelHome/app/payment/model"
	"TravelHome/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentSuccessRefundByOrderSnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaymentSuccessRefundByOrderSnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentSuccessRefundByOrderSnLogic {
	return &GetPaymentSuccessRefundByOrderSnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据订单sn查询流水记录
func (l *GetPaymentSuccessRefundByOrderSnLogic) GetPaymentSuccessRefundByOrderSn(in *pb.GetPaymentSuccessRefundByOrderSnReq) (*pb.GetPaymentSuccessRefundByOrderSnResp, error) {
	// todo: add your logic here and delete this line

	whereBuild := l.svcCtx.ThirdPaymentModel.SelectBuilder().Where(
		"order_sn = ? and (trade_state = ? or trade_state = ? )",
		in.OrderSn, model.ThirdPaymentPayTradeStateSuccess, model.ThirdPaymentPayTradeStateRefund,
	)

	thirdPayments, err := l.svcCtx.ThirdPaymentModel.FindAll(l.ctx, whereBuild, "id desc")
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "PaymentSuccessRefundByOrderSn fail findbysn error err%v,sn:%s", err, in.OrderSn)
	}

	var resp pb.PaymentDetail
	if len(thirdPayments) > 0 {

		thirdPayment := thirdPayments[0]
		if thirdPayment != nil {
			_ = copier.Copy(&resp, thirdPayment)

			resp.CreateTime = thirdPayment.CreateTime.Unix()
			resp.UpdateTime = thirdPayment.UpdateTime.Unix()
			resp.PayTime = thirdPayment.PayTime.Unix()

		}

	}

	return &pb.GetPaymentSuccessRefundByOrderSnResp{PaymentDetail: &resp}, nil
}
