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

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// register
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// todo: add your logic here and delete this line

	//rpc注册
	result, err := l.svcCtx.UsercenterRpc.Register(l.ctx, &usercenter.RegisterReq{
		Mobile:   req.Mobile,
		Password: req.Password,
		AuthKey:  req.Mobile,
		AuthType: model.UserAuthTypeSystem,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}

	//返回结构体

	var Registerresp types.RegisterResp
	_ = copier.Copy(&Registerresp, result)

	return &Registerresp, nil
}
