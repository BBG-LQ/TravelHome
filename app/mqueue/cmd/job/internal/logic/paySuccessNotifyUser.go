package logic

import (
	"TravelHome/app/mqueue/cmd/job/internal/svc"
	"context"

	"github.com/hibiken/asynq"
)

type PaySuccessNotifyUser struct {
	svcCtx *svc.ServiceContext
}

func NewPaySuccessNotifyUser(svcCtx *svc.ServiceContext) *PaySuccessNotifyUser {
	return &PaySuccessNotifyUser{

		svcCtx: svcCtx,
	}
}

func (l *PaySuccessNotifyUser) ProcessTask(ctx context.Context, t *asynq.Task) error {

	// close logic

	//
	return nil
}
