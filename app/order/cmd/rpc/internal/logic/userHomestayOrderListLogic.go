package logic

import (
	"context"

	"TravelHome/app/order/cmd/rpc/internal/svc"
	"TravelHome/app/order/cmd/rpc/pb"
	"TravelHome/app/order/model"
	"TravelHome/pkg/xerr"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserHomestayOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserHomestayOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserHomestayOrderListLogic {
	return &UserHomestayOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户民宿订单
func (l *UserHomestayOrderListLogic) UserHomestayOrderList(in *pb.UserHomestayOrderListReq) (*pb.UserHomestayOrderListResp, error) {
	// todo: add your logic here and delete this line

	whereBuilder := l.svcCtx.HomestayOrderModel.SelectBuilder().Where(squirrel.Eq{
		"user_id": in.UserId,
	})

	if in.TraderState >= model.HomestayOrderTradeStateCancel && in.TraderState <= model.HomestayOrderTradeStateExpire {
		whereBuilder = whereBuilder.Where(squirrel.Eq{"trade_state": in.TraderState})
	}

	list, err := l.svcCtx.HomestayOrderModel.FindPageListByIdDESC(l.ctx, whereBuilder, in.LastId, in.PageSize)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "Failed to get user's homestay order err : %v , in :%+v", err, in)
	}

	var resp []*pb.HomestayOrder
	if len(list) > 0 {

		for _, homestayOrder := range list {

			var tempItem pb.HomestayOrder

			_ = copier.Copy(&tempItem, homestayOrder)

			tempItem.CreateTime = homestayOrder.CreateTime.Unix()
			tempItem.LiveStartDate = homestayOrder.LiveStartDate.Unix()
			tempItem.LiveEndDate = homestayOrder.LiveEndDate.Unix()

			resp = append(resp, &tempItem)

		}

	}
	return &pb.UserHomestayOrderListResp{
		List: resp,
	}, nil
}
