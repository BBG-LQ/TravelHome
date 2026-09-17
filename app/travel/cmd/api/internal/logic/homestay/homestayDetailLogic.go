// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package homestay

import (
	"context"

	"TravelHome/app/travel/cmd/api/internal/svc"
	"TravelHome/app/travel/cmd/api/internal/types"
	"TravelHome/app/travel/cmd/rpc/travel"
	"TravelHome/pkg/tool"
	"TravelHome/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// homestay room detail
func NewHomestayDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayDetailLogic {
	return &HomestayDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayDetailLogic) HomestayDetail(req *types.HomestayDetailReq) (resp *types.HomestayDetailResp, err error) {
	// todo: add your logic here and delete this line

	homestayResp, err := l.svcCtx.TravelRpc.HomestayDetail(l.ctx, &travel.HomestayDetailReq{
		Id: req.Id,
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("get homestay detail fail"), " get homestay detail db err , id : %d , err : %v ", req.Id, err)
	}

	var tempHomestay types.Homestay
	if homestayResp.Homestay != nil {

		_ = copier.Copy(&tempHomestay, homestayResp.Homestay)

		//转换单位

		tempHomestay.FoodPrice = tool.Fen2Yuan(homestayResp.Homestay.FoodPrice)
		tempHomestay.HomestayPrice = tool.Fen2Yuan(int64(homestayResp.Homestay.HomestayPrice))
		tempHomestay.MarketHomestayPrice = tool.Fen2Yuan(int64(homestayResp.Homestay.MarketHomestayPrice))

	}

	return &types.HomestayDetailResp{Homestay: tempHomestay}, nil
}
