package svc

import (
	"TravelHome/app/travel/cmd/rpc/internal/config"
	"TravelHome/app/travel/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	BizRedis      *redis.Redis
	HomestayModel model.HomestayModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:        c,
		BizRedis:      redis.New(c.BizRedis.Host),
		HomestayModel: model.NewHomestayModel(sqlConn, c.Cache),
	}
}
