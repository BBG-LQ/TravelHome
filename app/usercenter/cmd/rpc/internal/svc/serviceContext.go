package svc

import (
	"TravelHome/app/usercenter/cmd/rpc/internal/config"
	"TravelHome/app/usercenter/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	UserModel     model.UserModel
	UserAuthModel model.UserAuthModel
	BizRedis      *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	SqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:        c,
		UserModel:     model.NewUserModel(SqlConn, c.Cache),
		UserAuthModel: model.NewUserAuthModel(SqlConn, c.Cache),
		BizRedis:      redis.New(c.BizRedis.Host),
	}
}
