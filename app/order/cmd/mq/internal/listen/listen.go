package listen

import (
	"TravelHome/app/order/cmd/mq/internal/config"
	"TravelHome/app/order/cmd/mq/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/service"
)

func Mqs(c config.Config) []service.Service {

	svcContext := svc.NewServiceContext(c)

	ctx := context.Background()

	var services []service.Service

	services = append(services, KqMqs(c, ctx, svcContext)...)

	return services

}
