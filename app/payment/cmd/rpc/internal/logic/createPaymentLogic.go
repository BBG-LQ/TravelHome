package logic

import (
	"context"
	"time"

	"TravelHome/app/payment/cmd/rpc/internal/svc"
	"TravelHome/app/payment/cmd/rpc/pb"
	"TravelHome/app/payment/model"
	"TravelHome/pkg/uniqueid"
	"TravelHome/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentLogic {
	return &CreatePaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建微信支付预处理订单
func (l *CreatePaymentLogic) CreatePayment(in *pb.CreatePaymentReq) (*pb.CreatePaymentResp, error) {
	// todo: add your logic here and delete this line
	data := new(model.ThirdPayment)

	data.Sn = uniqueid.GenSn(uniqueid.SN_PREFIX_THIRD_PAYMENT)
	data.OrderSn = in.OrderSn
	data.UserId = in.UserId
	data.PayMode = in.PayMode
	data.PayTotal = in.PayTotal
	data.ServiceType = model.ThirdPaymentServiceTypeHomestayOrder

	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	data.PayTime = time.Now()

	_, err := l.svcCtx.ThirdPaymentModel.Insert(l.ctx, nil, data)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "create insert error err%v,data%+v", err, data)
	}

	return &pb.CreatePaymentResp{Sn: data.Sn}, nil
}
