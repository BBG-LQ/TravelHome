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

type UserHomestayOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户订单列表
func NewUserHomestayOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserHomestayOrderListLogic {
	return &UserHomestayOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserHomestayOrderListLogic) UserHomestayOrderList(req *types.UserHomestayOrderListReq) (resp *types.UserHomestayOrderListResp, err error) {
	// todo: add your logic here and delete this line

	userId := ctxdata.GetUidFromCtx(l.ctx)

	res, err := l.svcCtx.OrderRpc.UserHomestayOrderList(l.ctx, &order.UserHomestayOrderListReq{
		UserId:      userId,
		TraderState: req.TradeState,
		LastId:      req.LastId,
		PageSize:    req.PageSize,
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("get homestay order list fail"), "get homestay order detail fail  rpc Homestaylist fail ,  err : %v,req:%+v", err, req)

	}

	var respList []types.UserHomestayOrderListView
	if len(res.List) > 0 {
		for _, homestayOrder := range res.List {

			var temphomestayOrder types.UserHomestayOrderListView
			_ = copier.Copy(&temphomestayOrder, homestayOrder)

			temphomestayOrder.OrderTotalPrice = tool.Fen2Yuan(homestayOrder.OrderTotalPrice)

			respList = append(respList, temphomestayOrder)

		}

	}

	return &types.UserHomestayOrderListResp{List: respList}, nil
}
