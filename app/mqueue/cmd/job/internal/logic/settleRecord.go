package logic

import (
	"TravelHome/app/mqueue/cmd/job/internal/svc"
	"context"
	"fmt"

	"github.com/hibiken/asynq"
)

type SettleRecordHander struct {
	svcCtx *svc.ServiceContext
}

func NewSettleRecordHander(svcCtx *svc.ServiceContext) *SettleRecordHander {
	return &SettleRecordHander{

		svcCtx: svcCtx,
	}
}

func (l *SettleRecordHander) ProcessTask(ctx context.Context, _ *asynq.Task) error {

	fmt.Printf("shcedule job demo -----> every one minute exec \n")
	return nil
}
