package logic

import (
	"TravelHome/app/mqueue/cmd/job/internal/svc"
	"TravelHome/app/mqueue/cmd/job/jobtype"
	"context"

	"github.com/hibiken/asynq"
)

type CronJob struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCronJob(ctx context.Context, svcCtx *svc.ServiceContext) *CronJob {
	return &CronJob{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CronJob) Register() *asynq.ServeMux {
	mux := asynq.NewServeMux()

	mux.Handle(jobtype.ScheduleSettleRecord, NewSettleRecordHander(l.svcCtx))

	mux.Handle(jobtype.DeferCloseHomestayOrder, NewCloseOrderHander(l.svcCtx))

	mux.Handle(jobtype.MsgPaySuccessNotifyUser, NewPaySuccessNotifyUser(l.svcCtx))

	return mux

}
