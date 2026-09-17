// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package thirdPayment

import (
	"context"
	"net/http"

	"TravelHome/app/payment/cmd/api/internal/svc"
	"TravelHome/app/payment/cmd/api/internal/types"
	"TravelHome/app/payment/cmd/rpc/payment"
	"TravelHome/app/payment/model"
	"TravelHome/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/zeromicro/go-zero/core/logx"
)

var ErrWxPayCallbackError = xerr.NewErrMsg("wechat pay callback fail")

type ThirdPaymentWxPayCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// third payment：wechat pay callback
func NewThirdPaymentWxPayCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThirdPaymentWxPayCallbackLogic {
	return &ThirdPaymentWxPayCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThirdPaymentWxPayCallbackLogic) ThirdPaymentWxPayCallback(rw http.ResponseWriter, req *http.Request) (resp *types.ThirdPaymentWxPayCallbackResp, err error) {
	// todo: add your logic here and delete this line

	// 1. 初始化微信支付客户端V3（加载商户私钥等配置）
	_, err = svc.NewWxPayClientV3(l.svcCtx.Config)
	if err != nil {
		return nil, err
	}

	// 2. 获取平台证书访问器（用于验证微信回调签名）
	// 通过证书管理器获取对应商户号的平台证书，用于验签

	certVisitor := downloader.MgrInstance().GetCertificateVisitor(l.svcCtx.Config.WxPayConf.MchId)

	// 3. 创建通知处理器，传入APIv3密钥和SHA256-RSA验签器

	handler, err := notify.NewRSANotifyHandler(l.svcCtx.Config.WxPayConf.APIv3Key, verifiers.NewSHA256WithRSAVerifier(certVisitor))
	if err != nil {
		return nil, err
	}

	// 4. 解析回调请求，验证签名并反序列化交易数据到Transaction结构

	transaction := new(payments.Transaction)

	_, err = handler.ParseNotifyRequest(context.Background(), req, transaction)
	if err != nil {
		// 解析失败，返回错误（微信会收到非SUCCESS状态）
		return nil, errors.Wrapf(ErrWxPayCallbackError, "Failed to parse data ,err:%v", err)
	}
	// 5. 默认返回成功状态，后续根据业务处理结果决定是否变更
	returnCode := "SUCCESS"
	// 6. 执行业务校验和状态更新
	err = l.verifyAndUpdateState(transaction)
	if err != nil {
		returnCode = "FAIL"
	}

	return &types.ThirdPaymentWxPayCallbackResp{ReturnCode: returnCode}, err
}

// 校验回调数据并更新本地订单状态
func (l *ThirdPaymentWxPayCallbackLogic) verifyAndUpdateState(notifyTrasaction *payments.Transaction) error {

	// 1. 根据商户订单号（OutTradeNo）查询本地支付流水记录

	paymentResp, err := l.svcCtx.PaymentRpc.GetPaymentBySn(l.ctx, &payment.GetPaymentBySnReq{
		Sn: *notifyTrasaction.OutTradeNo,
	})
	if err != nil || paymentResp.PaymentDetail.Id == 0 {
		return errors.Wrapf(ErrWxPayCallbackError, "Failed to get payment flow record err:%v ,notifyTrasaction:%+v ", err, notifyTrasaction)
	}

	// 2. 核对订单金额（校验回调中的支付总金额是否与本地记录一致）

	notifyPayTotal := *notifyTrasaction.Amount.PayerTotal
	if paymentResp.PaymentDetail.PayTotal != notifyPayTotal {
		return errors.Wrapf(ErrWxPayCallbackError, "Order amount exception  notifyPayTotal:%v , notifyTrasaction:%v ", notifyPayTotal, notifyTrasaction)
	}

	// 3. 根据微信交易状态映射到本地支付状态

	payStatus := l.getPayStatusByWXPayTradeState(*notifyTrasaction.TradeState)

	// 4. 如果交易状态为支付成功,若本地已经是成功状态，则直接返回（幂等处理）

	if payStatus == model.ThirdPaymentPayTradeStateSuccess {
		// 若本地已经是成功状态，则直接返回（幂等处理）
		if paymentResp.PaymentDetail.PayStatus != model.ThirdPaymentPayTradeStateWait {
			return nil
		}

		// 5. 更新本地流水状态：包括交易状态、微信订单号、交易类型、状态描述、本地状态

		if _, err = l.svcCtx.PaymentRpc.UpdateTradeState(l.ctx, &payment.UpdateTradeStateReq{
			Sn:             *notifyTrasaction.OutTradeNo,
			TradeState:     *notifyTrasaction.TradeState,
			TransactionId:  *notifyTrasaction.TransactionId,
			TradeType:      *notifyTrasaction.TradeType,
			TradeStateDesc: *notifyTrasaction.TradeStateDesc,
			PayStatus:      l.getPayStatusByWXPayTradeState(*notifyTrasaction.TradeState),
		}); err != nil {
			return errors.Wrapf(ErrWxPayCallbackError, "更新流水状态失败  err:%v , notifyTrasaction:%v ", err, notifyTrasaction)
		}

	} else if payStatus == model.ThirdPaymentPayTradeStateWait {
		// 6. 如果状态为支付中（如用户支付中），此处暂不处理（退款通知等留待后续）
		// 实际业务中可在此处理支付中或退款等其他状态
	}

	return nil
}

// 微信支付API返回的交易状态常量
const (
	SUCCESS    = "SUCCESS"    // 支付成功
	REFUND     = "REFUND"     // 转入退款
	NOTPAY     = "NOTPAY"     // 未支付
	CLOSED     = "CLOSED"     // 已关闭
	REVOKED    = "REVOKED"    // 已撤销（付款码支付）
	USERPAYING = "USERPAYING" // 用户支付中（付款码支付）
	PAYERROR   = "PAYERROR"   // 支付失败（其他原因）
)

// 将微信交易状态映射为本地支付状态（int64）

func (l *ThirdPaymentWxPayCallbackLogic) getPayStatusByWXPayTradeState(wxPayTradeState string) int64 {
	switch wxPayTradeState {
	case SUCCESS: // 支付成功 → 成功状态
		return model.ThirdPaymentPayTradeStateSuccess
	case USERPAYING: // 支付中 → 等待状态
		return model.ThirdPaymentPayTradeStateWait
	case REFUND: // 已退款 → 暂视为等待状态（可根据业务调整）
		return model.ThirdPaymentPayTradeStateWait
	default: // 其他（未支付、已关闭、撤销、失败等）→ 失败状态
		return model.ThirdPaymentPayTradeStateFAIL
	}
}
