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
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var ErrUserAlreadyRegisterError = xerr.NewErrCodeMsg(100009, "user has been registered")

// 注册
func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterResp, error) {
	// todo: add your logic here and delete this line

	//1.判断在不在，根据mobel查询
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "Register user exists mobile:%s,err:%v", in.Mobile, err)
	}
	if user != nil {
		return nil, errors.Wrapf(ErrUserAlreadyRegisterError, "Register user exists mobile:%s,err:%v", in.Mobile, err)
	}

	var UserID int64
	//2.事务trans 进行insert

	if err := l.svcCtx.UserModel.Trans(l.ctx, func(context context.Context, session sqlx.Session) error {
		user := new(model.User)

		//1.判断输入合法 Nikename password
		user.Mobile = in.Mobile

		if len(in.Nickname) == 0 {
			user.Nickname = tool.Krand(8, tool.KC_RAND_KIND_ALL)
		}

		if len(in.Password) > 0 {
			user.Password = tool.Md5ByString(in.Password)
		}

		//2.insert进入user
		insertresult, err := l.svcCtx.UserModel.Insert(l.ctx, session, user)
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DbError), "Register user insert err :%v,err:%v", user, err)
		}

		// 获取userid
		lastId, err := insertresult.LastInsertId()
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DbError), "Register user insertResult.lastID err :%v,err:%v", user, err)
		}
		UserID = lastId

		//*******userAuth****

		//插入结构体
		userAuth := new(model.UserAuth)
		userAuth.UserId = uint64(lastId)
		userAuth.AuthKey = in.AuthKey
		userAuth.AuthType = in.AuthType

		//insert

		_, err = l.svcCtx.UserAuthModel.Insert(l.ctx, session, userAuth)
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DbError), "Register userAuth insert err :%v,err:%v", userAuth, err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	//3.生成token 返回
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)

	tokenResult, err := generateTokenLogic.GenerateToken(&usercenter.GenerateTokenReq{
		UserId: UserID,
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.TokenGenerateError), "TokenGenerateError err :%v,err:%v", err, UserID)
	}

	return &pb.RegisterResp{
		AccessToken:  tokenResult.AccessToken,
		AccessExpire: tokenResult.AccessExpire,
		RefreshAfter: tokenResult.RefreshAfter,
	}, nil
}
