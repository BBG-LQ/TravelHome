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

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GuessListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// guess homestay room
func NewGuessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GuessListLogic {
	return &GuessListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GuessListLogic) GuessList(req *types.GuessListReq) (resp *types.GuessListResp, err error) {
	// todo: add your logic here and delete this line

	list, err := l.svcCtx.HomestayModel.FindPageListByIdDESC(l.ctx, l.svcCtx.HomestayModel.SelectBuilder(), 0, 5)

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "GuessList db err req : %+v , err : %v", req, err)
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

	return &types.GuessListResp{List: res}, nil
}
