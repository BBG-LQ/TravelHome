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

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayBussinessDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// boss detail
func NewHomestayBussinessDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayBussinessDetailLogic {
	return &HomestayBussinessDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayBussinessDetailLogic) HomestayBussinessDetail(req *types.HomestayBussinessDetailReq) (resp *types.HomestayBussinessDetailResp, err error) {
	// todo: add your logic here and delete this line

	homestayBusiness, err := l.svcCtx.HomestayBusinessModel.FindOne(l.ctx, req.Id)

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "HomestayBussinessDetail  FindOne db fail ,id  : %d , err : %v", req.Id, err)
	}

	var typeHomestayBusinessBoss types.HomestayBusinessBoss
	if homestayBusiness != nil {
		userResp, err := l.svcCtx.UsercenterRpc.GetUserInfo(l.ctx, &usercenter.GetUserInfoReq{
			Id: int64(homestayBusiness.UserId),
		})

		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("get boss info fail"), "get boss info fail ,  userId : %d ,err:%v", homestayBusiness.UserId, err)
		}

		if userResp != nil && userResp.User.Id > 0 {
			_ = copier.Copy(&typeHomestayBusinessBoss, userResp.User)
		}
	}

	return &types.HomestayBussinessDetailResp{
		Boss: typeHomestayBusinessBoss,
	}, nil
}
