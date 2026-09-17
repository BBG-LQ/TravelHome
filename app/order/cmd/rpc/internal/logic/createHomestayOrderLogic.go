package logic

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"TravelHome/app/mqueue/cmd/job/jobtype"
	"TravelHome/app/order/cmd/rpc/internal/svc"
	"TravelHome/app/order/cmd/rpc/pb"
	"TravelHome/app/order/model"
	"TravelHome/app/travel/cmd/rpc/travel"
	"TravelHome/pkg/tool"
	"TravelHome/pkg/uniqueid"
	"TravelHome/pkg/xerr"

	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

const CloseOrderTimeMinutes = 30 //defer close order time

type CreateHomestayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateHomestayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateHomestayOrderLogic {
	return &CreateHomestayOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 民宿下订单
func (l *CreateHomestayOrderLogic) CreateHomestayOrder(in *pb.CreateHomestayOrderReq) (*pb.CreateHomestayOrderResp, error) {
	// todo: add your logic here and delete this line

	//判断时间

	if in.LiveEndTime <= in.LiveStartTime {
		return nil, errors.Wrapf(xerr.NewErrMsg("Stay at least one night"), "Place an order at a B&B. The end time of your stay must be greater than the start time. in : %+v", in)
	}

	//获取民宿并判断

	resp, err := l.svcCtx.TravelRpc.HomestayDetail(l.ctx, &travel.HomestayDetailReq{Id: in.HomestayId})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("Failed to query the record"), "Failed to query the record  rpc HomestayDetail fail , homestayId : %d , err : %v", in.HomestayId, err)

	}

	if resp.Homestay == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("This record does not exist"), "This record does not exist HomestayDetail fail , homestayId : %d , ", in.HomestayId)

	}

	//订单数据，判断

	var cover string

	if len(resp.Homestay.Banner) > 0 {
		cover = strings.Split(resp.Homestay.Banner, ",")[0]
	}

	order := new(model.HomestayOrder)

	order.Sn = uniqueid.GenSn(uniqueid.SN_PREFIX_HOMESTAY_ORDER)
	order.UserId = in.UserId
	order.HomestayId = in.HomestayId
	order.Title = resp.Homestay.Title
	order.SubTitle = resp.Homestay.SubTitle
	order.Cover = cover
	order.Info = resp.Homestay.Info
	order.PeopleNum = resp.Homestay.PeopleNum
	order.RowType = resp.Homestay.RowType
	order.HomestayPrice = resp.Homestay.HomestayPrice
	order.MarketHomestayPrice = resp.Homestay.MarketHomestayPrice
	order.HomestayBusinessId = resp.Homestay.HomestayBusinessId
	order.HomestayUserId = resp.Homestay.UserId
	order.LivePeopleNum = in.LivePeopleNum

	order.TradeState = model.HomestayOrderTradeStateWaitPay
	order.TradeCode = tool.Krand(8, tool.KC_RAND_KIND_ALL)
	order.Remark = in.Remark
	order.FoodInfo = resp.Homestay.FoodInfo
	order.FoodPrice = resp.Homestay.FoodPrice
	order.LiveStartDate = time.Unix(in.LiveStartTime, 0)
	order.LiveEndDate = time.Unix(in.LiveEndTime, 0)

	liveDays := int64(order.LiveEndDate.Sub(order.LiveStartDate).Seconds() / 86400)
	order.HomestayTotalPrice = int64(resp.Homestay.HomestayPrice * liveDays)

	if in.IsFood {
		order.NeedFood = model.HomestayOrderNeedFoodYes
		order.FoodTotalPrice = int64(resp.Homestay.FoodPrice * in.LivePeopleNum * liveDays)
	}

	order.OrderTotalPrice = order.HomestayTotalPrice + order.FoodTotalPrice
	//创建订单 insert

	_, err = l.svcCtx.HomestayOrderModel.Insert(l.ctx, nil, order)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "Order Database Exception order : %+v , err: %v", order, err)
	}

	//pay app/mqueue/cmd/job/internal/logic/closeOrder.go有定义asynq的延迟关闭订单任务

	payload, err := json.Marshal(jobtype.DeferCloseHomestayOrderPayload{
		Sn: order.Sn,
	})
	if err != nil {
		// logx.WithContext(l.ctx).Errorf("create defer close order task json Marshal fail err :%+v , sn : %s", err, order.Sn)
	} else {
		_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(jobtype.DeferCloseHomestayOrder, payload), asynq.ProcessIn(CloseOrderTimeMinutes*time.Minute))
		if err != nil {
			// logx.WithContext(l.ctx).Errorf("create defer close order task insert queue fail err :%+v , sn : %s", err, order.Sn)
		}
	}

	return &pb.CreateHomestayOrderResp{Sn: order.Sn}, nil
}
