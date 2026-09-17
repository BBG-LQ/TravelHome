package svc

import (
	"TravelHome/app/payment/cmd/rpc/internal/config"
	"TravelHome/app/payment/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	ThirdPaymentModel model.ThirdPaymentModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		ThirdPaymentModel: model.NewThirdPaymentModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}
