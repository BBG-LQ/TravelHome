// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package thirdPayment

import (
	"context"

	"TravelHome/app/order/cmd/rpc/order"
	"TravelHome/app/payment/cmd/api/internal/svc"
	"TravelHome/app/payment/cmd/api/internal/types"
	"TravelHome/app/payment/cmd/rpc/payment"
	"TravelHome/app/payment/model"
	"TravelHome/app/usercenter/cmd/rpc/usercenter"
	usercentmodel "TravelHome/app/usercenter/model"
	"TravelHome/pkg/ctxdata"
	"TravelHome/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/zeromicro/go-zero/core/logx"
)

type ThirdPaymentwxPayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

var ErrWxPayError = xerr.NewErrMsg("wechat pay fail")

// third payment：wechat pay
func NewThirdPaymentwxPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThirdPaymentwxPayLogic {
	return &ThirdPaymentwxPayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThirdPaymentwxPayLogic) ThirdPaymentwxPay(req *types.ThirdPaymentWxPayReq) (resp *types.ThirdPaymentWxPayResp, err error) {
	// todo: add your logic here and delete this line

	var totalPrice int64   // Total amount paid for current order(cent)
	var description string // Current Payment Description.

	switch req.ServiceType {
	case model.ThirdPaymentServiceTypeHomestayOrder:
		//order 信息
		homestayTotalPrice, homestayDescription, err := l.getPayHomestayPriceDescription(req.OrderSn)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError),
				"getPayHomestayPriceDescription fail  error err%v,sn:%s", err, req.OrderSn)
		}

		totalPrice = homestayTotalPrice
		description = homestayDescription

	default:
		return nil, errors.Wrapf(xerr.NewErrMsg("Payment for this business type is not supported"), "Payment for this business type is not supported req: %+v", req)
	}

	// wechatPrepayRsp,
	_, err = l.createWxPrePayOrder(req.ServiceType, req.OrderSn, totalPrice, description)
	if err != nil {
		return nil, err
	}

	return &types.ThirdPaymentWxPayResp{}, nil
}

//*********************
//创建微信订单

func (l *ThirdPaymentwxPayLogic) createWxPrePayOrder(serviceType, orderSn string, totalPrice int64, description string) (*jsapi.PrepayWithRequestPaymentResponse, error) {

	//userRpc获得openid

	userid := ctxdata.GetUidFromCtx(l.ctx)

	userResp, err := l.svcCtx.UsercenterRpc.GetUserAuthByUserId(l.ctx, &usercenter.GetUserAuthByUserIdReq{
		UserId:   userid,
		AuthType: usercentmodel.UserAuthTypeSmallWX,
	})

	if err != nil {
		return nil, errors.Wrapf(ErrWxPayError,
			"Get user wechat openid err : %v , userId: %d , orderSn:%s", err, userid, orderSn)
	}

	if userResp.UserAuth == nil || userResp.UserAuth.Id == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("Get user wechat openid fail，Please pay before authorization by weChat"),
			"Get user WeChat openid does not exist  userId: %d , orderSn:%s", userid, orderSn)
	}

	openId := userResp.UserAuth.AuthKey

	//paymentRpc creat 第三方订单

	createPaymentResp, err := l.svcCtx.PaymentRpc.CreatePayment(l.ctx, &payment.CreatePaymentReq{
		UserId:      userid,
		PayMode:     model.ThirdPaymentPayModelWechatPay,
		PayTotal:    totalPrice,
		OrderSn:     orderSn,
		ServiceType: serviceType,
	})

	if err != nil || createPaymentResp.Sn == "" {
		return nil, errors.Wrapf(ErrWxPayError,
			"create local third payment record fail : err: %v , userId: %d,totalPrice: %d , orderSn: %s",
			err, userid, totalPrice, orderSn)
	}

	// 3、create wechat pay pre pay order

	wxPayClient, err := svc.NewWxPayClientV3(l.svcCtx.Config)
	if err != nil {
		return nil, err
	}

	jsApiSvc := jsapi.JsapiApiService{Client: wxPayClient}

	resp, _, err := jsApiSvc.PrepayWithRequestPayment(l.ctx, jsapi.PrepayRequest{
		Appid:       core.String(l.svcCtx.Config.WxMiniConf.AppId),
		Mchid:       core.String(l.svcCtx.Config.WxPayConf.MchId),
		Description: core.String(description),
		OutTradeNo:  core.String(createPaymentResp.Sn),
		Attach:      core.String(description),
		NotifyUrl:   core.String(l.svcCtx.Config.WxPayConf.NotifyUrl),

		Amount: &jsapi.Amount{
			Total: core.Int64(totalPrice),
		},
		Payer: &jsapi.Payer{
			Openid: core.String(openId),
		},
	})

	if err != nil {
		return nil, errors.Wrapf(ErrWxPayError, "Failed to initiate WeChat payment pre-order err : %v , userId: %d , orderSn:%s", err, userid, orderSn)
	}

	return resp, nil
}

// ***********************
// 获取order订单
func (l *ThirdPaymentwxPayLogic) getPayHomestayPriceDescription(orderSn string) (int64, string, error) {
	// todo: add your logic here and delete this line

	description := "homestay pay"

	//order 信息
	resp, err := l.svcCtx.OrderRpc.HomestayOrderDetail(l.ctx, &order.HomestayOrderDetailReq{
		Sn: orderSn,
	})

	if err != nil {
		return 0, description, errors.Wrapf(xerr.NewErrCode(xerr.DbError),
			"getPayHomestayPriceDescription fail  error err%v,sn:%s", err, orderSn)
	}

	if resp.HomestayOrder == nil || resp.HomestayOrder.Id == 0 {
		return 0, description, errors.Wrapf(xerr.NewErrMsg("order no exists"),
			"order no exists,sn:%s", orderSn)
	}

	return resp.HomestayOrder.HomestayTotalPrice, description, nil
}
