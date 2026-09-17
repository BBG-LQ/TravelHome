// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package homestay

import (
	"context"

	"TravelHome/app/travel/cmd/api/internal/svc"
	"TravelHome/app/travel/cmd/api/internal/types"
	"TravelHome/app/usercenter/model"
	"TravelHome/pkg/tool"
	"TravelHome/pkg/xerr"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type BusinessListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// boss all homestay room
func NewBusinessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BusinessListLogic {
	return &BusinessListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BusinessListLogic) BusinessList(req *types.BusinessListReq) (resp *types.BusinessListResp, err error) {
	// todo: add your logic here and delete this line

	//构建查询条件
	whereBuilder := l.svcCtx.HomestayModel.SelectBuilder().Where(squirrel.Eq{"homestay_business_id": req.HomestayBusinessId})

	//根据条件，分页查询
	list, err := l.svcCtx.HomestayModel.FindPageListByIdDESC(l.ctx, whereBuilder, req.LastId, req.PageSize)

	//判断查询内容是否出错
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "HomestayBusinessId: %d ,err : %v", req.HomestayBusinessId, err)
	}

	//返回结果，只拷贝单个项目
	var res []types.Homestay
	if len(list) > 0 {

		for _, homestay := range list {

			var tempHomestay types.Homestay
			_ = copier.Copy(&tempHomestay, homestay)

			//转换单位

			tempHomestay.FoodPrice = tool.Fen2Yuan(int64(homestay.FoodPrice))
			tempHomestay.HomestayPrice = tool.Fen2Yuan(int64(homestay.HomestayPrice))
			tempHomestay.MarketHomestayPrice = tool.Fen2Yuan(int64(homestay.MarketHomestayPrice))

			res = append(res, tempHomestay)

		}

	}

	return &types.BusinessListResp{
		List: res,
	}, nil
}
