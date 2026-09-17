// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package homestayBussiness

import (
	"context"

	"TravelHome/app/travel/cmd/api/internal/svc"
	"TravelHome/app/travel/cmd/api/internal/types"
	"TravelHome/app/travel/model"
	"TravelHome/app/usercenter/cmd/rpc/usercenter"
	"TravelHome/pkg/xerr"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type GoodBossLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// good boss
func NewGoodBossLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GoodBossLogic {
	return &GoodBossLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GoodBossLogic) GoodBoss(req *types.GoodBossReq) (resp *types.GoodBossResp, err error) {
	// todo: add your logic here and delete this line

	//构建查询条件

	whereBuilder := l.svcCtx.HomestayActivityModel.SelectBuilder().Where(squirrel.Eq{
		"row_type":   model.HomestayActivityPreferredType,
		"row_status": model.HomestayActivityUpStatus,
	})

	//分页查询

	list, err := l.svcCtx.HomestayActivityModel.FindPageListByPage(l.ctx, whereBuilder, 0, 10, "data_id desc")

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "get activity homestay id set fail rowType: %s ,err : %v", model.HomestayActivityPreferredType, err)
	}
	//xiang应切片
	var res []types.HomestayBusinessBoss

	if len(list) > 0 {

		//mapreduce
		mr.MapReduceVoid(

			//map阶段
			func(source chan<- interface{}) {
				for _, homestayActivity := range list {
					source <- homestayActivity.DataId
				}
			}, func(item interface{}, writer mr.Writer[*usercenter.User], cancel func(error)) {
				id := item.(int64)

				userRes, err := l.svcCtx.UsercenterRpc.GetUserInfo(l.ctx, &usercenter.GetUserInfoReq{
					Id: id,
				})

				if err != nil {
					// logx.WithContext(l.ctx).Errorf("ActivityHomestayListLogic ActivityHomestayList 获取活动数据失败 id : %d ,err : %v", id, err)

					return
				}

				if userRes.User != nil && userRes.User.Id > 0 {
					writer.Write(userRes.User)
				}

			}, func(pipe <-chan *usercenter.User, cancel func(error)) {

				for item := range pipe {
					var typesHomestayBusiness types.HomestayBusinessBoss
					_ = copier.Copy(&typesHomestayBusiness, item)

					res = append(res, typesHomestayBusiness)
				}
			})

	}

	return &types.GoodBossResp{List: res}, nil
}
