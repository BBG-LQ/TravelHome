// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"TravelHome/app/mqueue/cmd/scheduler/internal/logic"
	"TravelHome/app/mqueue/cmd/scheduler/internal/svc"
	"TravelHome/app/mqueue/cmd/scheduler/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SchedulerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSchedulerLogic(r.Context(), svcCtx)
		resp, err := l.Scheduler(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
