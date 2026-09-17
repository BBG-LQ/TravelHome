package logic

import (
	"context"

	"TravelHome/app/usercenter/cmd/rpc/internal/svc"
	"TravelHome/app/usercenter/cmd/rpc/pb"
	"TravelHome/app/usercenter/cmd/rpc/usercenter"
	"TravelHome/app/usercenter/model"
	"TravelHome/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserAuthByAuthKeyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserAuthByAuthKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAuthByAuthKeyLogic {
	return &GetUserAuthByAuthKeyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserAuthByAuthKeyLogic) GetUserAuthByAuthKey(in *pb.GetUserAuthByAuthKeyReq) (*pb.GetUserAuthByAuthKeyResp, error) {
	// todo: add your logic here and delete this line

	//1.查询是否存在 根据AuthByUserId
	user, err := l.svcCtx.UserAuthModel.FindOneByAuthKeyAuthType(l.ctx, in.AuthKey, in.AuthType)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "GetUserAuthByAuthKey db err , in:%v , err:%v", in, err)
	}

	// 2.返回结构体
	var respUserAuth usercenter.UserAuth
	_ = copier.Copy(&respUserAuth, user)
	return &pb.GetUserAuthByAuthKeyResp{
		UserAuth: &respUserAuth,
	}, nil
}
