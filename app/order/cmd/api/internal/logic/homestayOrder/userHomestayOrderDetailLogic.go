// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package homestayOrder

import (
	"context"

	"TravelHome/app/order/cmd/api/internal/svc"
	"TravelHome/app/order/cmd/api/internal/types"
	"TravelHome/app/order/cmd/rpc/order"
	"TravelHome/pkg/ctxdata"
	"TravelHome/pkg/tool"
	"TravelHome/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserHomestayOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户订单明细
func NewUserHomestayOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserHomestayOrderDetailLogic {
	return &UserHomestayOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserHomestayOrderDetailLogic) UserHomestayOrderDetail(req *types.UserHomestayOrderDetailReq) (resp *types.UserHomestayOrderDetailResp, err error) {
	// todo: add your logic here and delete this line

	userId := ctxdata.GetUidFromCtx(l.ctx)

	res, err := l.svcCtx.OrderRpc.HomestayOrderDetail(l.ctx, &order.HomestayOrderDetailReq{
		Sn: req.Sn,
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("get homestay order detail fail"), "get homestay order detail fail  rpc HomestayDetail fail ,  err : %v,sn:%s", err, req.Sn)

	}

	var typesOrderDetail types.UserHomestayOrderDetailResp
	if res.HomestayOrder != nil && res.HomestayOrder.UserId == userId {

		_ = copier.Copy(&typesOrderDetail, res.HomestayOrder)

		// format price.
		typesOrderDetail.OrderTotalPrice = tool.Fen2Yuan(res.HomestayOrder.OrderTotalPrice)
		typesOrderDetail.FoodTotalPrice = tool.Fen2Yuan(res.HomestayOrder.FoodTotalPrice)
		typesOrderDetail.HomestayTotalPrice = tool.Fen2Yuan(res.HomestayOrder.HomestayTotalPrice)
		typesOrderDetail.HomestayPrice = tool.Fen2Yuan(res.HomestayOrder.HomestayPrice)
		typesOrderDetail.FoodPrice = tool.Fen2Yuan(res.HomestayOrder.FoodPrice)
		typesOrderDetail.MarketHomestayPrice = tool.Fen2Yuan(res.HomestayOrder.MarketHomestayPrice)

		/*

			paymentRpc getSuccessdetail


		*/

		return &typesOrderDetail, nil

	}

	return nil, nil
}
