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

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var ErrUserNoExistsError = xerr.NewErrMsg("用户不存在")

// 获取用户详情（需认证，用户 ID 从上下文获取）
func (l *GetUserInfoLogic) GetUserInfo(in *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	// todo: add your logic here and delete this line

	//1.查询是否存在 根据id
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "GetUserInfo find user db err , id:%d , err:%v", in.Id, err)
	}
	if user == nil {
		return nil, errors.Wrapf(ErrUserNoExistsError, " id:%d ", in.Id)
	}

	// 2.返回结构体
	var respUser usercenter.User
	_ = copier.Copy(&respUser, user)

	return &pb.GetUserInfoResp{
		User: &respUser,
	}, nil
}
