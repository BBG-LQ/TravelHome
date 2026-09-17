// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"TravelHome/app/usercenter/cmd/api/internal/svc"
	"TravelHome/app/usercenter/cmd/api/internal/types"
	"TravelHome/app/usercenter/cmd/rpc/usercenter"
	"TravelHome/app/usercenter/model"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// login
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line

	//rpc 登录
	result, err := l.svcCtx.UsercenterRpc.Login(l.ctx, &usercenter.LoginReq{

		Password: req.Password,
		AuthKey:  req.Mobile,
		AuthType: model.UserAuthTypeSystem,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}

	//返回结构体

	var loginresp types.LoginResp
	_ = copier.Copy(&loginresp, result)

	return &loginresp, nil
}
