package listen

import (
	"TravelHome/app/order/cmd/mq/internal/config"
	kqmq "TravelHome/app/order/cmd/mq/internal/mqs/kq"
	"TravelHome/app/order/cmd/mq/internal/svc"
	"context"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

func KqMqs(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {

	return []service.Service{
		kq.MustNewQueue(c.PaymentUpdateStatusConf, kqmq.NewPaymentUpdateStatusMq(svcContext)),
	}

}
