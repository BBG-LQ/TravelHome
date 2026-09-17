// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"TravelHome/app/usercenter/cmd/api/internal/svc"
	"TravelHome/app/usercenter/cmd/api/internal/types"
	"TravelHome/app/usercenter/cmd/rpc/usercenter"
	"TravelHome/pkg/ctxdata"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get user info
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	// todo: add your logic here and delete this line

	//rpc 获取detail
	result, err := l.svcCtx.UsercenterRpc.GetUserInfo(l.ctx, &usercenter.GetUserInfoReq{

		Id: ctxdata.GetUidFromCtx(l.ctx),
	})

	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}

	//返回结构体

	var userInforesp types.UserInfoResp
	_ = copier.Copy(&userInforesp, result)

	return &userInforesp, nil
}
