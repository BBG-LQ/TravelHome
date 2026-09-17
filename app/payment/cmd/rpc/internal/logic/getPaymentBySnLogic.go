package logic

import (
	"context"

	"TravelHome/app/payment/cmd/rpc/internal/svc"
	"TravelHome/app/payment/cmd/rpc/pb"
	"TravelHome/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentBySnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaymentBySnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentBySnLogic {
	return &GetPaymentBySnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据sn查询流水记录
func (l *GetPaymentBySnLogic) GetPaymentBySn(in *pb.GetPaymentBySnReq) (*pb.GetPaymentBySnResp, error) {
	// todo: add your logic here and delete this line

	thirdPayment, err := l.svcCtx.ThirdPaymentModel.FindOneBySn(l.ctx, in.Sn)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "getpayment fail findbysn error err%v,sn:%s", err, in.Sn)
	}

	var resp pb.PaymentDetail
	if thirdPayment != nil {

		_ = copier.Copy(&resp, thirdPayment)

		resp.CreateTime = thirdPayment.CreateTime.Unix()
		resp.UpdateTime = thirdPayment.UpdateTime.Unix()
		resp.PayTime = thirdPayment.PayTime.Unix()
	}

	return &pb.GetPaymentBySnResp{
		PaymentDetail: &resp,
	}, nil
}
