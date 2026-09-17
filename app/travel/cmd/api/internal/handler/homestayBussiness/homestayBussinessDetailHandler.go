// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package homestayBussiness

import (
	"net/http"

	"TravelHome/app/travel/cmd/api/internal/logic/homestayBussiness"
	"TravelHome/app/travel/cmd/api/internal/svc"
	"TravelHome/app/travel/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// boss detail
func HomestayBussinessDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HomestayBussinessDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := homestayBussiness.NewHomestayBussinessDetailLogic(r.Context(), svcCtx)
		resp, err := l.HomestayBussinessDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
