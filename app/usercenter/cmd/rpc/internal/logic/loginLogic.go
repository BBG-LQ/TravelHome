package logic

import (
	"context"

	"TravelHome/app/usercenter/cmd/rpc/internal/svc"
	"TravelHome/app/usercenter/cmd/rpc/pb"
	"TravelHome/app/usercenter/cmd/rpc/usercenter"
	"TravelHome/app/usercenter/model"
	"TravelHome/pkg/tool"
	"TravelHome/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var ErrGenerateTokenError = xerr.NewErrMsg("生成token失败")

var ErrUsernamePwdError = xerr.NewErrMsg("账号密码错误")

// 登录
func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	// todo: add your logic here and delete this line

	var UserID int64
	var err error
	//1.平台登录

	switch in.AuthType {
	case model.UserAuthTypeSystem:
		UserID, err = l.LoginByMobel(in.AuthKey, in.Password)
	// case model.UserAuthTypeSmallWX:
	// 	err = l.loginBySmallWx()
	default:
		return nil, xerr.NewErrCode(xerr.ServerCommonError)
	}

	//2.生成token 返回
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)

	tokenResult, err := generateTokenLogic.GenerateToken(&usercenter.GenerateTokenReq{
		UserId: UserID,
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.TokenGenerateError), "TokenGenerateError err :%v,err:%v", UserID, err)
	}

	return &pb.LoginResp{
		AccessToken:  tokenResult.AccessToken,
		AccessExpire: tokenResult.AccessExpire,
		RefreshAfter: tokenResult.RefreshAfter,
	}, nil
}

func (l *LoginLogic) LoginByMobel(mobile string, password string) (userID int64, err error) {

	//1.查询mobile是否存在
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, mobile)
	if err != nil && err != model.ErrNotFound {
		return 0, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "查询失败 mobile:%s,err:%v", mobile, err)
	}
	if user == nil {
		return 0, errors.Wrapf(ErrUserAlreadyRegisterError, "mobile:%s,err:%v", mobile, err)
	}

	//2.验证possword

	if !(user.Password == tool.Md5ByString(password)) {
		return 0, errors.Wrapf(ErrUsernamePwdError, "密码匹配错误")
	}

	return int64(user.Id), nil
}
func (l *LoginLogic) loginBySmallWx() error {
	return nil
}
