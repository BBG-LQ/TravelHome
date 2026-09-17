// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package thirdPayment

import (
	"net/http"

	"TravelHome/app/payment/cmd/api/internal/logic/thirdPayment"
	"TravelHome/app/payment/cmd/api/internal/svc"
	"TravelHome/app/payment/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// third payment：wechat pay callback
func ThirdPaymentWxPayCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ThirdPaymentWxPayCallbackReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := thirdPayment.NewThirdPaymentWxPayCallbackLogic(r.Context(), svcCtx)
		resp, err := l.ThirdPaymentWxPayCallback(w, r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
