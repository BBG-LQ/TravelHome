// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"TravelHome/app/mqueue/cmd/scheduler/internal/svc"
	"TravelHome/app/mqueue/cmd/scheduler/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchedulerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSchedulerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchedulerLogic {
	return &SchedulerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SchedulerLogic) Scheduler(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
