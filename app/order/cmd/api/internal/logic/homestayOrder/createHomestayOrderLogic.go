// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package homestayOrder

import (
	"context"

	"TravelHome/app/order/cmd/api/internal/svc"
	"TravelHome/app/order/cmd/api/internal/types"
	"TravelHome/app/order/cmd/rpc/order"
	"TravelHome/app/travel/cmd/rpc/travel"
	"TravelHome/pkg/ctxdata"
	"TravelHome/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateHomestayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建民宿订单
func NewCreateHomestayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateHomestayOrderLogic {
	return &CreateHomestayOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateHomestayOrderLogic) CreateHomestayOrder(req *types.CreateHomestayOrderReq) (resp *types.CreateHomestayOrderResp, err error) {
	// todo: add your logic here and delete this line

	//获取民宿并判断

	homestayDetail, err := l.svcCtx.TravelRpc.HomestayDetail(l.ctx, &travel.HomestayDetailReq{
		Id: req.HomestayId,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("Failed to query the record"), "Failed to query the record  rpc HomestayDetail fail , homestayId : %d , err : %v", req.HomestayId, err)

	}

	if homestayDetail.Homestay == nil || homestayDetail.Homestay.Id == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("This record does not exist"), "This record does not exist HomestayDetail fail , homestayId : %d , ", req.HomestayId)

	}

	//uid

	userId := ctxdata.GetUidFromCtx(l.ctx)

	//orderRpc creatOrder
	res, err := l.svcCtx.OrderRpc.CreateHomestayOrder(l.ctx, &order.CreateHomestayOrderReq{
		HomestayId:    req.HomestayId,
		IsFood:        req.IsFood,
		LiveStartTime: req.LiveStartTime,
		LiveEndTime:   req.LiveEndTime,
		UserId:        userId,
		LivePeopleNum: req.LivePeopleNum,
		Remark:        req.Remark,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("create homestay order fail"), "create homestay order rpc CreateHomestayOrder fail req: %+v , err : %v ", req, err)
	}
	return &types.CreateHomestayOrderResp{
		OrderSn: res.Sn,
	}, nil
}
