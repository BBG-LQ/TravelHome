package logic

import (
	"context"

	"TravelHome/app/travel/cmd/rpc/internal/svc"
	"TravelHome/app/travel/cmd/rpc/pb"
	"TravelHome/app/travel/model"
	"TravelHome/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestayDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayDetailLogic {
	return &HomestayDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// homestayDetail
func (l *HomestayDetailLogic) HomestayDetail(in *pb.HomestayDetailReq) (*pb.HomestayDetailResp, error) {
	// todo: add your logic here and delete this line

	homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, in.Id)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), " HomestayDetail db err , id : %d ", in.Id)
	}

	var pbhomestay pb.Homestay
	if homestay != nil {
		_ = copier.Copy(&pbhomestay, homestay)
	}

	return &pb.HomestayDetailResp{
		Homestay: &pbhomestay,
	}, nil
}
