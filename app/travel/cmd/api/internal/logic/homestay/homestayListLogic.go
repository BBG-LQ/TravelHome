// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package homestay

import (
	"context"

	"TravelHome/app/travel/cmd/api/internal/svc"
	"TravelHome/app/travel/cmd/api/internal/types"
	"TravelHome/app/travel/model"
	"TravelHome/pkg/tool"
	"TravelHome/pkg/xerr"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type HomestayListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// homestay room list
func NewHomestayListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayListLogic {
	return &HomestayListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayListLogic) HomestayList(req *types.HomestayListReq) (resp *types.HomestayListResp, err error) {
	// todo: add your logic here and delete this line

	//构建查询条件

	whereBuilder := l.svcCtx.HomestayActivityModel.SelectBuilder().Where(squirrel.Eq{
		"row_type":   model.HomestayActivityPreferredType,
		"row_status": model.HomestayActivityUpStatus,
	})

	//分页查询

	list, err := l.svcCtx.HomestayActivityModel.FindPageListByPage(l.ctx, whereBuilder, req.Page, req.PageSize, "data_id desc")

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "get activity homestay id set fail rowType: %s ,err : %v", model.HomestayActivityPreferredType, err)
	}
	//xiang应切片
	var res []types.Homestay

	if len(list) > 0 {

		//mapreduce
		mr.MapReduceVoid(

			//map阶段
			func(source chan<- interface{}) {
				for _, homestayActivity := range list {
					source <- homestayActivity.DataId
				}
			},

			//Reduce阶段
			func(item interface{}, writer mr.Writer[*model.Homestay], cancel func(error)) {
				id := item.(int64)

				homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, id)

				if err != nil && err != model.ErrNotFound {
					logx.WithContext(l.ctx).Errorf("ActivityHomestayListLogic ActivityHomestayList 获取活动数据失败 id : %d ,err : %v", id, err)

					return
				}

				writer.Write(homestay)
			},

			//收集阶段
			func(pipe <-chan *model.Homestay, cancel func(error)) {

				for homestay := range pipe {
					var tyHomestay types.Homestay
					_ = copier.Copy(&tyHomestay, homestay)

					tyHomestay.FoodPrice = tool.Fen2Yuan(int64(homestay.FoodPrice))
					tyHomestay.HomestayPrice = tool.Fen2Yuan(int64(homestay.HomestayPrice))
					tyHomestay.MarketHomestayPrice = tool.Fen2Yuan(int64(homestay.MarketHomestayPrice))

					res = append(res, tyHomestay)
				}
			},
		)

	}

	//返回结果
	return &types.HomestayListResp{List: res}, nil
}
