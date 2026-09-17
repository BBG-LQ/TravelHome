// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package homestayBussiness

import (
	"context"

	"TravelHome/app/travel/cmd/api/internal/svc"
	"TravelHome/app/travel/cmd/api/internal/types"
	"TravelHome/app/travel/model"
	"TravelHome/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayBussinessListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// business list
func NewHomestayBussinessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayBussinessListLogic {
	return &HomestayBussinessListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayBussinessListLogic) HomestayBussinessList(req *types.HomestayBussinessListReq) (resp *types.HomestayBussinessListResp, err error) {
	// todo: add your logic here and delete this line

	//构建查询条件
	whereBuilder := l.svcCtx.HomestayBusinessModel.SelectBuilder()

	//根据条件，分页查询
	list, err := l.svcCtx.HomestayBusinessModel.FindPageListByIdDESC(l.ctx, whereBuilder, req.LastId, req.PageSize)

	//判断查询内容是否出错
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "HomestayBussinessList FindPageListByIdDESC db fail ,  req : %+v , err:%v", req, err)
	}

	//返回结果，只拷贝单个项目
	var res []types.HomestayBusinessListInfo
	if len(list) > 0 {

		for _, item := range list {

			var tempHomestayBusinessListInfo types.HomestayBusinessListInfo
			_ = copier.Copy(&tempHomestayBusinessListInfo, item)

			res = append(res, tempHomestayBusinessListInfo)

		}

	}
	return &types.HomestayBussinessListResp{List: res}, nil
}
